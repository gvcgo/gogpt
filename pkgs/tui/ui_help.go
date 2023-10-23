package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HelpModel struct {
	helpStyle lipgloss.Style
}

func NewHelpModel() (h *HelpModel) {
	h = &HelpModel{
		// "#FFA500" "#808080"
		helpStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")),
	}
	return
}

func (that *HelpModel) Init() tea.Cmd {
	return nil
}

func (that *HelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return that, nil
}

func (that *HelpModel) View() string {
	pattern := "%-10s  %s"
	helpList := []string{
		fmt.Sprintf(pattern, "enter", "Submit your message to gpt."),
		fmt.Sprintf(pattern, "↑", "Scroll up."),
		fmt.Sprintf(pattern, "↓", "Scroll down."),
		fmt.Sprintf(pattern, "ctrl+p", "Show the previous QA."),
		fmt.Sprintf(pattern, "ctrl+f", "Show the next QA."),
		fmt.Sprintf(pattern, "ctrl+s", "Save conversation."),
		fmt.Sprintf(pattern, "ctrl+l", "Load conversation."),
		fmt.Sprintf(pattern, "→", "Switch to the next Tab."),
		fmt.Sprintf(pattern, "←", "Switch to the previous Tab."),
	}
	r := []string{}
	for _, str := range helpList {
		r = append(r, that.helpStyle.Render(str))
	}
	return lipgloss.JoinVertical(lipgloss.Left, r...)
}
