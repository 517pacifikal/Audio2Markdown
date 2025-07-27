package component

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cloudwego/eino/compose"
	"github.com/google/uuid"
)

type ASRConfig struct {
	AppKey     string
	AccessKey  string
	UID        string
	Format     string
	OutputFile string
}

// NewASRLambda 用于语音识别的 Lambda 函数
func NewASRLambda(cfg ASRConfig) *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, audioURL string) (string, error) {
		log.Default().Printf("[NewASRLambda] Starting ASR for audio URL: %s", audioURL)
		requestID, err := submitASRTask(ctx, cfg, audioURL)
		if err != nil {
			return "", err
		}
		log.Default().Printf("[NewASRLambda] Task submitted successfully, request ID: %s", requestID)
		result, err := queryASRResultRaw(ctx, cfg, requestID)
		if err != nil {
			return "", err
		}

		output := formatASRResultWithSpeaker(result)

		log.Default().Printf("[NewASRLambda] Writing ASR result to path: %s", cfg.OutputFile)

		// 导出到指定文件
		if cfg.OutputFile != "" {
			if err := os.MkdirAll(getDir(cfg.OutputFile), 0755); err != nil {
				return "", fmt.Errorf("failed to create output dir: %v", err)
			}
			if err := os.WriteFile(cfg.OutputFile, []byte(output), 0644); err != nil {
				return "", fmt.Errorf("failed to write output file: %v", err)
			}
			log.Default().Printf("[NewASRLambda] Output written to %s", cfg.OutputFile)
		} else {
			log.Default().Printf("[NewASRLambda] Empty output, nothing written to file!")
			return "", fmt.Errorf("output file path is empty")
		}
		return cfg.OutputFile, nil
	})
}

// submitASRTask 提交ASR任务，返回requestID
func submitASRTask(ctx context.Context, cfg ASRConfig, audioURL string) (string, error) {
	submitURL := "https://openspeech.bytedance.com/api/v3/auc/bigmodel/submit"
	requestID := uuid.New().String()
	headers := map[string]string{
		"Content-Type":      "application/json",
		"X-Api-App-Key":     cfg.AppKey,
		"X-Api-Access-Key":  cfg.AccessKey,
		"X-Api-Resource-Id": "volc.bigasr.auc",
		"X-Api-Request-Id":  requestID,
		"X-Api-Sequence":    "-1",
	}
	payload := map[string]interface{}{
		"user": map[string]interface{}{
			"uid": cfg.UID,
		},
		"audio": map[string]interface{}{
			"format": cfg.Format,
			"url":    audioURL,
		},
		"request": map[string]interface{}{
			"model_name":          "bigmodel",
			"enable_itn":          true,
			"enable_punc":         true,
			"show_utterances":     true,
			"enable_speaker_info": true, // 识别说话人
		},
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", submitURL, bytes.NewReader(body))
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.Header.Get("X-Api-Status-Code") != "20000000" {
		return "", fmt.Errorf("submit failed: %s", resp.Header.Get("X-Api-Message"))
	}
	return requestID, nil
}

// queryASRResultRaw 查询ASR任务结果，返回完整JSON
func queryASRResultRaw(ctx context.Context, cfg ASRConfig, requestID string) (map[string]interface{}, error) {
	queryURL := "https://openspeech.bytedance.com/api/v3/auc/bigmodel/query"
	headers := map[string]string{
		"Content-Type":      "application/json",
		"X-Api-App-Key":     cfg.AppKey,
		"X-Api-Access-Key":  cfg.AccessKey,
		"X-Api-Resource-Id": "volc.bigasr.auc",
		"X-Api-Request-Id":  requestID,
		"X-Api-Sequence":    "-1",
	}
	maxTries := 180
	interval := 10 * time.Second
	for i := 0; i < maxTries; i++ {
		req, _ := http.NewRequestWithContext(ctx, "POST", queryURL, bytes.NewReader([]byte("{}")))
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		statusCode := resp.Header.Get("X-Api-Status-Code")
		if statusCode == "20000000" {
			var result map[string]interface{}
			body, _ := io.ReadAll(resp.Body)
			json.Unmarshal(body, &result)
			return result, nil
		} else if statusCode == "20000001" || statusCode == "20000002" {
			time.Sleep(interval)
			continue
		} else {
			return nil, fmt.Errorf("query failed: %s", resp.Header.Get("X-Api-Message"))
		}
	}
	return nil, fmt.Errorf("query timeout")
}

// formatASRResultWithSpeaker 格式化输出文本（带说话人信息）
func formatASRResultWithSpeaker(result map[string]interface{}) string {
	var output string
	if res, ok := result["result"].(map[string]interface{}); ok {
		if utterances, ok := res["utterances"].([]interface{}); ok {
			for _, u := range utterances {
				if utt, ok := u.(map[string]interface{}); ok {
					speaker := ""
					if additions, ok := utt["additions"].(map[string]interface{}); ok {
						if s, ok := additions["speaker"].(float64); ok {
							speaker = fmt.Sprintf("Speaker%d", int(s))
						} else if s, ok := additions["speaker"].(string); ok {
							speaker = s
						}
					}
					if speaker == "" {
						if s, ok := utt["speaker_id"].(string); ok {
							speaker = s
						} else if s, ok := utt["speaker_id"].(float64); ok {
							speaker = fmt.Sprintf("Speaker%d", int(s))
						}
					}
					text := utt["text"]
					if speaker != "" {
						output += fmt.Sprintf("%s: %s\n", speaker, text)
					} else {
						output += fmt.Sprintf("%s\n", text)
					}
				}
			}
		} else if text, ok := res["text"].(string); ok {
			output = text
		}
	}
	return output
}

// getDir 获取文件所在目录
func getDir(path string) string {
	if idx := len(path) - 1; idx >= 0 && path[idx] == '/' {
		return path
	}
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[:i]
		}
	}
	return "."
}
