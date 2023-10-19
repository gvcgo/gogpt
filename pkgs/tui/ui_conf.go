package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/moqsien/goutils/pkgs/gtea/input"
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

func GetGoGPTConfigModel() tea.Model {
	mi := input.NewInputMultiModel()
	mi.SetSubmitCmd(func() tea.Msg {
		k := tea.Key{
			Type:  tea.KeyTab,
			Runes: []rune("tab"),
		}
		return tea.KeyMsg(k)
	})

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
