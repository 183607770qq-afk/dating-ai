package com.datingai.config;

import com.datingai.llm.DatingAdviceAssistant;
import dev.langchain4j.model.chat.ChatModel;
import dev.langchain4j.model.ollama.OllamaChatModel;
import dev.langchain4j.service.AiServices;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.time.Duration;

@Configuration
public class LangChain4jConfig {

    @Bean
    public ChatModel datingChatModel(
            @Value("${llm.base-url:http://localhost:11434}") String baseUrl,
            @Value("${llm.model:qwen3-vl:4b}") String modelName,
            @Value("${llm.temperature:0.7}") double temperature,
            @Value("${llm.timeout:PT60S}") String timeout) {

        /*
         * ChatModel 是 LangChain4j 对聊天模型的统一抽象。
         * 这里使用 OllamaChatModel，所以业务代码不需要自己拼 HTTP 请求、
         * 解析 JSON，也不依赖 Ollama 的 /api/chat 响应格式。
         */
        return OllamaChatModel.builder()
                .baseUrl(baseUrl)
                .modelName(modelName)
                .temperature(temperature)
                .timeout(Duration.parse(timeout))
                .build();
    }

    @Bean
    public DatingAdviceAssistant datingAdviceAssistant(ChatModel datingChatModel) {
        /*
         * AiServices 会为接口动态创建实现类：
         * 1. 读取接口上的 @SystemMessage / @UserMessage 等提示词注解
         * 2. 把方法参数转换为模型消息
         * 3. 调用 ChatModel 并把回复转换成方法返回值
         */
        return AiServices.builder(DatingAdviceAssistant.class)
                .chatModel(datingChatModel)
                .build();
    }
}
