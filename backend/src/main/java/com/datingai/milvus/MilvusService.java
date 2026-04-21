package com.datingai.milvus;

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
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import jakarta.annotation.PostConstruct;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

@Service
public class MilvusService {

    @Value("${milvus.host}")
    private String host;

    @Value("${milvus.port}")
    private int port;

    private MilvusServiceClient milvusClient;
    private final String collectionName = "dating_knowledge";

    @Autowired
    private EmbeddingService embeddingService;

    @PostConstruct
    public void init() {
        ConnectParam connectParam = ConnectParam.newBuilder()
                .withHost(host)
                .withPort(port)
                .build();
        milvusClient = new MilvusServiceClient(connectParam);
        createCollectionIfNotExists();
    }

    private void createCollectionIfNotExists() {
        try {
            milvusClient.describeCollection(DescribeCollectionParam.newBuilder()
                    .withCollectionName(collectionName)
                    .build());
        } catch (Exception e) {
            // Collection doesn't exist, create it
            CreateCollectionParam createParam = CreateCollectionParam.newBuilder()
                    .withCollectionName(collectionName)
                    .withShardsNum(2)
                    .build();
            milvusClient.createCollection(createParam);
        }
    }

    public void insertKnowledge(VectorizedKnowledge knowledge) {
        InsertParam insertParam = InsertParam.newBuilder()
                .withCollectionName(collectionName)
                .withFields(List.of(
                       new InsertParam.Field("id", List.of(knowledge.getId().toString())) ,
                                new InsertParam.Field("embedding", List.of((Object) knowledge.getEmbedding())),
                                        new InsertParam.Field( "title", List.of(knowledge.getTitle())),
                                                new InsertParam.Field("content", List.of(knowledge.getContent()))
                ))
                .build();
        milvusClient.insert(insertParam);
    }

    public List<VectorizedKnowledge> searchKnowledge(float[] queryVector, int topK) {
        SearchParam searchParam = SearchParam.newBuilder()
                .withCollectionName(collectionName)
                .withVectorFieldName("embedding")
                .withVectors(List.of(queryVector))
                .withTopK(topK)
                .withMetricType(MetricType.L2)
                .withParams("{\"nprobe\": 10}")
                .build();

        R<SearchResults> searchResponse = milvusClient.search(searchParam);
        List<VectorizedKnowledge> results = new ArrayList<>();

        if (searchResponse.getStatus() == R.Status.Success.getCode()) {
            SearchResultsWrapper wrapper = new SearchResults(searchResponse.getData());
            for (int i = 0; i < wrapper.getResultCount(); i++) {
                VectorizedKnowledge knowledge = new VectorizedKnowledge();
                knowledge.setId(Long.parseLong(wrapper.getFieldData("id", i).toString()));
                knowledge.setTitle(wrapper.getFieldData("title", i).toString());
                knowledge.setContent(wrapper.getFieldData("content", i).toString());
                results.add(knowledge);
            }
        }
        return results;
    }

    public float[] getEmbedding(String text) throws IOException {
        return embeddingService.getEmbedding(text);
    }
}
