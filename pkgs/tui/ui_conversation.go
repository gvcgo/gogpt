package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/moqsien/gogpt/pkgs/config"
)

type ConversationModel struct {
	Viewport     viewport.Model
	TextArea     textarea.Model
	Spinner      spinner.Model
	CNF          *config.Config
	WindowHeight int
	WindowWidth  int
}

func NewConversationModel(cnf *config.Config) (cvm *ConversationModel) {
	cvm = &ConversationModel{
		CNF: cnf,
	}
	cvm.Spinner = spinner.New()
	cvm.TextArea = textarea.New()
	cvm.Viewport = viewport.New(20, 10)
	return
}

func (that *ConversationModel) Init() tea.Cmd {
	return textarea.Blink
}

func (that *ConversationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		that.WindowWidth = msg.Width
		that.WindowHeight = msg.Height
		that.Viewport.SetContent(fmt.Sprintf("height: %d, width: %d", that.WindowHeight, that.WindowWidth))
		that.Viewport.Width = msg.Width
		that.Viewport.Height = msg.Height - that.TextArea.Height() - lipgloss.Height(that.Spinner.View()) - lipgloss.Height(lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Render("title\n"))
	}
	return that, nil
}

func (that *ConversationModel) View() string {
	if that.WindowWidth == 0 || that.WindowHeight == 0 {
		return "Initializing..."
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		that.Viewport.View(),
		that.TextArea.View(),
		that.Spinner.View(),
	)
}
