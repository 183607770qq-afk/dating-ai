package com.datingai.controller;

import com.datingai.llm.LLMService;
import com.datingai.model.User;
import com.datingai.repository.UserRepository;
import com.datingai.service.ChatHistoryService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.*;

import java.io.IOException;
import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/llm")
public class LLMController {

    @Autowired
    private LLMService llmService;

    @Autowired
    private ChatHistoryService chatHistoryService;

    @Autowired
    private UserRepository userRepository;

    @PostMapping("/advice")
    public ResponseEntity<Map<String, String>> getDatingAdvice(@RequestBody Map<String, String> request) {
        try {
            String question = request.get("question");
            String advice = llmService.getDatingAdvice(question);

            // 持久化聊天记录
            User user = getCurrentUser();
            if (user != null) {
                chatHistoryService.saveMessage(user, "user", question);
                chatHistoryService.saveMessage(user, "assistant", advice);
            }

            return ResponseEntity.ok(Map.of("advice", advice));
        } catch (IOException e) {
            return ResponseEntity.status(500).body(Map.of("error", "Failed to get advice"));
        }
    }

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
