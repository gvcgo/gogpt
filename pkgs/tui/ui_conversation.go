package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type ConversationModel struct {
	Viewport viewport.Model
	TextArea textarea.Model
	Spinner  spinner.Model
}

func NewConversationModel() (cvm *ConversationModel) {
	cvm = &ConversationModel{
		Viewport: viewport.Model{},
		TextArea: textarea.Model{},
		Spinner:  spinner.Model{},
	}
	return
}

func (that *ConversationModel) Init() tea.Cmd {
	return nil
}

func (that *ConversationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return that, nil
}

func (that *ConversationModel) View() string {
	return ""
}
