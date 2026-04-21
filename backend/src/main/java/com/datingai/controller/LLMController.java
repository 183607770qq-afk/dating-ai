package com.datingai.controller;

import com.datingai.llm.LLMService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.io.IOException;
import java.util.Map;

@RestController
@RequestMapping("/api/llm")
public class LLMController {

    @Autowired
    private LLMService llmService;

    @PostMapping("/advice")
    public ResponseEntity<Map<String, String>> getDatingAdvice(@RequestBody Map<String, String> request) {
        try {
            String question = request.get("question");
            String advice = llmService.getDatingAdvice(question);
            return ResponseEntity.ok(Map.of("advice", advice));
        } catch (IOException e) {
            return ResponseEntity.status(500).body(Map.of("error", "Failed to get advice"));
        }
    }
}
