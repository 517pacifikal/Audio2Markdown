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
	AudioConfigs AudioConfigs   `json:"AUDIO_CONFIGS"`
	Indexing     IndexingConfig `json:"INDEXING"`
}

type AudioConfigs struct {
	ModelSrc  string          `json:"MODEL_SRC"`
	Bytedance BytedanceConfig `json:"BYTEDANCE"`
}

type IndexerConfig struct {
	Type  string      `json:"TYPE"`
	Redis RedisConfig `json:"REDIS"`
	Faiss FaissConfig `json:"FAISS"`
}

type RedisConfig struct {
	Addr      string `json:"ADDR"`
	KeyPrefix string `json:"KEY_PREFIX"`
	BatchSize int    `json:"BATCH_SIZE"`
}

type FaissConfig struct {
	IndexPath string `json:"INDEX_PATH"`
	BatchSize int    `json:"BATCH_SIZE"`
}

type EmbeddingConfig struct {
	BaseURL string `json:"BASE_URL"`
	APIKey  string `json:"API_KEY"`
	Model   string `json:"MODEL"`
}

type TransformerConfig struct {
	Headers     map[string]string `json:"HEADERS"`
	TrimHeaders bool              `json:"TRIM_HEADERS"`
}

type IndexingConfig struct {
	FilePath    string                 `json:"FILE_PATH"`
	Embedding   EmbeddingConfig        `json:"EMBEDDING"`
	Indexer     IndexerConfig          `json:"INDEXER"`
	Loader      map[string]interface{} `json:"LOADER"`
	Transformer TransformerConfig      `json:"TRANSFORMER"`
}

type BytedanceConfig struct {
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
