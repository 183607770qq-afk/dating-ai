package com.datingai.repository;

import com.datingai.model.ConversationSummary;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.transaction.annotation.Transactional;

public interface ConversationSummaryRepository extends JpaRepository<ConversationSummary, Long> {

    ConversationSummary findByUserId(Long userId);

    @Modifying
    @Transactional
    @Query("DELETE FROM ConversationSummary cs WHERE cs.user.id = ?1")
    void deleteByUserId(Long userId);
}
