package graph

import (
	component "audio2markdown/component/audio"
	"audio2markdown/config"
	"context"

	"github.com/cloudwego/eino/compose"
)

// BuildAudio2TextGraph 构建音频转文本的Graph
func BuildAudio2TextGraph(ctx context.Context) (compose.Runnable[string, string], error) {
	const (
		TOSUploaderNode = "TOSUploader"
		ASRNode         = "ASR"
	)

	g := compose.NewGraph[string, string]()

	tosCfg := component.TOSUploaderConfig{
		Endpoint:  config.LoadConfig().AudioConfigs.Bytedance.TOSEndpoint,
		AccessKey: config.LoadConfig().AudioConfigs.Bytedance.TOSAccessKey,
		SecretKey: config.LoadConfig().AudioConfigs.Bytedance.TOSSecretKey,
		Bucket:    config.LoadConfig().AudioConfigs.Bytedance.TOSBucket,
		Region:    config.LoadConfig().AudioConfigs.Bytedance.TOSRegion,
	}

	asrCfg := component.ASRConfig{
		AppKey:     config.LoadConfig().AudioConfigs.Bytedance.AppKey,
		AccessKey:  config.LoadConfig().AudioConfigs.Bytedance.AccessKey,
		UID:        config.LoadConfig().AudioConfigs.Bytedance.UID,
		OutputFile: config.LoadConfig().AudioConfigs.Bytedance.OutputFile,
	}

	// TOS上传节点
	tosUploader := component.NewTOSUploaderLambda(tosCfg)
	if err := g.AddLambdaNode(TOSUploaderNode, tosUploader); err != nil {
		return nil, err
	}

	// ASR节点
	asr := component.NewASRLambda(asrCfg)
	if err := g.AddLambdaNode(ASRNode, asr); err != nil {
		return nil, err
	}

	// Edges
	_ = g.AddEdge(compose.START, TOSUploaderNode)
	_ = g.AddEdge(TOSUploaderNode, ASRNode)
	_ = g.AddEdge(ASRNode, compose.END)

	// 编译 Graph
	r, err := g.Compile(ctx, compose.WithGraphName("Audio2TextGraph"))
	if err != nil {
		return nil, err
	}
	return r, nil
}
