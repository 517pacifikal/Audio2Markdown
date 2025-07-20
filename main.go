package main

import (
	"audio2markdown/graph"
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino/components/document"
)

func main() {

	// 构建 RAG 子图
	ragGraph, err := graph.BuildRagGraph(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build RAG graph: %v\n", err)
		os.Exit(1)
	}

	// 构造输入（如指定目录路径）
	src := document.Source{
		URI: "./your_input_dir", // 替换为你的目录路径
	}

	// 执行图
	result, err := ragGraph.Invoke(context.Background(), src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run RAG graph: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("RAG Result: %+v\n", result)
}
