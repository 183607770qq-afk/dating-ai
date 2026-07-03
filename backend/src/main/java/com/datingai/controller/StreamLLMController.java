package com.datingai.controller;

import com.datingai.llm.StreamLLMService;
import com.datingai.model.User;
import com.datingai.repository.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.MediaType;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.servlet.mvc.method.annotation.SseEmitter;

import java.util.Map;

@RestController
@RequestMapping("/api/llm/stream")
public class StreamLLMController {

    @Autowired
    private StreamLLMService streamLLMService;

    @Autowired
    private UserRepository userRepository;

    @PostMapping(value = "/advice", produces = MediaType.TEXT_EVENT_STREAM_VALUE)
    public SseEmitter getStreamingAdvice(@RequestBody Map<String, String> request) {
        String question = request.get("question");
        Long userId = extractUserId();
        return streamLLMService.getStreamingAdvice(question, userId);
    }

    /**
     * 从 SecurityContext 提取当前登录用户的 ID，未登录则返回 null
     */
    private Long extractUserId() {
        Authentication auth = SecurityContextHolder.getContext().getAuthentication();
        if (auth == null || !auth.isAuthenticated()) {
            return null;
        }
        String username = auth.getName();
        if (username == null || "anonymousUser".equals(username)) {
            return null;
        }
        User user = userRepository.findByUsername(username);
        return user != null ? user.getId() : null;
    }
}
