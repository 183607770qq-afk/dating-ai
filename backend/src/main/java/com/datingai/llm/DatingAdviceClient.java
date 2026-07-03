package com.datingai.llm;

import java.io.IOException;

public interface DatingAdviceClient {

    /*
     * 统一的业务接口。
     * 不管底层是手写 HTTP、LangChain4j，还是以后接 Spring AI，
     * Controller 和 LLMService 都只依赖这个方法。
     */
    String getAdvice(String question) throws IOException;
}
