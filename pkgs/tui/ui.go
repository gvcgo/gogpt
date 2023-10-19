package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/moqsien/gogpt/pkgs/config"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
)

var (
	returnFirst = ReturnFirst("")
)

type ExtraModel interface {
	tea.Model
	SetSubmitCmd(tea.Cmd)
}

type GPTUI struct {
	Program *tea.Program
	GVM     *GPTViewModel
	CNF     *config.Config
}

func NewGPTUI(cnf *config.Config) (g *GPTUI) {
	g = &GPTUI{
		GVM: NewGPTViewModel(),
		CNF: cnf,
	}
	uconf := GetGoGPTConfigModel()
	uconf.SetSubmitCmd(func() tea.Msg {
		return returnFirst
	})
	g.GVM.AddTab("Configuration", uconf)

	uprompt := GetPromptModel(cnf)
	uprompt.SetSubmitCmd(func() tea.Msg {
		return returnFirst
	})
	g.GVM.AddTab("Prompt", uprompt)

	ugmodel := GetModelSelector()
	ugmodel.SetSubmitCmd(func() tea.Msg {
		return returnFirst
	})
	g.GVM.AddTab("GPTModel", ugmodel)
	return
}

func (that *GPTUI) Run() {
	if that.Program == nil {
		that.Program = tea.NewProgram(that.GVM)
	}
	if _, err := that.Program.Run(); err != nil {
		gprint.PrintError("%+v", err)
	}
}
