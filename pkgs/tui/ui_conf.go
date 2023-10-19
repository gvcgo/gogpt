package tui

import (
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gogf/gf/util/gconv"
	"github.com/moqsien/gogpt/pkgs/config"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/gtea/input"
	"github.com/moqsien/goutils/pkgs/gutils"
	openai "github.com/sashabaranov/go-openai"
)

/*
Gogpt Config Model
*/
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

func GetGoGPTConfigModel() ExtraModel {
	mi := input.NewInputMultiModel()
	mi.AddOneInput(baseUrl, input.MWithPlaceholder("base_url"), input.MWithWidth(150))
	mi.AddOneInput(apiKey, input.MWithPlaceholder("api_key"), input.MWithWidth(100))
	mi.AddOneInput(proxy, input.MWithPlaceholder("proxy"), input.MWithWidth(150))
	mi.AddOneInput(apiVersion, input.MWithPlaceholder("api_version"), input.MWithWidth(100))
	mi.AddOneInput(orgID, input.MWithPlaceholder("org_id"), input.MWithWidth(100))
	mi.AddOneInput(engine, input.MWithPlaceholder("engine"), input.MWithWidth(100))
	mi.AddOneInput(limit, input.MWithPlaceholder("empty_message_limit"), input.MWithWidth(100))
	mi.AddOneInput(timeout, input.MWithPlaceholder("timeout"), input.MWithWidth(100))
	mi.AddOneInput(maxTokens, input.MWithPlaceholder("max_tokens"), input.MWithWidth(100))
	mi.AddOneInput(ctxLen, input.MWithPlaceholder("context_length"), input.MWithWidth(100))
	mi.AddOneInput(temperature, input.MWithPlaceholder("temperature"), input.MWithWidth(100))
	return mi
}

func SetConfig(cfg *config.Config, values map[string]string) {
	if cfg == nil {
		gprint.PrintError("conf object is nil!")
		return
	}
	cfg.Reload()
	if len(values) > 0 {
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
	}
	cfg.OpenAI.PromptMsgUrl = "https://gitlab.com/moqsien/gpt_resources/-/raw/main/prompt.json"
	cfg.Save()
}

func GetDefaultConfig() (conf *config.Config) {
	homeDir, _ := os.UserHomeDir()
	workDir := filepath.Join(homeDir, ".gogpt")
	confPath := filepath.Join(workDir, config.ConfigFileName)
	cfg := config.NewConf(workDir)
	if ok, _ := gutils.PathIsExist(confPath); !ok {
		m := GetGoGPTConfigModel()
		pgm := tea.NewProgram(m)
		if _, err := pgm.Run(); err != nil {
			gprint.PrintError("%+v", err)
		}
		cfg.OpenAI.Model = openai.GPT3Dot5Turbo
		SetConfig(cfg, m.Values())
	}
	return cfg
}
