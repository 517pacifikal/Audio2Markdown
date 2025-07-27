package logic

import (
	"context"
	"strings"

	"github.com/cloudwego/eino/schema"
)

// MessagesToQueryLambda 将消息列表转换为查询字符串
func MessagesToQueryLambda(ctx context.Context, messages []*schema.Message) (string, error) {
	var sb strings.Builder
	for _, m := range messages {
		sb.WriteString(m.Content)
		sb.WriteString("\n")
	}
	return sb.String(), nil
}
