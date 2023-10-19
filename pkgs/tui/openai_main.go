package tui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	Tabs       []string
	TabContent []string
	activeTab  int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l", "n", "tab":
			m.activeTab = min(m.activeTab+1, len(m.Tabs)-1)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}
	}

	return m, nil
}

var (
	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
)

func (m model) View() string {
	doc := strings.Builder{}
	var newTabs []string
	var style lipgloss.Style
	for i, t := range m.Tabs {
		if i == m.activeTab {
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("229"))
		} else {
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("57"))
		}
		newTabs = append(newTabs, style.Render(t))
	}
	// row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	row := strings.Join(newTabs, " | ")
	doc.WriteString(row)
	doc.WriteString("\n")
	return docStyle.Render(doc.String())
}

func RunTab() {
	tabs := []string{"Lip Gloss", "Blush", "Eye Shadow", "Mascara", "Foundation"}
	tabContent := []string{"Lip Gloss Tab", "Blush Tab", "Eye Shadow Tab", "Mascara Tab", "Foundation Tab"}
	m := model{Tabs: tabs, TabContent: tabContent}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
