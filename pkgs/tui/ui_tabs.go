package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/*
tab
*/
type Tab struct {
	Title string
	Model tea.Model
}

/*
Customed messgaes
*/
type ReturnFirst string

/*
GPT UI Model
*/
type GPTViewModel struct {
	TabList   []*Tab
	ActiveTab int
}

func NewGPTViewModel() (gm *GPTViewModel) {
	gm = &GPTViewModel{
		TabList: []*Tab{},
	}
	return
}

func (that *GPTViewModel) Init() tea.Cmd {
	teaCmdList := []tea.Cmd{}
	for _, m := range that.TabList {
		teaCmd := m.Model.Init()
		if teaCmd != nil {
			teaCmdList = append(teaCmdList, teaCmd)
		}
	}
	if len(teaCmdList) > 0 {
		return tea.Batch(teaCmdList...)
	}
	return nil
}

func (that *GPTViewModel) GetCurrentModel() tea.Model {
	return that.TabList[that.ActiveTab].Model
}

func (that *GPTViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	currentModel := that.GetCurrentModel()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q", "esc":
			return that, tea.Quit
		case "right":
			if that.ActiveTab < len(that.TabList)-1 {
				that.ActiveTab++
			} else {
				that.ActiveTab = 0
			}
		case "left":
			if that.ActiveTab > 0 {
				that.ActiveTab--
			} else {
				that.ActiveTab = len(that.TabList) - 1
			}
		default:
			_, cmd := currentModel.Update(msg)
			return that, cmd
		}
	case ReturnFirst:
		that.ActiveTab = 0
		return that, nil
	default:
		_, cmd := currentModel.Update(msg)
		return that, cmd
	}
	return that, nil
}

func (that *GPTViewModel) View() string {
	doc := strings.Builder{}
	var newTabs []string
	var style lipgloss.Style
	for i, t := range that.TabList {
		if i == that.ActiveTab {
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("229"))
		} else {
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("57"))
		}
		newTabs = append(newTabs, style.Render(t.Title))
	}
	row := strings.Join(newTabs, " | ")
	doc.WriteString(row)
	currentModel := that.GetCurrentModel()
	return lipgloss.JoinVertical(lipgloss.Left, doc.String()+"\n", currentModel.View())
}

func (that *GPTViewModel) AddTab(title string, model tea.Model) {
	that.TabList = append(that.TabList, &Tab{Title: title, Model: model})
}
