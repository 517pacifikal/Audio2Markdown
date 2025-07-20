package main

import (
	"audio2markdown/config"
	"audio2markdown/graph"
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino/components/document"
)

// runRagGraph 构建并运行 RAG 图
func runRagGraph() {
	// 构建 RAG 子图
	ragGraph, err := graph.BuildRagGraph(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build RAG graph: %v\n", err)
		os.Exit(1)
	}

	// 构造输入
	src := document.Source{
		URI: config.LoadConfig().Indexing.FilePath,
	}

	// 执行图
	result, err := ragGraph.Invoke(context.Background(), src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run RAG graph: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("RAG Result: %+v\n", result)
}

// runAudio2TextGraph 构建并运行音频转文本图
func runA2TGraph() {

	// 构建音频转文本子图
	audioGraph, err := graph.BuildAudio2TextGraph(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build audio to text graph: %v\n", err)
		os.Exit(1)
	}

	// 构造输入
	audioFilePath := config.LoadConfig().AudioConfigs.Bytedance.AudioFile

	// 执行图
	text, err := audioGraph.Invoke(context.Background(), audioFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run audio to text graph: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Transcribed Text: %s\n", text)
}

func main() {
	runA2TGraph()
}
