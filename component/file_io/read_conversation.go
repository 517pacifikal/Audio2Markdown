package fileio

import (
	"context"
	"os"
	"regexp"
	"strings"

	"github.com/cloudwego/eino/schema"
)

// NewReadConversationLambda 解析对话文本为 []*schema.Message，支持任意数量的角色
func NewReadConversationLambda(ctx context.Context, filePath string) ([]*schema.Message, error) {
	bs, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(bs), "\n")
	var messages []*schema.Message
	re := regexp.MustCompile(`^(\d+):\s*(.*)$`)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		matches := re.FindStringSubmatch(line)
		if len(matches) != 3 {
			// 跳过无法识别的行
			continue
		}
		roleNum := matches[1]
		content := matches[2]
		role := schema.RoleType("role_" + roleNum)
		messages = append(messages, &schema.Message{
			Role:    role,
			Content: strings.TrimSpace(content),
		})
	}
	return messages, nil
}
