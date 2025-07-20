package component

import (
	"audio2markdown/config"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cloudwego/eino-ext/components/indexer/redis"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	redisCli "github.com/redis/go-redis/v9"

	redispkg "github.com/cloudwego/eino-examples/quickstart/eino_assistant/pkg/redis"
)

func init() {
	err := redispkg.Init()
	if err != nil {
		log.Fatalf("failed to init redis index: %v", err)
	}
}

// NewIndexer 索引器
func NewIndexer(ctx context.Context) (idr indexer.Indexer, err error) {
	cfg := config.LoadConfig().Indexing.Indexer
	if cfg.Type == "REDIS" {
		redisCfg := cfg.Redis
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
						redispkg.ContentField:  {Value: doc.Content, EmbedKey: redispkg.VectorField},
						redispkg.MetadataField: {Value: metadataBytes},
					},
				}, nil
			},
		}

		embeddingIns11, err := NewEmbedding(ctx)
		if err != nil {
			return nil, err
		}
		config.Embedding = embeddingIns11
		idr, err = redis.NewIndexer(ctx, config)
		if err != nil {
			return nil, err
		}
		return idr, nil
	}
	// TODO: 增加 FAISS 支持
	return nil, fmt.Errorf("unsupported indexer type: %s", cfg.Type)
}
