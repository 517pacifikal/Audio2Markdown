package component

import (
	"audio2markdown/config"
	"context"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino/components/document"
)

// NewDocumentTransformer 对加载的文档进行分段/分词处理
func NewDocumentTransformer(ctx context.Context) (tfr document.Transformer, err error) {
	cfg := config.LoadConfig().Indexing.Transformer
	config := &markdown.HeaderConfig{
		Headers:     cfg.Headers,
		TrimHeaders: cfg.TrimHeaders,
	}
	tfr, err = markdown.NewHeaderSplitter(ctx, config)
	if err != nil {
		return nil, err
	}
	return tfr, nil
}
