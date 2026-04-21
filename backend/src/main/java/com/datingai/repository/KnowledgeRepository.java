package com.datingai.repository;

import com.datingai.model.Knowledge;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface KnowledgeRepository extends JpaRepository<Knowledge, Long> {
    List<Knowledge> findByIsPublishedTrue();
    List<Knowledge> findByCategoryAndIsPublishedTrue(String category);
}
