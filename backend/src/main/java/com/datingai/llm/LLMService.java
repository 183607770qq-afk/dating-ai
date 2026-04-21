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
import org.springframework.stereotype.Service;

import java.io.IOException;
import java.util.HashMap;
import java.util.Map;

@Service
public class LLMService {

    @Value("${llm.api.url}")
    private String apiUrl;

    @Value("${llm.api.key}")
    private String apiKey;

    private final ObjectMapper objectMapper = new ObjectMapper();

    public String getDatingAdvice(String question) throws IOException {
        try (CloseableHttpClient httpClient = HttpClients.createDefault()) {
            HttpPost httpPost = new HttpPost(apiUrl);
            httpPost.setHeader("Authorization", "Bearer " + apiKey);
            httpPost.setHeader("Content-Type", "application/json");

            Map<String, Object> requestBody = new HashMap<>();
            requestBody.put("model", "deepseek-chat");
            requestBody.put("messages", new Object[]{
                Map.of("role", "system", "content", "你是一个专业的情感顾问，专注于帮助用户解决脱单问题。请提供具体、实用的建议，帮助用户提升社交技能和建立健康的恋爱关系。"),
                Map.of("role", "user", "content", question)
            });
            requestBody.put("temperature", 0.7);

            String jsonBody = objectMapper.writeValueAsString(requestBody);
            httpPost.setEntity(new StringEntity(jsonBody, ContentType.APPLICATION_JSON));

            try (CloseableHttpResponse response = httpClient.execute(httpPost)) {
                JsonNode responseNode = objectMapper.readTree(response.getEntity().getContent());
                return responseNode.get("choices").get(0).get("message").get("content").asText();
            }
        }
    }
}
