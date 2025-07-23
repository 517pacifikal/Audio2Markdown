package retriever

import (
	"audio2markdown/component/indexing"
	"audio2markdown/config"
	"context"

	"github.com/cloudwego/eino-ext/components/retriever/redis"
	"github.com/cloudwego/eino/components/retriever"
	redisCli "github.com/redis/go-redis/v9"
)

// NewRedisRetriever 创建 Redis 检索器
func NewRedisRetriever(ctx context.Context) (retriever.Retriever, error) {
	redisCfg := config.LoadConfig().Indexing.Indexer.Redis
	client := redisCli.NewClient(&redisCli.Options{
		Addr:     redisCfg.Addr,
		Protocol: 2,
	})
	embedder, err := indexing.NewEmbedding(ctx)
	if err != nil {
		return nil, err
	}
	return redis.NewRetriever(ctx, &redis.RetrieverConfig{
		Client:    client,
		Index:     "rag",
		TopK:      redisCfg.BatchSize,
		Embedding: embedder,
	})
}
