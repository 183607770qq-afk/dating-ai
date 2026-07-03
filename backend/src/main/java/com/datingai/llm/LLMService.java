package com.datingai.llm;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;

import java.io.IOException;
import java.util.Map;

@Service
public class LLMService {

    private static final Logger logger = LoggerFactory.getLogger(LLMService.class);

    private final Map<String, DatingAdviceClient> adviceClients;
    private final String provider;

    public LLMService(
            Map<String, DatingAdviceClient> adviceClients,
            @Value("${llm.provider:langchain4j}") String provider) {
        this.adviceClients = adviceClients;
        this.provider = provider;
    }

    public String getDatingAdvice(String question) throws IOException {
        if (!StringUtils.hasText(question)) {
            return "请先告诉我你的具体困惑，比如聊天开场、约会安排、关系推进或冲突沟通。";
        }

        DatingAdviceClient adviceClient = adviceClients.get(provider);
        if (adviceClient == null) {
            throw new IOException("Unknown LLM provider: " + provider);
        }

        logger.info("Getting dating advice with provider: {}", provider);

        try {
            /*
             * LLMService 是门面层：Controller 只调用它，不关心底层使用哪种实现。
             * 通过 llm.provider 可以在 legacy 和 langchain4j 之间切换，便于学习和回滚。
             */
            return adviceClient.getAdvice(question);
        } catch (Exception e) {
            logger.error("Error calling LLM provider: {}", provider, e);
            throw new IOException("Failed to get advice from LLM: " + e.getMessage(), e);
        }
    }
}
