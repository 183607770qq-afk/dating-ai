package com.datingai.llm;

import org.springframework.stereotype.Component;

@Component("langchain4j")
public class LangChainDatingAdviceClient implements DatingAdviceClient {

    private final DatingAdviceAssistant datingAdviceAssistant;

    public LangChainDatingAdviceClient(DatingAdviceAssistant datingAdviceAssistant) {
        this.datingAdviceAssistant = datingAdviceAssistant;
    }

    @Override
    public String getAdvice(String question) {
        /*
         * LangChain4j 实现：业务代码只调用接口方法。
         * 提示词模板和模型调用细节都由 DatingAdviceAssistant + AiServices 处理。
         */
        return datingAdviceAssistant.getAdvice(question);
    }
}
