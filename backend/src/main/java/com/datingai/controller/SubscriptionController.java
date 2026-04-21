package com.datingai.controller;

import com.datingai.model.User;
import com.datingai.service.UserService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.web.bind.annotation.*;

import java.util.Map;

@RestController
@RequestMapping("/api/subscription")
public class SubscriptionController {

    @Autowired
    private UserService userService;

    @PostMapping("/subscribe")
    public ResponseEntity<User> subscribe(Authentication authentication, @RequestBody Map<String, String> request) {
        UserDetails userDetails = (UserDetails) authentication.getPrincipal();
        User user = userService.findByUsername(userDetails.getUsername())
                .orElseThrow(() -> new RuntimeException("User not found"));
        String subscriptionType = request.get("subscriptionType");
        User updatedUser = userService.subscribe(user, subscriptionType);
        return ResponseEntity.ok(updatedUser);
    }

    @GetMapping("/status")
    public ResponseEntity<Map<String, Object>> getSubscriptionStatus(Authentication authentication) {
        UserDetails userDetails = (UserDetails) authentication.getPrincipal();
        User user = userService.findByUsername(userDetails.getUsername())
                .orElseThrow(() -> new RuntimeException("User not found"));
        return ResponseEntity.ok(Map.of(
                "isSubscribed", user.isSubscribed(),
                "subscriptionEndDate", user.getSubscriptionEndDate()
        ));
    }
}
