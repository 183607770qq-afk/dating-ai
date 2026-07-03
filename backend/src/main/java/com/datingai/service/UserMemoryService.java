package com.datingai.service;

import com.datingai.llm.EmbeddingService;
import io.milvus.client.MilvusServiceClient;
import io.milvus.grpc.SearchResults;
import io.milvus.param.ConnectParam;
import io.milvus.param.MetricType;
import io.milvus.param.R;
import io.milvus.param.collection.CreateCollectionParam;
import io.milvus.param.collection.DescribeCollectionParam;
import io.milvus.param.dml.InsertParam;
import io.milvus.param.dml.SearchParam;
import io.milvus.response.SearchResultsWrapper;
import jakarta.annotation.PostConstruct;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

/**
 * 用户记忆向量服务 — 基于 Milvus 存储用户关键事实的向量，支持语义检索。
 * 与 conversation_summaries（全局摘要）互补：摘要管画像，Milvus 管细节。
 */
@Service
public class UserMemoryService {

    private static final Logger logger = LoggerFactory.getLogger(UserMemoryService.class);

    private static final String COLLECTION_NAME = "user_memory";
    private static final int EMBED_DIM = 1024; // Ollama qwen3-vl:4b 默认维度

    @Value("${milvus.host}")
    private String host;

    @Value("${milvus.port}")
    private int port;

    @Autowired
    private EmbeddingService embeddingService;

    private MilvusServiceClient milvusClient;

    @PostConstruct
    public void init() {
        ConnectParam connectParam = ConnectParam.newBuilder()
                .withHost(host)
                .withPort(port)
                .build();
        milvusClient = new MilvusServiceClient(connectParam);
        createCollectionIfNotExists();
        logger.info("UserMemoryService initialized, collection={}", COLLECTION_NAME);
    }

    private void createCollectionIfNotExists() {
        try {
            milvusClient.describeCollection(DescribeCollectionParam.newBuilder()
                    .withCollectionName(COLLECTION_NAME)
                    .build());
            logger.info("Collection {} already exists", COLLECTION_NAME);
        } catch (Exception e) {
            CreateCollectionParam createParam = CreateCollectionParam.newBuilder()
                    .withCollectionName(COLLECTION_NAME)
                    .withShardsNum(1)
                    .build();
            milvusClient.createCollection(createParam);
            logger.info("Created collection: {}", COLLECTION_NAME);
        }
    }

    /**
     * 批量存储用户的关键事实向量。
     * 每次摘要更新时，先删旧数据再插入最新的事实（避免重复和过期信息）。
     *
     * @param userId 用户 ID
     * @param facts  从摘要中提取的关键事实列表
     */
    public void storeUserFacts(Long userId, List<String> facts) {
        if (facts == null || facts.isEmpty()) return;

        // 先删除该用户旧的记忆向量
        deleteByUserId(userId);

        long timestamp = System.currentTimeMillis();
        for (String fact : facts) {
            try {
                float[] embedding = embeddingService.getEmbedding(fact);
                String id = "umem_" + userId + "_" + UUID.randomUUID().toString().substring(0, 8);

                InsertParam insertParam = InsertParam.newBuilder()
                        .withCollectionName(COLLECTION_NAME)
                        .withFields(List.of(
                                new InsertParam.Field("id", List.of(id)),
                                new InsertParam.Field("user_id", List.of(String.valueOf(userId))),
                                new InsertParam.Field("embedding", List.of((Object) embedding)),
                                new InsertParam.Field("fact", List.of(fact))
                        ))
                        .build();

                milvusClient.insert(insertParam);
                logger.debug("Stored user fact: userId={}, fact={}", userId,
                        fact.length() > 50 ? fact.substring(0, 50) + "..." : fact);
            } catch (Exception e) {
                logger.error("Failed to store user fact for userId={}: {}", userId, e.getMessage());
            }
        }

        // flush 确保数据可立即检索
        milvusClient.flush(io.milvus.param.collection.FlushParam.newBuilder()
                .addCollectionName(COLLECTION_NAME)
                .build());
    }

    /**
     * 语义检索用户的相关记忆。
     * 对当前查询文本做 embedding，然后搜用户记忆中语义最接近的 topK 条事实。
     *
     * @param userId    用户 ID
     * @param queryText 当前查询文本（用于语义匹配）
     * @param topK      返回条数
     * @return 相关事实文本列表，按相似度降序
     */
    public List<String> searchUserMemories(Long userId, String queryText, int topK) {
        try {
            float[] queryVector = embeddingService.getEmbedding(queryText);

            SearchParam searchParam = SearchParam.newBuilder()
                    .withCollectionName(COLLECTION_NAME)
                    .withVectorFieldName("embedding")
                    .withVectors(List.of(queryVector))
                    .withTopK(topK)
                    .withMetricType(MetricType.L2)
                    .withParams("{\"nprobe\": 10}")
                    .withExpr("user_id == \"" + userId + "\"")  // 只搜该用户的记忆
                    .build();

            R<SearchResults> searchResponse = milvusClient.search(searchParam);
            List<String> results = new ArrayList<>();

            if (searchResponse.getStatus() == R.Status.Success.getCode()) {
                SearchResultsWrapper wrapper = new SearchResultsWrapper(
                        searchResponse.getData().getResults());
                for (int i = 0; i < wrapper.getRowRecords().size(); i++) {
                    Object factObj = wrapper.getFieldData("fact", i);
                    if (factObj != null) {
                        results.add(factObj.toString());
                    }
                }
            }

            logger.debug("User memory search: userId={}, query='{}', found {} results",
                    userId, queryText.length() > 30 ? queryText.substring(0, 30) + "..." : queryText,
                    results.size());
            return results;
        } catch (Exception e) {
            logger.warn("Failed to search user memories for userId={}: {}", userId, e.getMessage());
            return List.of();
        }
    }

    /**
     * 删除某用户的所有记忆向量（用于摘要更新时刷新）
     */
    private void deleteByUserId(Long userId) {
        try {
            milvusClient.delete(io.milvus.param.dml.DeleteParam.newBuilder()
                    .withCollectionName(COLLECTION_NAME)
                    .withExpr("user_id == \"" + userId + "\"")
                    .build());
        } catch (Exception e) {
            logger.warn("Failed to delete old memories for userId={}: {}", userId, e.getMessage());
        }
    }
}
