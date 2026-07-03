package com.datingai.llm;

import com.datingai.model.ChatMessage;
import com.datingai.model.User;
import com.datingai.repository.UserRepository;
import com.datingai.service.ChatHistoryService;
import com.datingai.service.SummaryService;
import com.datingai.service.UserMemoryService;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.hc.client5.http.classic.methods.HttpPost;
import org.apache.hc.client5.http.impl.classic.CloseableHttpClient;
import org.apache.hc.client5.http.impl.classic.HttpClients;
import org.apache.hc.core5.http.ContentType;
import org.apache.hc.core5.http.io.entity.StringEntity;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.web.servlet.mvc.method.annotation.SseEmitter;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.nio.charset.StandardCharsets;
import java.util.*;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

@Service
public class StreamLLMService {

    private static final Logger logger = LoggerFactory.getLogger(StreamLLMService.class);

    private static final String SYSTEM_PROMPT = "你是一个专业的情感顾问，专注于帮助用户解决脱单问题。请提供具体、实用的建议，帮助用户提升社交技能和建立健康的恋爱关系。";

    /** 历史消息窗口的 token 预算上限（中英文混合估算，chars/2≈tokens） */
    private static final int MAX_HISTORY_TOKENS = 4000;

    /** 单条历史消息最长字符数，超出部分截断并追加省略标记 */
    private static final int MAX_SINGLE_MSG_CHARS = 1000;

    @Value("${llm.api.url}")
    private String apiUrl;

    @Value("${llm.model}")
    private String model;

    @Autowired
    private ChatHistoryService chatHistoryService;

    @Autowired
    private UserRepository userRepository;

    @Autowired
    private SummaryService summaryService;

    @Autowired
    private UserMemoryService userMemoryService;

    private final ObjectMapper objectMapper = new ObjectMapper();
    private final ExecutorService executorService = Executors.newCachedThreadPool();

    public SseEmitter getStreamingAdvice(String question, Long userId) {
        SseEmitter emitter = new SseEmitter(300000L);

        executorService.execute(() -> {
            StringBuilder fullResponse = new StringBuilder();

            try (CloseableHttpClient httpClient = HttpClients.createDefault()) {
                HttpPost httpPost = new HttpPost(apiUrl);
                httpPost.setHeader("Content-Type", "application/json");

                Map<String, Object> requestBody = new HashMap<>();
                requestBody.put("model", model);
                requestBody.put("stream", true);

                // 构建 messages 数组：system + 历史记录 + 当前问题
                List<Map<String, String>> messages = buildMessages(question, userId);
                requestBody.put("messages", messages);
                requestBody.put("temperature", 0.7);

                String jsonBody = objectMapper.writeValueAsString(requestBody);
                httpPost.setEntity(new StringEntity(jsonBody, ContentType.APPLICATION_JSON));

                try (var response = httpClient.execute(httpPost);
                     InputStream inputStream = response.getEntity().getContent();
                     BufferedReader reader = new BufferedReader(new InputStreamReader(inputStream, StandardCharsets.UTF_8))) {

                    String line;
                    while ((line = reader.readLine()) != null) {
                        if (line.isEmpty()) continue;

                        String data = line.startsWith("data:") ? line.substring(5).trim() : line.trim();

                        if ("[DONE]".equals(data)) {
                            emitter.send(SseEmitter.event().data("[DONE]"));
                            break;
                        }

                        String content = extractStreamingContent(data);
                        if (content != null && !content.isEmpty()) {
                            fullResponse.append(content);
                            Map<String, String> result = new HashMap<>();
                            result.put("content", content);
                            emitter.send(SseEmitter.event().data(result));
                        }
                    }

                    emitter.complete();
                }

            } catch (IOException e) {
                logger.error("Streaming LLM request failed", e);
                emitter.completeWithError(e);
                return; // 出错时不保存消息
            }

            // 流式完成后，持久化聊天记录
            if (userId != null && fullResponse.length() > 0) {
                saveChatHistory(userId, question, fullResponse.toString());
                // 异步检查是否需要更新对话摘要（长时记忆）
                User user = userRepository.findById(userId).orElse(null);
                if (user != null) {
                    summaryService.tryUpdateSummary(user);
                }
            }
        });

        return emitter;
    }

    /**
     * 构建发送给 LLM 的 messages 数组：
     * [system, ...摘要(长时记忆), ...Milvus相关记忆, ...历史消息(按 token 预算滑动窗口), user(当前问题)]
     *
     * 滑动窗口从最近一条开始向前取，直到 token 预算用完。
     * 单条消息超出 MAX_SINGLE_MSG_CHARS 会被截断。
     */
    private List<Map<String, String>> buildMessages(String question, Long userId) {
        List<Map<String, String>> messages = new ArrayList<>();

        // 1. system prompt
        messages.add(Map.of("role", "system", "content", SYSTEM_PROMPT));
        int usedTokens = estimateTokens(SYSTEM_PROMPT);

        // 2. 对话摘要（长时记忆，始终注入）
        if (userId != null) {
            String summary = summaryService.getLatestSummary(userId);
            if (summary != null && !summary.isBlank()) {
                String summaryContent = "【以下是关于这位用户的关键信息摘要，请在回答时参考，但不要主动提及「摘要」这两个字】\n" + summary;
                messages.add(Map.of("role", "system", "content", summaryContent));
                usedTokens += estimateTokens(summaryContent);
            }
        }

        // 3. Milvus 语义检索：从用户记忆向量库中拉取与当前问题相关的关键事实
        if (userId != null) {
            try {
                List<String> relevantFacts = userMemoryService.searchUserMemories(userId, question, 3);
                if (!relevantFacts.isEmpty()) {
                    StringBuilder factBlock = new StringBuilder("【以下是与用户当前问题相关的历史记忆，请在回答时自然引用，不要提「记忆库」等词】\n");
                    for (int i = 0; i < relevantFacts.size(); i++) {
                        factBlock.append("- ").append(relevantFacts.get(i)).append("\n");
                    }
                    messages.add(Map.of("role", "system", "content", factBlock.toString()));
                    usedTokens += estimateTokens(factBlock.toString());
                }
            } catch (Exception e) {
                logger.warn("Failed to search user memories for userId={}: {}", userId, e.getMessage());
            }
        }

        // 4. 按 token 预算加载历史消息（滑动窗口：从最新往前取）
        if (userId != null) {
            try {
                // 取足够多的候选消息（上限 50 条），再由 token 预算筛选
                List<ChatMessage> candidates = chatHistoryService.getRecentMessages(userId, 50);
                // 用 LinkedList 从最新往前插入，O(1) 头部插入
                java.util.LinkedList<Map<String, String>> historyWindow = new java.util.LinkedList<>();

                // 从最新往前遍历（candidates 已经是时间升序）
                for (int i = candidates.size() - 1; i >= 0; i--) {
                    ChatMessage msg = candidates.get(i);
                    String content = truncateContent(msg.getContent());
                    int tokens = estimateTokens(content);
                    if (usedTokens + tokens > MAX_HISTORY_TOKENS) {
                        break; // 预算用完，停止
                    }
                    historyWindow.addFirst(Map.of("role", msg.getRole(), "content", content));
                    usedTokens += tokens;
                }

                messages.addAll(historyWindow);
            } catch (Exception e) {
                logger.warn("Failed to load chat history for userId={}: {}", userId, e.getMessage());
            }
        }

        // 5. 当前用户问题（不做截断，确保完整传入）
        messages.add(Map.of("role", "user", "content", question));

        if (logger.isDebugEnabled()) {
            logger.debug("Messages built: {} messages, ~{} tokens (excl. current question)",
                    messages.size() - 1, usedTokens);
        }

        return messages;
    }

    /** 截断过长消息 */
    private String truncateContent(String content) {
        if (content == null || content.isEmpty()) return "";
        if (content.length() <= MAX_SINGLE_MSG_CHARS) return content;
        return content.substring(0, MAX_SINGLE_MSG_CHARS) + "…";
    }

    /** 简单 token 估算：中英文混合文本，约 2 个字符 ≈ 1 个 token */
    private int estimateTokens(String text) {
        if (text == null || text.isEmpty()) return 0;
        return (int) Math.ceil(text.length() / 2.0);
    }

    /**
     * 持久化本轮对话的用户消息和 AI 回复
     */
    private void saveChatHistory(Long userId, String question, String answer) {
        try {
            User user = userRepository.findById(userId).orElse(null);
            if (user == null) {
                logger.warn("User not found for id={}, skipping chat history save", userId);
                return;
            }
            chatHistoryService.saveMessage(user, "user", question);
            chatHistoryService.saveMessage(user, "assistant", answer);
            logger.debug("Saved chat history for userId={}", userId);
        } catch (Exception e) {
            logger.error("Failed to save chat history for userId={}: {}", userId, e.getMessage(), e);
        }
    }

    private String extractStreamingContent(String data) {
        try {
            JsonNode node = objectMapper.readTree(data);

            JsonNode messageNode = node.get("message");
            if (messageNode != null && messageNode.has("content")) {
                return messageNode.get("content").asText();
            }

            if (node.has("response")) {
                return node.get("response").asText();
            }

            JsonNode choicesNode = node.get("choices");
            if (choicesNode != null && choicesNode.isArray() && !choicesNode.isEmpty()) {
                JsonNode choice = choicesNode.get(0);
                JsonNode delta = choice.get("delta");
                if (delta != null && delta.has("content")) {
                    return delta.get("content").asText();
                }

                JsonNode message = choice.get("message");
                if (message != null && message.has("content")) {
                    return message.get("content").asText();
                }
            }
        } catch (Exception e) {
            return data;
        }

        return "";
    }
}
