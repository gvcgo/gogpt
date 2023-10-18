package config

import (
	"os"
	"path/filepath"

	"github.com/gogf/gf/v2/util/gconv"
	"github.com/moqsien/goutils/pkgs/gtea/input"
	"github.com/moqsien/goutils/pkgs/gtea/selector"
	"github.com/moqsien/goutils/pkgs/gutils"
	"github.com/moqsien/goutils/pkgs/koanfer"
	"github.com/sashabaranov/go-openai"
)

const (
	ConfigFileName string = "gogpt_conf.json"
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
		cfg.OpenAI.PromptMsgUrl = "https://gitlab.com/moqsien/gpt_resources/-/raw/main/prompt.json"
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

/*
Set configurations
*/
func SetConfig(workDir string) {
	cfg := NewConf(workDir)
	cfg.Reload()

	selectorItems := selector.NewItemList()
	selectorItems.Add(string(openai.APITypeOpenAI), openai.APITypeOpenAI)
	selectorItems.Add(string(openai.APITypeAzure), openai.APITypeAzure)
	selectorItems.Add(string(openai.APITypeAzureAD), openai.APITypeAzureAD)
	sel := selector.NewSelector(
		selectorItems,
		selector.WithTitle("Choose APIType"),
		selector.WithHeight(5),
		selector.WidthEnableMulti(false),
		selector.WithEnbleInfinite(true),
		selector.WithFilteringEnabled(false),
	)
	sel.Run()
	val := sel.Value()[0]
	cfg.OpenAI.ApiType = val.(openai.APIType)

	models := []string{
		openai.GPT4,
		openai.GPT432K0613,
		openai.GPT432K0314,
		openai.GPT432K,
		openai.GPT40613,
		openai.GPT40314,
		openai.GPT3Dot5Turbo,
		openai.GPT3Dot5Turbo0613,
		openai.GPT3Dot5Turbo0301,
		openai.GPT3Dot5Turbo16K,
		openai.GPT3Dot5Turbo16K0613,
		openai.GPT3Dot5TurboInstruct,
		openai.GPT3Davinci,
		openai.GPT3Davinci002,
		openai.GPT3Curie,
		openai.GPT3Curie002,
		openai.GPT3Ada,
		openai.GPT3Ada002,
		openai.GPT3Babbage,
		openai.GPT3Babbage002,
	}
	selectorItems = selector.NewItemList()
	for _, model := range models {
		selectorItems.Add(model, model)
	}
	sel = selector.NewSelector(
		selectorItems,
		selector.WithTitle("Choose model"),
		selector.WithHeight(15),
		selector.WidthEnableMulti(false),
		selector.WithEnbleInfinite(true),
		selector.WithFilteringEnabled(false),
	)
	sel.Run()
	val = sel.Value()[0]
	cfg.OpenAI.Model = val.(string)

	var (
		baseUrl     string = "baseUrl"
		apiKey      string = "apiKey"
		proxy       string = "proxy"
		apiVersion  string = "apiVersion"
		orgID       string = "orgID"
		engine      string = "engine"
		limit       string = "limit"
		timeout     string = "timeout"
		maxTokens   string = "maxTokens"
		ctxLen      string = "contextLength"
		temperature string = "temperature"
	)

	mi := input.NewMultiInput()
	mi.AddOneItem(baseUrl, input.MWithPlaceholder("base_url"), input.MWithWidth(150))
	mi.AddOneItem(apiKey, input.MWithPlaceholder("api_key"), input.MWithWidth(100))
	mi.AddOneItem(proxy, input.MWithPlaceholder("proxy"), input.MWithWidth(150))
	mi.AddOneItem(apiVersion, input.MWithPlaceholder("api_version"), input.MWithWidth(100))
	mi.AddOneItem(orgID, input.MWithPlaceholder("org_id"), input.MWithWidth(100))
	mi.AddOneItem(engine, input.MWithPlaceholder("engine"), input.MWithWidth(100))
	mi.AddOneItem(limit, input.MWithPlaceholder("empty_message_limit"), input.MWithWidth(100))
	mi.AddOneItem(timeout, input.MWithPlaceholder("timeout"), input.MWithWidth(100))
	mi.AddOneItem(maxTokens, input.MWithPlaceholder("max_tokens"), input.MWithWidth(100))
	mi.AddOneItem(ctxLen, input.MWithPlaceholder("context_length"), input.MWithWidth(100))
	mi.AddOneItem(temperature, input.MWithPlaceholder("temperature"), input.MWithWidth(100))
	mi.Run()

	values := mi.Values()
	if values[baseUrl] != "" {
		cfg.OpenAI.BaseUrl = values[baseUrl]
	}
	if values[apiKey] != "" {
		cfg.OpenAI.ApiKey = values[apiKey]
	}
	if values[proxy] != "" {
		cfg.OpenAI.Proxy = values[proxy]
	}
	if values[apiVersion] != "" {
		cfg.OpenAI.ApiVersion = values[apiVersion]
	}
	if values[orgID] != "" {
		cfg.OpenAI.OrgID = values[orgID]
	}
	if values[engine] != "" {
		cfg.OpenAI.Engine = values[engine]
	}

	cfg.OpenAI.EmptyMessagesLimit = gconv.Uint(values[limit])
	tt := gconv.Int64(values[timeout])
	if tt <= 0 {
		tt = 30
	}
	cfg.OpenAI.TimeOut = tt
	mTokens := gconv.Int(values[maxTokens])
	if mTokens == 0 {
		mTokens = 1024
	}
	cfg.OpenAI.MaxTokens = mTokens

	cLen := gconv.Int(values[ctxLen])
	if cLen == 0 {
		cLen = 6
	}
	cfg.OpenAI.ContextLen = cLen
	cfg.OpenAI.Temperature = gconv.Float32(values[temperature])

	cfg.Save()
}

func GetDefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	workDir := filepath.Join(homeDir, ".gogpt")
	confPath := filepath.Join(workDir, ConfigFileName)
	if ok, _ := gutils.PathIsExist(confPath); !ok {
		SetConfig(workDir)
	}
	return NewConf(workDir)
}
