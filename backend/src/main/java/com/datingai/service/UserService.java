package com.datingai.service;

import com.datingai.model.User;
import com.datingai.model.Subscription;
import com.datingai.repository.UserRepository;
import com.datingai.repository.SubscriptionRepository;
import com.datingai.security.JwtUtil;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.Optional;

@Service
public class UserService {

    @Autowired
    private UserRepository userRepository;

    @Autowired
    private SubscriptionRepository subscriptionRepository;

    @Autowired
    private PasswordEncoder passwordEncoder;

    @Autowired
    private JwtUtil jwtUtil;

    public User register(User user) {
        user.setPassword(passwordEncoder.encode(user.getPassword()));
        user.setRole("USER");
        user.setSubscribed(false);
        return userRepository.save(user);
    }

    public String login(String username, String password) {
        User user = userRepository.findByUsername(username);
        if (user == null || !passwordEncoder.matches(password, user.getPassword())) {
            throw new RuntimeException("Invalid username or password");
        }
        return jwtUtil.generateToken(username);
    }

    public User subscribe(User user, String subscriptionType) {
        LocalDateTime now = LocalDateTime.now();
        LocalDateTime endDate = now.plusMonths(1); // 默认为1个月订阅

        Subscription subscription = new Subscription();
        subscription.setUser(user);
        subscription.setSubscriptionType(subscriptionType);
        subscription.setStartDate(now);
        subscription.setEndDate(endDate);
        subscription.setStatus("ACTIVE");
        subscriptionRepository.save(subscription);

        user.setSubscribed(true);
        user.setSubscriptionEndDate(endDate);
        user.setRole("SUBSCRIBER");
        return userRepository.save(user);
    }

    public Optional<User> findByUsername(String username) {
        return Optional.ofNullable(userRepository.findByUsername(username));
    }

    public User updateUser(User user) {
        return userRepository.save(user);
    }
}
