package fileio

import (
	"audio2markdown/config"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/cloudwego/eino/schema"
)

// NewWriteMarkdownLambda 将 Markdown 内容写入指定目录
func NewWriteMarkdownLambda(ctx context.Context, msg *schema.Message) (string, error) {
	mdContent := msg.Content
	outputDir := config.LoadConfig().Common.OutputDir
	log.Default().Printf("[NewWriteMarkdownLambda] LLM message: %+v", msg)

	_ = os.MkdirAll(outputDir, 0755)
	filename := fmt.Sprintf("summary_%d.md", time.Now().UnixNano())
	path := filepath.Join(outputDir, filename)
	err := os.WriteFile(path, []byte(mdContent), 0644)
	if err != nil {
		return "", err
	}
	return path, nil
}
