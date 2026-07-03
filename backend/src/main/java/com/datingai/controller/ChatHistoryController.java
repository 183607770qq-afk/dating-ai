package com.datingai.controller;

import com.datingai.model.ChatMessage;
import com.datingai.model.User;
import com.datingai.repository.UserRepository;
import com.datingai.service.ChatHistoryService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.*;

import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

@RestController
@RequestMapping("/api/chat")
public class ChatHistoryController {

    @Autowired
    private ChatHistoryService chatHistoryService;

    @Autowired
    private UserRepository userRepository;

    /**
     * 获取当前用户的历史消息，分页返回（按时间倒序）
     */
    @GetMapping("/history")
    public ResponseEntity<Map<String, Object>> getHistory(
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size) {

        User user = getCurrentUser();
        if (user == null) {
            return ResponseEntity.status(401).body(Map.of("error", "请先登录"));
        }

        List<ChatMessage> messages = chatHistoryService.getHistoryMessages(user.getId(), page, size);
        long total = chatHistoryService.getMessageCount(user.getId());

        List<Map<String, Object>> messageList = messages.stream().map(msg -> {
            Map<String, Object> map = new HashMap<>();
            map.put("id", msg.getId());
            map.put("role", msg.getRole());
            map.put("content", msg.getContent());
            map.put("createdAt", msg.getCreatedAt().toString());
            return map;
        }).collect(Collectors.toList());

        Map<String, Object> result = new HashMap<>();
        result.put("messages", messageList);
        result.put("total", total);
        result.put("page", page);
        result.put("size", size);

        return ResponseEntity.ok(result);
    }

    /**
     * 获取用户最近的 N 条消息（按时间升序，用于页面加载时恢复对话）
     */
    @GetMapping("/recent")
    public ResponseEntity<Map<String, Object>> getRecent(
            @RequestParam(defaultValue = "20") int limit) {

        User user = getCurrentUser();
        if (user == null) {
            return ResponseEntity.status(401).body(Map.of("error", "请先登录"));
        }

        List<ChatMessage> messages = chatHistoryService.getRecentMessages(user.getId(), limit);

        List<Map<String, Object>> messageList = messages.stream().map(msg -> {
            Map<String, Object> map = new HashMap<>();
            map.put("id", msg.getId());
            map.put("role", msg.getRole());
            map.put("content", msg.getContent());
            map.put("createdAt", msg.getCreatedAt().toString());
            return map;
        }).collect(Collectors.toList());

        Map<String, Object> result = new HashMap<>();
        result.put("messages", messageList);

        return ResponseEntity.ok(result);
    }

    /**
     * 从 SecurityContext 获取当前登录用户
     */
    private User getCurrentUser() {
        Authentication auth = SecurityContextHolder.getContext().getAuthentication();
        if (auth == null || !auth.isAuthenticated()) {
            return null;
        }
        String username = auth.getName();
        if (username == null || "anonymousUser".equals(username)) {
            return null;
        }
        return userRepository.findByUsername(username);
    }
}
