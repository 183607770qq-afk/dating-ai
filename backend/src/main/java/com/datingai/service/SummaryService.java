package com.datingai.service;

import com.datingai.model.ChatMessage;
import com.datingai.model.ConversationSummary;
import com.datingai.model.User;
import com.datingai.repository.ConversationSummaryRepository;
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

import java.io.IOException;
import java.util.*;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.stream.Collectors;

@Service
public class SummaryService {

    private static final Logger logger = LoggerFactory.getLogger(SummaryService.class);

    /** 每隔多少条新消息触发一次摘要更新 */
    private static final int SUMMARY_INTERVAL = 4;

    /** 摘要用的系统提示词 */
    private static final String SUMMARY_SYSTEM_PROMPT = """
            你是一个信息提取助手。你的任务是从一段情感咨询对话中提炼用户档案。

            【第一步 — 生成摘要】
            用简洁的中文输出以下内容，每项 1-3 句话：
            1. 用户基本信息（年龄、性别、职业等，如果提到）
            2. 感情状态（单身/恋爱中/已婚/刚分手等）
            3. 核心困扰（用户当前最大的感情问题或焦虑）
            4. 性格特点（社恐/外向/敏感/理性等）
            5. 已讨论过的重要话题和给过的建议
            6. 用户未提及但未来值得探讨的方向

            只输出事实，不要编造。如果某项未被提及，写"未知"。
            总长度控制在 600 字以内。

            【第二步 — 提取关键事实】
            从对话中提取 3-8 条关键事实，每条独立成句，用于后续语义检索。
            每条事实应该是自包含的、可独立理解的短句，比如：
            - "用户提到他前女友喜欢徒步旅行"
            - "用户三个月前经历了一次分手，情绪低落"
            - "用户在泰国餐厅有过一次失败的约会经历"

            请按以下格式输出：
            [SUMMARY]
            摘要内容...
            [FACTS]
            - 事实1
            - 事实2
            - 事实3
            """;

    @Value("${llm.api.url}")
    private String apiUrl;

    @Value("${llm.model}")
    private String model;

    @Autowired
    private ChatHistoryService chatHistoryService;

    @Autowired
    private ConversationSummaryRepository summaryRepository;

    @Autowired
    private UserMemoryService userMemoryService;

    private final ObjectMapper objectMapper = new ObjectMapper();
    private final ExecutorService summaryExecutor = Executors.newSingleThreadExecutor(r -> {
        Thread t = new Thread(r, "summary-generator");
        t.setDaemon(true);
        return t;
    });

    /**
     * 获取用户的最新摘要，没有则返回 null
     */
    public String getLatestSummary(Long userId) {
        ConversationSummary summary = summaryRepository.findByUserId(userId);
        return summary != null ? summary.getSummaryText() : null;
    }

    /**
     * 检查是否需要更新摘要（异步触发）
     * 在主流程完成后调用，不阻塞响应
     */
    public void tryUpdateSummary(User user) {
        summaryExecutor.execute(() -> {
            try {
                long total = chatHistoryService.getMessageCount(user.getId());

                ConversationSummary existing = summaryRepository.findByUserId(user.getId());
                long lastSummaryCount = existing != null ? existing.getMessageCountAtSummary() : 0;
                long newMessages = total - lastSummaryCount;

                if (newMessages < SUMMARY_INTERVAL) {
                    logger.debug("Not enough new messages for summary: {} < {}", newMessages, SUMMARY_INTERVAL);
                    return;
                }

                logger.info("Triggering summary update for userId={}, total={}, new={}",
                        user.getId(), total, newMessages);

                // 取最近的对话用于生成摘要（取足够多但不过量）
                List<ChatMessage> recentMessages = chatHistoryService.getRecentMessages(user.getId(), 30);
                String conversationText = formatMessages(recentMessages);

                String previousSummary = existing != null ? existing.getSummaryText() : "";
                String newSummary = callLLMForSummary(previousSummary, conversationText);

                if (newSummary != null && !newSummary.isBlank()) {
                    // 解析 LLM 输出：分离摘要和关键事实
                    RawSummaryResult result = parseRawSummary(newSummary);

                    // upsert: 删除旧摘要，插入新摘要
                    summaryRepository.deleteByUserId(user.getId());
                    ConversationSummary summary = new ConversationSummary(user, result.summary, total);
                    summaryRepository.save(summary);
                    logger.info("Summary updated for userId={}, messages={}, length={}",
                            user.getId(), total, result.summary.length());

                    // 将关键事实向量化存入 Milvus
                    if (result.facts != null && !result.facts.isEmpty()) {
                        userMemoryService.storeUserFacts(user.getId(), result.facts);
                        logger.info("Stored {} user facts in Milvus for userId={}",
                                result.facts.size(), user.getId());
                    }
                }
            } catch (Exception e) {
                logger.error("Failed to update summary for userId={}: {}", user.getId(), e.getMessage(), e);
            }
        });
    }

    /**
     * 调用 LLM 生成对话摘要
     */
    private String callLLMForSummary(String previousSummary, String conversationText) {
        try (CloseableHttpClient httpClient = HttpClients.createDefault()) {
            HttpPost httpPost = new HttpPost(apiUrl);
            httpPost.setHeader("Content-Type", "application/json");

            // 构建摘要请求的 messages
            List<Map<String, String>> messages = new ArrayList<>();
            messages.add(Map.of("role", "system", "content", SUMMARY_SYSTEM_PROMPT));

            StringBuilder userContent = new StringBuilder();
            if (previousSummary != null && !previousSummary.isBlank()) {
                userContent.append("【上一轮摘要】\n").append(previousSummary).append("\n\n");
            }
            userContent.append("【最近对话】\n").append(conversationText);
            userContent.append("\n\n请根据以上内容更新用户档案摘要：");
            messages.add(Map.of("role", "user", "content", userContent.toString()));

            Map<String, Object> requestBody = new HashMap<>();
            requestBody.put("model", model);
            requestBody.put("messages", messages);
            requestBody.put("stream", false);
            requestBody.put("temperature", 0.3); // 低温度，确保稳定提取

            String jsonBody = objectMapper.writeValueAsString(requestBody);
            httpPost.setEntity(new StringEntity(jsonBody, ContentType.APPLICATION_JSON));

            try (var response = httpClient.execute(httpPost)) {
                String responseBody = new String(response.getEntity().getContent().readAllBytes(), "UTF-8");
                return extractContent(responseBody);
            }
        } catch (IOException e) {
            logger.error("Summary LLM call failed: {}", e.getMessage());
            return null;
        }
    }

    /**
     * 将消息列表格式化为可读文本
     */
    private String formatMessages(List<ChatMessage> messages) {
        return messages.stream()
                .map(m -> (m.getRole().equals("user") ? "用户：" : "AI：") + m.getContent())
                .collect(Collectors.joining("\n"));
    }

    /**
     * 从 LLM 响应中提取 content
     */
    private String extractContent(String responseBody) {
        try {
            JsonNode root = objectMapper.readTree(responseBody);

            // Ollama 格式
            if (root.has("message") && root.get("message").has("content")) {
                return root.get("message").get("content").asText();
            }

            // 通用 choices 格式
            if (root.has("choices") && root.get("choices").isArray()) {
                JsonNode choice = root.get("choices").get(0);
                if (choice.has("message") && choice.get("message").has("content")) {
                    return choice.get("message").get("content").asText();
                }
            }

            // 纯文本
            if (root.has("response")) {
                return root.get("response").asText();
            }

            logger.warn("Could not extract content from LLM response: {}", 
                    responseBody.length() > 200 ? responseBody.substring(0, 200) : responseBody);
        } catch (Exception e) {
            logger.warn("Failed to parse LLM response", e);
        }
        return null;
    }

    /**
     * 解析 LLM 返回的 [SUMMARY]...[FACTS] 结构
     */
    private RawSummaryResult parseRawSummary(String raw) {
        String summary = raw;
        List<String> facts = new ArrayList<>();

        int factsIdx = raw.indexOf("[FACTS]");
        if (factsIdx >= 0) {
            summary = raw.substring(0, factsIdx).trim();
            // 去掉 [SUMMARY] 标签
            int summaryIdx = summary.indexOf("[SUMMARY]");
            if (summaryIdx >= 0) {
                summary = summary.substring(summaryIdx + "[SUMMARY]".length()).trim();
            }

            // 解析事实列表
            String factsBlock = raw.substring(factsIdx + "[FACTS]".length()).trim();
            for (String line : factsBlock.split("\n")) {
                String trimmed = line.trim();
                // 去掉 "- "、"-"、"* " 等前缀
                if (trimmed.startsWith("- ")) {
                    trimmed = trimmed.substring(2).trim();
                } else if (trimmed.startsWith("-")) {
                    trimmed = trimmed.substring(1).trim();
                } else if (trimmed.startsWith("* ")) {
                    trimmed = trimmed.substring(2).trim();
                }
                if (!trimmed.isEmpty() && trimmed.length() > 3) {
                    facts.add(trimmed);
                }
            }
        }

        return new RawSummaryResult(summary, facts);
    }

    /** 摘要解析结果 */
    private record RawSummaryResult(String summary, List<String> facts) {}
}
