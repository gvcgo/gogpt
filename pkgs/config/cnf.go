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

type OpenAIConf struct {
	BaseUrl            string         `koanf,json:"base_url"`
	ApiKey             string         `koanf,json:"api_key"`
	ApiType            openai.APIType `koanf,json:"api_type"`
	ApiVersion         string         `koanf,json:"api_version"`
	OrgID              string         `koanf,json:"org_id"`
	Engine             string         `koanf,json:"engine"`
	EmptyMessagesLimit uint           `koanf,json:"empty_msg_limit"`
	Proxy              string         `koanf,json:"proxy"`
	TimeOut            int64          `koanf,json:"timeout_seconds"`
	Model              string         `koanf,json:"model"`
	MaxTokens          int            `koanf,json:"max_tokens"`
	ContextLen         int            `koanf,json:"context_length"`
	Temperature        float32        `koanf,json:"temperature"`
	PromptMsgUrl       string         `koanf,json:"prompt_msgs_url"`
	PromptStr          string         `koanf,json:"prompt"`
}

type Config struct {
	OpenAI  *OpenAIConf `koanf,json:"openai"`
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
