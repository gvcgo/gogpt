package tui

import (
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gogf/gf/util/gconv"
	"github.com/moqsien/gogpt/pkgs/config"
	"github.com/moqsien/gogpt/pkgs/gpt"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/gtea/input"
	"github.com/moqsien/goutils/pkgs/gutils"
	openai "github.com/sashabaranov/go-openai"
)

type PromptString string

func (that PromptString) Less(prompt gutils.IComparable) bool {
	p := prompt.(PromptString)
	r := strings.Compare(string(that), string(p))
	return r == -1
}

/*
Gogpt Config Model
*/
var (
	baseUrl        string = "baseUrl"
	gptModel       string = "select_model"
	apiKey         string = "apiKey"
	proxy          string = "proxy"
	apiType        string = "select_APIType"
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
	mi.SetInputPromptFormat("%-15s")
	placeHolderStyle := input.MWithPlaceholderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#BEBEBE")))
	mi.AddOneInput(apiKey, input.MWithPlaceholder("ChatGPT auth token"), input.MWithWidth(100), placeHolderStyle)
	mi.AddOneInput(proxy, input.MWithPlaceholder("Local proxy"), input.MWithWidth(150), placeHolderStyle)
	mi.AddOneInput(ctxLen, input.MWithPlaceholder("Conversation context length"), input.MWithWidth(100), placeHolderStyle)

	// Select ChatGPT API type
	gptApiTypeList := []string{
		string(openai.APITypeOpenAI),
		string(openai.APITypeAzure),
		string(openai.APITypeAzureAD),
	}
	mi.AddOneOption(apiType, gptApiTypeList, input.MWithPlaceholder("ChatGPT Api Type."), input.MWithWidth(100), placeHolderStyle)

	// Select ChatGPT Model
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
	mi.AddOneOption(gptModel, gptModelList, input.MWithPlaceholder("ChatGPT Model."), input.MWithWidth(100), placeHolderStyle)

	// Select ChatGPT Prompt
	gptPromptList := []gutils.IComparable{}
	for _, item := range *prompt.PromptList {
		gptPromptList = append(gptPromptList, PromptString(item.Title))
	}
	gutils.QuickSort(gptPromptList, 0, len(gptPromptList)-1)
	pList := []string{}
	for _, p := range gptPromptList {
		pStr := p.(PromptString)
		pList = append(pList, string(pStr))
	}
	mi.AddOneOption(gptPrompt, pList, input.MWithPlaceholder("gpt_prompt"), input.MWithWidth(100), placeHolderStyle)
	// Enter you own ChatGPT Prompt
	mi.AddOneInput(gptPromptValue, input.MWithPlaceholder("Enter your own prompt info instead of a selection from above."), input.MWithWidth(100), placeHolderStyle)
	// Some configs
	mi.AddOneInput(limit, input.MWithPlaceholder("Max empty message limit. Int."), input.MWithWidth(100), placeHolderStyle)
	mi.AddOneInput(maxTokens, input.MWithPlaceholder("Max tokens. Int."), input.MWithWidth(100), placeHolderStyle)
	mi.AddOneInput(temperature, input.MWithPlaceholder("Temperautue. Float."), input.MWithWidth(100), placeHolderStyle)

	// Custom baseUrl
	mi.AddOneInput(baseUrl, input.MWithPlaceholder("defaul:https://api.openai.com/v1"), input.MWithWidth(150), placeHolderStyle)
	// For AzureGPT
	mi.AddOneInput(apiVersion, input.MWithPlaceholder("API version."), input.MWithWidth(100), placeHolderStyle)
	mi.AddOneInput(orgID, input.MWithPlaceholder("Organization ID."), input.MWithWidth(100), placeHolderStyle)
	mi.AddOneInput(engine, input.MWithPlaceholder("Engine."), input.MWithWidth(100), placeHolderStyle)
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
