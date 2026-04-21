package com.datingai.repository;

import com.datingai.model.Subscription;
import com.datingai.model.User;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface SubscriptionRepository extends JpaRepository<Subscription, Long> {
    List<Subscription> findByUser(User user);
    Subscription findByUserAndStatus(User user, String status);
}
