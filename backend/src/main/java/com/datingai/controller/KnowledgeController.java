package com.datingai.controller;

import com.datingai.model.Knowledge;
import com.datingai.service.KnowledgeService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/knowledge")
public class KnowledgeController {

    @Autowired
    private KnowledgeService knowledgeService;

    @GetMapping("/all")
    public ResponseEntity<List<Knowledge>> getAllKnowledge() {
        List<Knowledge> knowledgeList = knowledgeService.getAllPublishedKnowledge();
        return ResponseEntity.ok(knowledgeList);
    }

    @GetMapping("/category/{category}")
    public ResponseEntity<List<Knowledge>> getKnowledgeByCategory(@PathVariable String category) {
        List<Knowledge> knowledgeList = knowledgeService.getKnowledgeByCategory(category);
        return ResponseEntity.ok(knowledgeList);
    }

    @GetMapping("/{id}")
    public ResponseEntity<Knowledge> getKnowledgeById(@PathVariable Long id) {
        return knowledgeService.getKnowledgeById(id)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }

    @PostMapping("/create")
    public ResponseEntity<Knowledge> createKnowledge(@RequestBody Knowledge knowledge) {
        Knowledge createdKnowledge = knowledgeService.createKnowledge(knowledge);
        return ResponseEntity.ok(createdKnowledge);
    }

    @PutMapping("/{id}")
    public ResponseEntity<Knowledge> updateKnowledge(@PathVariable Long id, @RequestBody Knowledge knowledge) {
        knowledge.setId(id);
        Knowledge updatedKnowledge = knowledgeService.updateKnowledge(knowledge);
        return ResponseEntity.ok(updatedKnowledge);
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteKnowledge(@PathVariable Long id) {
        knowledgeService.deleteKnowledge(id);
        return ResponseEntity.noContent().build();
    }
}
