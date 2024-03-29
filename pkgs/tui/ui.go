package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gvcgo/goutils/pkgs/gtea/gprint"
	"github.com/gvcgo/gogpt/pkgs/config"
	"github.com/gvcgo/gogpt/pkgs/gpt"
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
	g.AddHelpInfo()
	return
}

func (that *GPTUI) AddConversationUI() {
	uconv := NewConversationModel(that.CNF)
	that.GVM.AddTab("Conversation", uconv)
}

func (that *GPTUI) AddConfUI() {
	uconf := GetGoGPTConfigModel(that.Prompt, that.CNF)
	uconf.SetSubmitCmd(func() tea.Msg {
		vals := uconf.Values()
		vals[gptPrompt] = that.Prompt.GetPromptByTile(vals[gptPrompt])
		SetConfig(that.CNF, vals)
		return returnFirst
	})
	that.GVM.AddTab("Configuration", uconf)
}

func (that *GPTUI) AddHelpInfo() {
	helpInfo := NewHelpModel()
	that.GVM.AddTab("HelpInfo", helpInfo)
}

func (that *GPTUI) Run() {
	if that.Program == nil {
		that.Program = tea.NewProgram(that.GVM, tea.WithAltScreen())
	}
	if _, err := that.Program.Run(); err != nil {
		gprint.PrintError("Run bubbletea failed: %+v", err)
	}
}
