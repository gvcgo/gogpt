package tui

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/moqsien/gogpt/pkgs/config"
	cvsation "github.com/moqsien/gogpt/pkgs/conversation"
	"github.com/moqsien/gogpt/pkgs/gpt"
	"github.com/moqsien/gogpt/pkgs/iflytek"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"
	openai "github.com/sashabaranov/go-openai"
)

type AnswerContinue string

type Bot interface {
	SendMsg(msgs []openai.ChatCompletionMessage) (m string, err error)
	RecvMsg() (m string, err error)
	Close()
	GetTokens() int64
}

type ConversationModel struct {
	Viewport     viewport.Model
	TextArea     textarea.Model
	Spinner      spinner.Model
	R            *glamour.TermRenderer
	CNF          *config.Config
	WindowHeight int
	WindowWidth  int
	GPT          *gpt.GPT
	Spark        *iflytek.Spark
	Conversation *cvsation.Conversation
	Receiving    bool
	Error        error
}

func NewConversationModel(cnf *config.Config) (cvm *ConversationModel) {
	cvm = &ConversationModel{
		CNF:          cnf,
		GPT:          gpt.NewGPT(cnf),
		Spark:        iflytek.NewSpark(cnf),
		Conversation: cvsation.NewConversation(cnf),
	}
	cvm.Conversation.SetBotType(cvsation.BotGPT) // ChatGPT by default
	cvm.Spinner = spinner.New(spinner.WithSpinner(spinner.Meter))
	cvm.TextArea = textarea.New()
	cvm.TextArea.Cursor.SetMode(cursor.CursorBlink)
	cvm.TextArea.Placeholder = "enter you message"
	cvm.TextArea.CharLimit = -1
	cvm.TextArea.FocusedStyle.CursorLine = lipgloss.NewStyle()
	cvm.TextArea.ShowLineNumbers = false
	cvm.TextArea.Focus()
	cvm.TextArea.CursorEnd()
	cvm.TextArea.SetHeight(2)
	cvm.Viewport = viewport.Model{}
	cvm.R, _ = glamour.NewTermRenderer(
		glamour.WithEnvironmentConfig(),
		glamour.WithWordWrap(0),
	)
	return
}

func (that *ConversationModel) GetBot() Bot {
	switch that.Conversation.BotType {
	case cvsation.BotSpark:
		if that.Spark == nil {
			that.Spark = iflytek.NewSpark(that.CNF)
		}
		return that.Spark
	default:
		if that.GPT == nil {
			that.GPT = gpt.NewGPT(that.CNF)
		}
		return that.GPT
	}
}

func (that *ConversationModel) SwitchBot() {
	if that.Conversation.BotType == cvsation.BotGPT {
		that.Conversation.SetBotType(cvsation.BotSpark)
		if that.GPT != nil {
			that.GPT.Close()
			that.GPT = nil
		}
	} else {
		that.Conversation.SetBotType(cvsation.BotGPT)
		if that.Spark != nil {
			that.Spark.Close()
			that.Spark = nil
		}
	}
}

func (that *ConversationModel) Init() tea.Cmd {
	return tea.Batch(that.Spinner.Tick, textarea.Blink)
}

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
		that.Viewport.Height = msg.Height - that.TextArea.Height() - lipgloss.Height(that.RenderFooter()) - lipgloss.Height(lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Render("title\n"))
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
				cmds = append(
					cmds, func() tea.Msg {
						return that.Spinner.Tick()
					},
				)
				answerStr, err := that.GetBot().SendMsg(msgList)

				if err == io.EOF {
					that.Receiving = false
				} else {
					cmds = append(cmds, func() tea.Msg {
						var msg AnswerContinue
						return msg
					})
				}
				that.Conversation.AddAnswer(answerStr, !that.Receiving)
				if err != nil && err != io.EOF {
					that.Error = err
					// clear errored answer, continue to Q&A
					that.Conversation.ClearCurrentAnswer()
					that.Receiving = false
				}
				that.Viewport.SetContent(that.RenderQA(*that.Conversation.Current))
				that.Viewport.GotoBottom()
			}
		case "up", "down":
			if !that.Receiving {
				that.Viewport, cmd = that.Viewport.Update(msg)
				cmds = append(cmds, cmd)
			}
		case "ctrl+p":
			if !that.Receiving {
				qa := that.Conversation.GetPrevQA()
				that.Viewport.SetContent(that.RenderQA(qa))
			}
		case "ctrl+f":
			if !that.Receiving {
				qa := that.Conversation.GetNextQA()
				that.Viewport.SetContent(that.RenderQA(qa))
			}
		case "ctrl+s":
			if !that.Receiving {
				that.Conversation.Save()
			}
		case "ctrl+l":
			if !that.Receiving {
				that.Conversation.Load()
			}
		case "ctrl+w":
			that.SwitchBot() // switch bot
		default:
			if !that.TextArea.Focused() && !that.Receiving {
				cmd = that.TextArea.Focus()
				cmds = append(cmds, cmd)
			}
			that.TextArea, cmd = that.TextArea.Update(msg)
			cmds = append(cmds, cmd)
		}
	case AnswerContinue:
		answerStr, err := that.GetBot().RecvMsg()
		if err == io.EOF {
			that.Receiving = false
		} else {
			cmds = append(cmds, func() tea.Msg {
				var msg AnswerContinue
				return msg
			})
		}
		if err != nil && err != io.EOF {
			that.Error = err
		}
		that.Conversation.AddAnswer(answerStr, !that.Receiving)

		if err != nil && err != io.EOF {
			that.Error = err
			// clear errored answer, continue to Q&A
			that.Conversation.ClearCurrentAnswer()
			that.Receiving = false
		}
		if that.Conversation.Current != nil {
			that.Viewport.SetContent(that.RenderQA(*that.Conversation.Current))
			that.Viewport.GotoBottom()
		}
	}
	return that, tea.Batch(cmds...)
}

func (that *ConversationModel) ContainsCJK(s string) bool {
	for _, r := range s {
		if unicode.In(r, unicode.Han, unicode.Hangul, unicode.Hiragana, unicode.Katakana) {
			return true
		}
	}
	return false
}

func (that *ConversationModel) EnsureTrailingNewline(s string) string {
	if !strings.HasSuffix(s, "\n") {
		return s + "\n"
	}
	return s
}

var (
	senderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5"))
	botStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	errorStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
	footerStyle = lipgloss.NewStyle().Height(1).Foreground(lipgloss.Color("#00FFFF")).Faint(true)
)

func (that *ConversationModel) RenderQA(qa cvsation.QuesAnsw) string {
	var (
		b       strings.Builder
		content string
	)
	b.WriteString(senderStyle.Render("You: "))

	content = qa.Q
	if that.ContainsCJK(content) {
		content = wrap.String(content, that.WindowWidth-5)
	} else {
		content = wordwrap.String(content, that.WindowWidth-5)
	}
	content, _ = that.R.Render(content)
	b.WriteString(that.EnsureTrailingNewline(content))

	b.WriteString(botStyle.Render("Bot: "))
	content = qa.A
	if that.ContainsCJK(content) {
		content = wrap.String(content, that.WindowWidth-5)
	} else {
		content = wordwrap.String(content, that.WindowWidth-5)
	}
	content, _ = that.R.Render(content)
	b.WriteString(that.EnsureTrailingNewline(content))
	return b.String()
}

func (that *ConversationModel) RenderFooter() string {
	if that.Error != nil {
		return footerStyle.Render(errorStyle.Render(fmt.Sprintf("error: %+v", that.Error)))
	}
	var columns []string

	// spinner
	if that.Receiving {
		columns = append(columns, that.Spinner.View())
	} else {
		columns = append(columns, that.Spinner.Spinner.Frames[0])
	}

	// bot type: ChatGPT/Spark
	if that.Conversation.BotType == cvsation.BotGPT {
		columns = append(columns, cvsation.BotGPT)
	} else {
		columns = append(columns, cvsation.BotSpark)
	}

	// conversation indicator
	if that.Conversation.Len() > 1 {
		conversationIdx := fmt.Sprintf("%s %d/%d", "Q&A", that.Conversation.Cursor+1, that.Conversation.Len())
		columns = append(columns, conversationIdx)
	}

	// tokens
	msgs := that.Conversation.GetMessages()
	if len(msgs) > 0 {
		token := cvsation.NumTokensFromMessages(msgs, that.CNF.OpenAI.Model)
		columns = append(columns, fmt.Sprintf("Tokens %d", token))
	}

	// switch tab
	columns = append(columns, "Tab ←/→")

	l := len(columns)
	length := that.WindowWidth / l
	for i := 0; i < l; i++ {
		columns[i] = footerStyle.Render(columns[i] + strings.Repeat(" ", length-len(columns[i])))
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, columns...)
}

func (that *ConversationModel) View() string {
	if that.WindowWidth == 0 || that.WindowHeight == 0 {
		return "Initializing..."
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		that.Viewport.View(),
		that.TextArea.View(),
		that.RenderFooter(),
	)
}

func (that *ConversationModel) CloseConversation() {
	if that.GPT != nil {
		that.GPT.Close()
	}
	if that.Spark != nil {
		that.Spark.Close()
	}
}
