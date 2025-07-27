package component

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cloudwego/eino/compose"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
)

// TOSUploaderConfig TOS上传配置
type TOSUploaderConfig struct {
	Bucket    string
	Region    string
	Endpoint  string
	AccessKey string
	SecretKey string
}

// NewTOSUploaderLambda 上传文件到 TOS 的 Lambda 节点
func NewTOSUploaderLambda(cfg TOSUploaderConfig) *compose.Lambda {
	client, err := tos.NewClientV2(
		cfg.Endpoint,
		tos.WithRegion(cfg.Region),
		tos.WithCredentials(tos.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey)),
	)
	if err != nil {
		panic(fmt.Sprintf("TOS client init failed: %v", err))
	}
	bucket := cfg.Bucket

	return compose.InvokableLambda(func(ctx context.Context, localPath string) (string, error) {
		file, err := os.Open(localPath)
		if err != nil {
			return "", err
		}
		defer file.Close()

		log.Default().Printf("[NewTOSUploaderLambda] Starting upload for file: %s", localPath)

		objectKey := fmt.Sprintf("audio2md/%s", filepath.Base(localPath))
		_, err = client.PutObjectV2(ctx, &tos.PutObjectV2Input{
			PutObjectBasicInput: tos.PutObjectBasicInput{
				Bucket: bucket,
				Key:    objectKey,
			},
			Content: file,
		})
		if err != nil {
			return "", err
		}
		resp, err := client.PreSignedURL(&tos.PreSignedURLInput{
			Bucket:     bucket,
			Key:        objectKey,
			HTTPMethod: "GET",
			Expires:    3600,
		})
		if err != nil {
			return "", err
		}
		log.Default().Printf("[NewTOSUploaderLambda] File uploaded successfully, object url: %s", resp.SignedUrl)
		return resp.SignedUrl, nil
	})
}
