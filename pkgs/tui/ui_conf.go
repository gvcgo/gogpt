package tui

import (
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gogf/gf/util/gconv"
	"github.com/moqsien/gogpt/pkgs/config"
	"github.com/moqsien/gogpt/pkgs/gpt"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/gtea/input"
	"github.com/moqsien/goutils/pkgs/gutils"
	openai "github.com/sashabaranov/go-openai"
)

/*
Gogpt Config Model
*/
var (
	baseUrl        string = "baseUrl"
	gptModel       string = "select_model"
	apiKey         string = "apiKey"
	proxy          string = "proxy"
	apiType        string = "apiType"
	apiVersion     string = "apiVersion"
	orgID          string = "orgID"
	engine         string = "engine"
	limit          string = "empty_limit"
	maxTokens      string = "maxTokens"
	ctxLen         string = "contextLength"
	temperature    string = "temperature"
	gptPrompt      string = "select_prompt"
	gptPromptValue string = "enter_prompt"
)

func GetGoGPTConfigModel(prompt *gpt.GPTPrompt) ExtraModel {
	mi := input.NewInputMultiModel()
	mi.AddOneInput(baseUrl, input.MWithPlaceholder("base_url"), input.MWithWidth(150))
	gptModelList := []string{
		openai.GPT3Dot5Turbo0613,
		openai.GPT3Dot5Turbo,
		openai.GPT432K0613,
		openai.GPT4,
		openai.GPT432K0314,
		openai.GPT432K,
		openai.GPT40613,
		openai.GPT40314,
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
	mi.AddOneOption(gptModel, gptModelList, input.MWithPlaceholder("gpt_model"), input.MWithWidth(100))
	gptPromptList := []string{}
	for _, item := range *prompt.PromptList {
		gptPromptList = append(gptPromptList, item.Title)
	}
	mi.AddOneOption(gptPrompt, gptPromptList, input.MWithPlaceholder("gpt_prompt"), input.MWithWidth(100))
	mi.AddOneInput(gptPromptValue, input.MWithPlaceholder("enter_gpt_prompt"), input.MWithWidth(100))

	mi.AddOneInput(apiKey, input.MWithPlaceholder("api_key"), input.MWithWidth(100))
	mi.AddOneInput(proxy, input.MWithPlaceholder("proxy"), input.MWithWidth(150))
	mi.AddOneInput(apiVersion, input.MWithPlaceholder("api_version"), input.MWithWidth(100))

	gptApiTypeList := []string{
		string(openai.APITypeOpenAI),
		string(openai.APITypeAzure),
		string(openai.APITypeAzureAD),
	}
	mi.AddOneOption(apiType, gptApiTypeList, input.MWithPlaceholder("gpt_api_type"), input.MWithWidth(100))

	mi.AddOneInput(orgID, input.MWithPlaceholder("org_id"), input.MWithWidth(100))
	mi.AddOneInput(engine, input.MWithPlaceholder("engine"), input.MWithWidth(100))
	mi.AddOneInput(limit, input.MWithPlaceholder("empty_message_limit"), input.MWithWidth(100))
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
		if values[apiType] != "" {
			cfg.OpenAI.ApiType = openai.APIType(values[apiType])
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
		if values[gptPrompt] != "" {
			cfg.OpenAI.PromptStr = values[gptPrompt]
		}
		if values[gptPromptValue] != "" {
			cfg.OpenAI.PromptStr = values[gptPromptValue]
		}

		cfg.OpenAI.EmptyMessagesLimit = gconv.Uint(values[limit])
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
	cfg.OpenAI.PromptMsgUrl = config.PromptUrl
	cfg.Save()
}

func GetDefaultConfig() (conf *config.Config) {
	homeDir, _ := os.UserHomeDir()
	workDir := filepath.Join(homeDir, ".gogpt")
	confPath := filepath.Join(workDir, config.ConfigFileName)
	cfg := config.NewConf(workDir)
	prompt := gpt.NewGPTPrompt(cfg)
	if ok, _ := gutils.PathIsExist(confPath); !ok {
		m := GetGoGPTConfigModel(prompt)
		pgm := tea.NewProgram(m)
		if _, err := pgm.Run(); err != nil {
			gprint.PrintError("%+v", err)
		}
		cfg.OpenAI.Model = openai.GPT3Dot5Turbo
		SetConfig(cfg, m.Values())
	}
	return cfg
}
