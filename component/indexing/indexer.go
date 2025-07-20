package component

import (
	"audio2markdown/config"
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino-ext/components/indexer/redis"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	redisCli "github.com/redis/go-redis/v9"
)

// Redis索引器
func newRedisIndexer(ctx context.Context, redisCfg config.RedisConfig) (indexer.Indexer, error) {
	redisClient := redisCli.NewClient(&redisCli.Options{
		Addr:     redisCfg.Addr,
		Protocol: 2,
	})

	config := &redis.IndexerConfig{
		Client:    redisClient,
		KeyPrefix: redisCfg.KeyPrefix,
		BatchSize: redisCfg.BatchSize,
		DocumentToHashes: func(ctx context.Context, doc *schema.Document) (*redis.Hashes, error) {
			if doc.ID == "" {
				doc.ID = uuid.New().String()
			}
			key := doc.ID

			metadataBytes, err := json.Marshal(doc.MetaData)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal metadata: %+v", err)
			}

			return &redis.Hashes{
				Key: key,
				Field2Value: map[string]redis.FieldValue{
					"content":  {Value: doc.Content, EmbedKey: "vector"},
					"metadata": {Value: metadataBytes},
				},
			}, nil
		},
	}

	embeddingIns, err := NewEmbedding(ctx)
	if err != nil {
		return nil, err
	}
	config.Embedding = embeddingIns
	return redis.NewIndexer(ctx, config)
}

// FAISS索引器
func newFaissIndexer(ctx context.Context, faissCfg config.FaissConfig) (indexer.Indexer, error) {
	// TODO
	return nil, fmt.Errorf("FAISS indexer is not implemented yet")
}

// NewIndexer 索引器
func NewIndexer(ctx context.Context) (indexer.Indexer, error) {
	cfg := config.LoadConfig().Indexing.Indexer
	switch cfg.Type {
	case "REDIS":
		return newRedisIndexer(ctx, cfg.Redis)
	case "FAISS":
		return newFaissIndexer(ctx, cfg.Faiss)
	default:
		return nil, fmt.Errorf("unsupported indexer type: %s", cfg.Type)
	}
}
