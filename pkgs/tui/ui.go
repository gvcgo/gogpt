package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/moqsien/gogpt/pkgs/config"
	"github.com/moqsien/gogpt/pkgs/gpt"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
)

var (
	returnFirst = ReturnFirst("")
)

type ExtraModel interface {
	tea.Model
	SetSubmitCmd(tea.Cmd)
	Values() map[string]string
}

type GPTUI struct {
	Program *tea.Program
	GVM     *GPTViewModel
	CNF     *config.Config
	Prompt  *gpt.GPTPrompt
}

func NewGPTUI(cnf *config.Config) (g *GPTUI) {
	g = &GPTUI{
		GVM:    NewGPTViewModel(),
		CNF:    cnf,
		Prompt: gpt.NewGPTPrompt(cnf),
	}
	g.AddConversationUI()
	g.AddConfUI()
	g.AddGPTModelSelectorUI()
	g.AddPromptSelectorUI()
	return
}

func (that *GPTUI) AddConversationUI() {}

func (that *GPTUI) AddConfUI() {
	uconf := GetGoGPTConfigModel()
	uconf.SetSubmitCmd(func() tea.Msg {
		SetConfig(that.CNF, uconf.Values())
		return returnFirst
	})
	that.GVM.AddTab("Configuration", uconf)
}

func (that *GPTUI) AddPromptSelectorUI() {
	uprompt := GetPromptModel(that.Prompt)
	uprompt.SetSubmitCmd(func() tea.Msg {
		SetGPTPrompt(that.Prompt, uprompt.Values())
		return returnFirst
	})
	that.GVM.AddTab("Prompt", uprompt)
}

func (that *GPTUI) AddGPTModelSelectorUI() {
	ugmodel := GetModelSelector()
	ugmodel.SetSubmitCmd(func() tea.Msg {
		SetGPTModel(that.CNF, ugmodel.Values())
		return returnFirst
	})
	that.GVM.AddTab("GPTModel", ugmodel)
}

func (that *GPTUI) Run() {
	if that.Program == nil {
		that.Program = tea.NewProgram(that.GVM)
	}
	if _, err := that.Program.Run(); err != nil {
		gprint.PrintError("%+v", err)
	}
}
