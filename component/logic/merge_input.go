package logic

import "context"

// MergeInputAndDocs 合并输入
func NewMergeInputNode(ctx context.Context, in map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		"read_input":     in["read_input"],
		"knowledge_docs": in["knowledge_docs"],
	}, nil
}
