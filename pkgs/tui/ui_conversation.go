package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/moqsien/gogpt/pkgs/config"
	"github.com/moqsien/gogpt/pkgs/gpt"
)

type AnswerContinue string

type ConversationModel struct {
	Viewport     viewport.Model
	TextArea     textarea.Model
	Spinner      spinner.Model
	CNF          *config.Config
	WindowHeight int
	WindowWidth  int
	GPT          *gpt.GPT
	Conversation *gpt.Conversation
	Receiving    bool
}

func NewConversationModel(cnf *config.Config) (cvm *ConversationModel) {
	cvm = &ConversationModel{
		CNF:          cnf,
		GPT:          gpt.NewGPT(cnf),
		Conversation: gpt.NewConversation(cnf),
	}
	cvm.Spinner = spinner.New()
	cvm.TextArea = textarea.New()
	cvm.TextArea.Cursor.SetMode(cursor.CursorBlink)
	cvm.TextArea.Placeholder = "enter you message"
	cvm.TextArea.CharLimit = -1
	cvm.TextArea.FocusedStyle.CursorLine = lipgloss.NewStyle()
	cvm.TextArea.ShowLineNumbers = false
	cvm.TextArea.Focus()
	cvm.TextArea.CursorEnd()
	cvm.Viewport = viewport.Model{}
	return
}

func (that *ConversationModel) Init() tea.Cmd {
	return textarea.Blink
}

func (that *ConversationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		that.WindowWidth = msg.Width
		that.WindowHeight = msg.Height
		that.TextArea.SetWidth(that.WindowWidth)
		that.Viewport.SetContent(fmt.Sprintf("height: %d, width: %d", that.WindowHeight, that.WindowWidth))
		that.Viewport.Width = msg.Width
		that.Viewport.Height = msg.Height - that.TextArea.Height() - lipgloss.Height(that.Spinner.View()) - lipgloss.Height(lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Render("title\n"))
	case tea.KeyMsg:
		switch keyPress := msg.String(); keyPress {
		case "enter":
			messageStr := that.TextArea.Value()
			if messageStr != "" {
				that.Conversation.AddQuestion(messageStr)
				msgList := that.Conversation.GetMessages()
				answerStr, err := that.GPT.SendMsg(msgList)
				isCompleted := false
				if err == io.EOF {
					isCompleted = true
				} else {
					cmds = append(cmds, func() tea.Msg {
						return AnswerContinue("")
					})
				}
				that.Conversation.AddAnswer(answerStr, isCompleted)

				that.Viewport.SetContent(that.Conversation.Current.A)
			}
			that.TextArea.Reset()
			that.TextArea.Blur()
		default:
			if !that.TextArea.Focused() {
				cmd = that.TextArea.Focus()
				cmds = append(cmds, cmd)
			}
			that.TextArea, cmd = that.TextArea.Update(msg)
			cmds = append(cmds, cmd)
		}
	case AnswerContinue:
		answerStr, err := that.GPT.RecvMsg()
		isCompleted := false
		if err == io.EOF {
			isCompleted = true
		} else {
			cmds = append(cmds, func() tea.Msg {
				return AnswerContinue("")
			})
		}
		that.Conversation.AddAnswer(answerStr, isCompleted)
		that.Viewport.SetContent(that.Conversation.Current.A)
	}
	return that, tea.Batch(cmds...)
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
