package com.datingai.service;

import com.datingai.model.Knowledge;
import com.datingai.repository.KnowledgeRepository;
import com.datingai.milvus.MilvusService;
import com.datingai.milvus.VectorizedKnowledge;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Optional;

@Service
public class KnowledgeService {

    @Autowired
    private KnowledgeRepository knowledgeRepository;

    @Autowired
    private MilvusService milvusService;

    public Knowledge createKnowledge(Knowledge knowledge) {
        Knowledge savedKnowledge = knowledgeRepository.save(knowledge);
        // 向Milvus中插入向量化的知识
        // 这里需要实现文本向量化的逻辑
         VectorizedKnowledge vectorizedKnowledge = new VectorizedKnowledge();
         vectorizedKnowledge.setId(savedKnowledge.getId());
         vectorizedKnowledge.setTitle(savedKnowledge.getTitle());
         vectorizedKnowledge.setContent(savedKnowledge.getContent());
         vectorizedKnowledge.setEmbedding(generateEmbedding(savedKnowledge.getContent()));
         milvusService.insertKnowledge(vectorizedKnowledge);
        return savedKnowledge;
    }

    public List<Knowledge> getAllPublishedKnowledge() {
        return knowledgeRepository.findByIsPublishedTrue();
    }

    public List<Knowledge> getKnowledgeByCategory(String category) {
        return knowledgeRepository.findByCategoryAndIsPublishedTrue(category);
    }

    public Optional<Knowledge> getKnowledgeById(Long id) {
        return knowledgeRepository.findById(id);
    }

    public Knowledge updateKnowledge(Knowledge knowledge) {
        return knowledgeRepository.save(knowledge);
    }

    public void deleteKnowledge(Long id) {
        knowledgeRepository.deleteById(id);
    }

    // 生成文本嵌入向量的方法，需要使用实际的嵌入模型
    private float[] generateEmbedding(String text) {
        // 这里需要实现文本向量化的逻辑
        return new float[128]; // 假设使用128维向量
    }
}
