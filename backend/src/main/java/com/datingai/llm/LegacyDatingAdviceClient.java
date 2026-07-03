package com.datingai.llm;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.hc.client5.http.classic.methods.HttpPost;
import org.apache.hc.client5.http.impl.classic.CloseableHttpClient;
import org.apache.hc.client5.http.impl.classic.CloseableHttpResponse;
import org.apache.hc.client5.http.impl.classic.HttpClients;
import org.apache.hc.core5.http.ContentType;
import org.apache.hc.core5.http.io.entity.StringEntity;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import java.io.IOException;
import java.util.HashMap;
import java.util.Map;

@Component("legacy")
public class LegacyDatingAdviceClient implements DatingAdviceClient {

    @Value("${llm.api.url}")
    private String apiUrl;

    @Value("${llm.model}")
    private String model;

    @Value("${llm.temperature:0.7}")
    private double temperature;

    private final ObjectMapper objectMapper = new ObjectMapper();

    @Override
    public String getAdvice(String question) throws IOException {
        /*
         * legacy 实现保留原来的手写 HTTP 方式。
         * 它适合学习模型接口原始请求结构，也方便在 LangChain4j 出问题时快速切回。
         */
        try (CloseableHttpClient httpClient = HttpClients.createDefault()) {
            HttpPost httpPost = new HttpPost(apiUrl);
            httpPost.setHeader("Content-Type", "application/json");

            Map<String, Object> requestBody = new HashMap<>();
            requestBody.put("model", model);
            requestBody.put("stream", false);
            requestBody.put("messages", new Object[]{
                    Map.of("role", "system", "content", "你是一个专业的情感顾问，专注于帮助用户解决脱单问题。请提供具体、实用的建议，帮助用户提升社交技能和建立健康的恋爱关系。"),
                    Map.of("role", "user", "content", question)
            });
            requestBody.put("temperature", temperature);

            String jsonBody = objectMapper.writeValueAsString(requestBody);
            httpPost.setEntity(new StringEntity(jsonBody, ContentType.APPLICATION_JSON));

            try (CloseableHttpResponse response = httpClient.execute(httpPost)) {
                JsonNode responseNode = objectMapper.readTree(response.getEntity().getContent());
                return extractContent(responseNode);
            }
        }
    }

    private String extractContent(JsonNode responseNode) {
        JsonNode messageNode = responseNode.get("message");
        if (messageNode != null && messageNode.has("content")) {
            return messageNode.get("content").asText();
        }

        if (responseNode.has("response")) {
            return responseNode.get("response").asText();
        }

        JsonNode choicesNode = responseNode.get("choices");
        if (choicesNode != null && choicesNode.isArray() && !choicesNode.isEmpty()) {
            JsonNode message = choicesNode.get(0).get("message");
            if (message != null && message.has("content")) {
                return message.get("content").asText();
            }
        }

        return "AI 服务返回数据格式异常，请稍后再试。";
    }
}
