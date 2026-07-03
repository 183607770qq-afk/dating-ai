package com.datingai.llm;

import dev.langchain4j.service.SystemMessage;
import dev.langchain4j.service.UserMessage;

public interface DatingAdviceAssistant {

    /*
     * SystemMessage 相当于给 AI 设定“角色”和“回答边界”。
     * 这类长期稳定的规则建议放在系统提示词里，用户每次提问只传具体问题。
     */
    @SystemMessage("""
            你是一个专业、温和、务实的情感顾问，专注于帮助用户提升社交能力、
            识别健康关系、建立真诚沟通。回答必须具体、可执行、尊重边界，
            不鼓励PUA、操控、骚扰或侵犯隐私的行为。
            """)
    @UserMessage("""
            用户的问题是：
            {{question}}

            请用中文回答，并尽量包含：
            1. 对问题的简短判断
            2. 可以立刻执行的建议
            3. 需要避免的风险
            """)
    String getAdvice(String question);
}
