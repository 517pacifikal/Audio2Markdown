package chatmodel

import (
	"audio2markdown/config"
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
)

// NewOpenAIChatModel 创建LLM客户端
func NewOpenAIChatModel(ctx context.Context) model.ToolCallingChatModel {
	key := config.LoadConfig().ChatModelConfig.OpenAIAPIKey
	modelName := config.LoadConfig().ChatModelConfig.OpenAIModel
	baseURL := config.LoadConfig().ChatModelConfig.OpenAIBaseURL
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: baseURL,
		Model:   modelName,
		APIKey:  key,
	})
	if err != nil {
		log.Fatalf("create openai chat model failed, err=%v", err)
	}
	return chatModel
}
