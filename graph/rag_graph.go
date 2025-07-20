package graph

import (
	component "audio2markdown/component/indexing"
	"context"

	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/compose"
)

// BuildRagGraph 构建 RAG 图
func BuildRagGraph(ctx context.Context) (r compose.Runnable[document.Source, []string], err error) {
	const (
		FileLoaderNode   = "FileLoader"
		MarkdownSplitter = "MarkdownSplitter"
		IndexerNode      = "Indexer"
	)
	g := compose.NewGraph[document.Source, []string]()

	// Loader
	loader, err := component.NewLoader(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddLoaderNode(FileLoaderNode, loader)

	// Transformer
	transformer, err := component.NewDocumentTransformer(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddDocumentTransformerNode(MarkdownSplitter, transformer)

	// Indexer
	indexer, err := component.NewIndexer(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddIndexerNode(IndexerNode, indexer)

	// Edges
	_ = g.AddEdge(compose.START, FileLoaderNode)
	_ = g.AddEdge(FileLoaderNode, MarkdownSplitter)
	_ = g.AddEdge(MarkdownSplitter, IndexerNode)
	_ = g.AddEdge(IndexerNode, compose.END)

	// 编译 Graph
	r, err = g.Compile(ctx, compose.WithGraphName("RAGGraph"), compose.WithNodeTriggerMode(compose.AnyPredecessor))
	if err != nil {
		return nil, err
	}

	return r, nil
}
