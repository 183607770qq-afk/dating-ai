package com.datingai.service;

import com.datingai.model.ChatMessage;
import com.datingai.model.User;
import com.datingai.repository.ChatMessageRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.PageRequest;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

@Service
public class ChatHistoryService {

    @Autowired
    private ChatMessageRepository chatMessageRepository;

    /**
     * 保存一条聊天消息
     */
    public ChatMessage saveMessage(User user, String role, String content) {
        ChatMessage message = new ChatMessage(user, role, content);
        return chatMessageRepository.save(message);
    }

    /**
     * 获取用户最近的 N 条历史消息，按时间升序排列（用于拼入 LLM 上下文）
     */
    public List<ChatMessage> getRecentMessages(Long userId, int limit) {
        List<ChatMessage> messages = chatMessageRepository
                .findByUserIdOrderByCreatedAtDesc(userId, PageRequest.of(0, limit));
        List<ChatMessage> result = new ArrayList<>(messages);
        Collections.reverse(result);
        return result;
    }

    /**
     * 分页查询用户历史消息，按时间倒序（用于前端展示）
     */
    public List<ChatMessage> getHistoryMessages(Long userId, int page, int size) {
        return chatMessageRepository
                .findByUserIdOrderByCreatedAtDesc(userId, PageRequest.of(page, size));
    }

    /**
     * 获取用户消息总数
     */
    public long getMessageCount(Long userId) {
        return chatMessageRepository.countByUserId(userId);
    }
}
