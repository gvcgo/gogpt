package config

import (
	"os"
	"path/filepath"

	"github.com/moqsien/goutils/pkgs/gutils"
	"github.com/moqsien/goutils/pkgs/koanfer"
	"github.com/sashabaranov/go-openai"
)

const (
	ConfigFileName string = "gogpt_conf.json"
	PromptUrl      string = "https://gitlab.com/moqsien/gpt_resources/-/raw/main/prompt.json"
)

// ChatGPT
type OpenAIConf struct {
	BaseUrl            string         `koanf,json:"base_url"`
	ApiKey             string         `koanf,json:"api_key"`
	ApiType            openai.APIType `koanf,json:"api_type"`
	ApiVersion         string         `koanf,json:"api_version"`
	OrgID              string         `koanf,json:"org_id"`
	Engine             string         `koanf,json:"engine"`
	EmptyMessagesLimit uint           `koanf,json:"empty_msg_limit"`
	Proxy              string         `koanf,json:"proxy"`
	Model              string         `koanf,json:"model"`
	MaxTokens          int            `koanf,json:"max_tokens"`
	ContextLen         int            `koanf,json:"context_length"`
	Temperature        float32        `koanf,json:"temperature"`
	PromptMsgUrl       string         `koanf,json:"prompt_msgs_url"`
	PromptStr          string         `koanf,json:"prompt"`
}

type SparkAPIVersion string

const (
	SparkAPIV1Dot1 string          = "wss://spark-api.xf-yun.com/v1.1/chat"
	SparkAPIV2Dot1 string          = "wss://spark-api.xf-yun.com/v2.1/chat"
	SparkAPIV3Dot1 string          = "wss://spark-api.xf-yun.com/v3.1/chat"
	SparkDomainV1  string          = "general"
	SparkDomainV2  string          = "general2"
	SparkDomainV3  string          = "general3"
	SparkAPIV1     SparkAPIVersion = "v1.1"
	SparkAPIV2     SparkAPIVersion = "v2.1"
	SparkAPIV3     SparkAPIVersion = "v3.1"
)

// IFlyTek Spark
type IflySparkConf struct {
	APIVersion  SparkAPIVersion `koanf,json:"spark_api_version"`
	APPID       string          `koanf,json:"spark_app_id"`
	UID         string          `koanf,json:"spark_user_id"`
	APPKey      string          `koanf,json:"spark_app_key"`
	APPSecrete  string          `koanf,json:"spark_app_secrete"`
	MaxTokens   int64           `koanf,json:"spark_max_tokens"`
	Temperature float64         `koanf,json:"spark_temperature"`
	TopK        int64           `koanf,json:"spark_topk"`
	ChatID      string          `koanf,json:"spark_chat_id"`
	Timeout     int             `koanf,json:"spark_timeout"` // seconds
}

type Config struct {
	OpenAI  *OpenAIConf    `koanf,json:"openai"`
	Spark   *IflySparkConf `koanf,json:"spark"`
	path    string
	workDir string
	koanfer *koanfer.JsonKoanfer
}

func NewConf(workDir string) (cfg *Config) {
	if ok, _ := gutils.PathIsExist(workDir); !ok {
		os.MkdirAll(workDir, os.ModePerm)
	}
	cfg = &Config{
		OpenAI:  &OpenAIConf{},
		Spark:   &IflySparkConf{},
		workDir: workDir,
	}
	cfg.path = filepath.Join(workDir, ConfigFileName)
	cfg.koanfer, _ = koanfer.NewKoanfer(cfg.path)
	if cfg.koanfer != nil {
		cfg.Reload()
	}
	if cfg.OpenAI.PromptMsgUrl == "" {
		cfg.OpenAI.PromptMsgUrl = PromptUrl
	}
	return
}

func (that *Config) GetWorkDir() string {
	return that.workDir
}

func (that *Config) Reload() {
	that.koanfer.Load(that)
}

func (that *Config) Save() {
	that.koanfer.Save(that)
}
