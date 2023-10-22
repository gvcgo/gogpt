package tui

import (
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
	return tea.Batch(that.Spinner.Tick, textarea.Blink)
}

// TODO: keymap & viewpord render
func (that *ConversationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		that.WindowWidth = msg.Width
		that.WindowHeight = msg.Height
		that.TextArea.SetWidth(that.WindowWidth)
		that.Viewport.Width = msg.Width - 5
		that.Viewport.MouseWheelEnabled = true
		that.Viewport.Height = msg.Height - that.TextArea.Height() - lipgloss.Height(that.Spinner.View()) - lipgloss.Height(lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Render("title\n"))
	case spinner.TickMsg:
		if that.Receiving {
			that.Spinner, cmd = that.Spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	case tea.KeyMsg:
		switch keyPress := msg.String(); keyPress {
		case "enter":
			messageStr := that.TextArea.Value()
			that.TextArea.Reset()
			that.TextArea.Blur()
			if messageStr != "" {
				that.Conversation.AddQuestion(messageStr)
				msgList := that.Conversation.GetMessages()
				that.Receiving = true
				answerStr, err := that.GPT.SendMsg(msgList)
				if err == io.EOF {
					that.Receiving = false
				} else {
					cmds = append(cmds, func() tea.Msg {
						var msg AnswerContinue
						return msg
					})
				}
				that.Conversation.AddAnswer(answerStr, !that.Receiving)
				that.Viewport.SetContent(that.Conversation.Current.A)
				that.Viewport.GotoBottom()
			}
		default:
			if !that.TextArea.Focused() && !that.Receiving {
				cmd = that.TextArea.Focus()
				cmds = append(cmds, cmd)
			}
			that.TextArea, cmd = that.TextArea.Update(msg)
			cmds = append(cmds, cmd)
		}
	case AnswerContinue:
		answerStr, err := that.GPT.RecvMsg()
		if err == io.EOF {
			that.Receiving = false
		} else {
			cmds = append(cmds, func() tea.Msg {
				var msg AnswerContinue
				return msg
			})
		}
		that.Conversation.AddAnswer(answerStr, !that.Receiving)
		if that.Conversation.Current != nil {
			that.Viewport.SetContent(that.Conversation.Current.A)
			that.Viewport.GotoBottom()
		}
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
