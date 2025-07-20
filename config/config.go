package config

import (
	"encoding/json"
	"os"
	"sync"
)

var (
	config     *Config
	configOnce sync.Once
)

// Config 配置结构体
type Config struct {
	AudioConfigs struct {
		ModelSrc  string `json:"MODEL_SRC"`
		Bytedance struct {
			AudioFile    string `json:"AUDIO_FILE"`
			OutputFile   string `json:"OUTPUT_FILE"`
			AppKey       string `json:"APP_KEY"`
			AccessKey    string `json:"ACCESS_KEY"`
			UID          string `json:"UID"`
			AudioFormat  string `json:"AUDIO_FORMAT"`
			TOSBucket    string `json:"TOS_BUCKET"`
			TOSRegion    string `json:"TOS_REGION"`
			TOSEndpoint  string `json:"TOS_ENDPOINT"`
			TOSAccessKey string `json:"TOS_ACCESS_KEY"`
			TOSSecretKey string `json:"TOS_SECRET_KEY"`
		} `json:"BYTEDANCE"`
	} `json:"AUDIO_CONFIGS"`
	Indexing struct {
		Embedding struct {
			BaseURL string `json:"BASE_URL"`
			APIKey  string `json:"API_KEY"`
			Model   string `json:"MODEL"`
		} `json:"EMBEDDING"`
		Indexer struct {
			Type  string `json:"TYPE"`
			Redis struct {
				Addr      string `json:"ADDR"`
				KeyPrefix string `json:"KEY_PREFIX"`
				BatchSize int    `json:"BATCH_SIZE"`
			} `json:"REDIS"`
			Faiss struct {
				IndexPath string `json:"INDEX_PATH"`
				BatchSize int    `json:"BATCH_SIZE"`
			} `json:"FAISS"`
		} `json:"INDEXER"`
		Loader      map[string]interface{} `json:"LOADER"`
		Transformer struct {
			Headers     map[string]string `json:"HEADERS"`
			TrimHeaders bool              `json:"TRIM_HEADERS"`
		} `json:"TRANSFORMER"`
	} `json:"INDEXING"`
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	configOnce.Do(func() {
		f, err := os.Open("config.json")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		decoder := json.NewDecoder(f)
		cfg := &Config{}
		if err := decoder.Decode(cfg); err != nil {
			panic(err)
		}
		config = cfg
	})
	return config
}
