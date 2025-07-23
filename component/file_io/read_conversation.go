package fileio

import (
	"context"
	"os"
)

// NewReadConversationLambda 读入对话内容Lambda（输入为单个文件路径）
func NewReadConversationLambda(ctx context.Context, filePath string) (map[string]interface{}, error) {
	bs, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"read_input": string(bs)}, nil
}
