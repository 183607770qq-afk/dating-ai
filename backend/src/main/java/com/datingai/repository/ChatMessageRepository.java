package com.datingai.repository;

import com.datingai.model.ChatMessage;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface ChatMessageRepository extends JpaRepository<ChatMessage, Long> {

    /**
     * 按时间倒序查询用户聊天记录（用于前端分页展示）
     */
    List<ChatMessage> findByUserIdOrderByCreatedAtDesc(Long userId, Pageable pageable);

    /**
     * 统计用户消息总数
     */
    long countByUserId(Long userId);
}
