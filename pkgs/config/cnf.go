package config

import (
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/moqsien/goutils/pkgs/gtea/input"
	"github.com/moqsien/goutils/pkgs/gtea/selector"
	"github.com/moqsien/goutils/pkgs/koanfer"
	"github.com/sashabaranov/go-openai"
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
	MaxTokens          int            `koanf,json:"max_tokens"`
	ContextLen         int            `koanf,json:"context_length"`
	Temperature        float32        `koanf,json:"temperature"`
}

type Config struct {
	OpenAI  *OpenAIConf `koanf,json:"openai"`
	path    string
	koanfer *koanfer.JsonKoanfer
}

func NewConf(cfgPath string) (cfg *Config) {
	cfg = &Config{}
	cfg.path = cfgPath
	cfg.koanfer, _ = koanfer.NewKoanfer(cfgPath)
	if cfg.koanfer != nil {
		cfg.Reload()
	}
	return
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
func SetConfig(cfgPath string) {
	cfg := NewConf(cfgPath)
	cfg.Reload()

	selectorItems := selector.NewItemList()
	selectorItems.Add(string(openai.APITypeOpenAI), openai.APITypeOpenAI)
	selectorItems.Add(string(openai.APITypeAzure), openai.APITypeAzure)
	selectorItems.Add(string(openai.APITypeAzureAD), openai.APITypeAzureAD)
	sel := selector.NewSelector(
		selectorItems,
		selector.WidthEnableMulti(false),
		selector.WithEnbleInfinite(true),
		selector.WithFilteringEnabled(false),
	)
	sel.Run()
	val := sel.Value()[0]
	cfg.OpenAI.ApiType = val.(openai.APIType)

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
	cfg.OpenAI.BaseUrl = values[baseUrl]
	cfg.OpenAI.ApiKey = values[apiKey]
	cfg.OpenAI.Proxy = values[proxy]
	cfg.OpenAI.ApiVersion = values[apiVersion]
	cfg.OpenAI.OrgID = values[orgID]
	cfg.OpenAI.Engine = values[engine]
	cfg.OpenAI.EmptyMessagesLimit = gconv.Uint(values[limit])
	tt := gconv.Int64(values[timeout])
	if tt <= 0 {
		tt = 30
	}
	cfg.OpenAI.TimeOut = tt
	cfg.OpenAI.MaxTokens = gconv.Int(values[maxTokens])
	cfg.OpenAI.ContextLen = gconv.Int(values[ctxLen])
	cfg.OpenAI.Temperature = gconv.Float32(values[temperature])

	cfg.Save()
}
