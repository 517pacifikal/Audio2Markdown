package component

import (
	"audio2markdown/config"
	"context"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino/components/embedding"
)

// NewEmbedding 嵌入模型
func NewEmbedding(ctx context.Context) (eb embedding.Embedder, err error) {
	cfg := config.LoadConfig().Indexing.Embedding
	config := &ark.EmbeddingConfig{
		BaseURL: cfg.BaseURL,
		APIKey:  cfg.APIKey,
		Model:   cfg.Model,
	}
	eb, err = ark.NewEmbedder(ctx, config)
	if err != nil {
		return nil, err
	}
	return eb, nil
}
