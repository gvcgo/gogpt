package gpt

import (
	"path/filepath"

	"github.com/moqsien/gogpt/pkgs/config"
	"github.com/moqsien/goutils/pkgs/koanfer"
	"github.com/sashabaranov/go-openai"
)

/*
Manage Chatgpt conversation
*/

const (
	ConversationFileName string = "gpt_conversation.json"
)

type QuesAnsw struct {
	Q string `koanf,json:"question"` // question
	A string `koanf,json:"answer"`   // answer
}

type ConversationSaver struct {
	QAList []QuesAnsw `koanf,json:"qa_list"`
	Prompt string     `koanf,json:"prompt"`
}

type Conversation struct {
	Context []QuesAnsw
	History []QuesAnsw
	Current *QuesAnsw
	Saver   *ConversationSaver `koanf,json:"conversation"`
	Tokens  int
	CNF     *config.Config
	Cursor  int
	path    string
}

func NewConversation(cnf *config.Config) (conv *Conversation) {
	conv = &Conversation{
		Context: []QuesAnsw{},
		History: []QuesAnsw{},
		CNF:     cnf,
		Saver: &ConversationSaver{
			QAList: []QuesAnsw{},
		},
		path: filepath.Join(cnf.GetWorkDir(), ConversationFileName),
	}
	return
}

func (that *Conversation) AddQuestion(ques string) {
	that.Current = &QuesAnsw{
		Q: ques,
	}
	that.Tokens = 0
	that.ResetCursor()
}

func (that *Conversation) AddAnswer(answ string, completed bool) {
	if that.Current == nil {
		return
	}
	that.Current.A += answ
	if completed {
		that.Context = append(that.Context, *that.Current)
		that.Tokens = 0
		if len(that.Context) > that.CNF.OpenAI.ContextLen {
			that.History = append(that.History, that.Context[0])
			that.Context = that.Context[1:]
		}
		that.Current = nil
	}
}

func (that *Conversation) GetMessages() []openai.ChatCompletionMessage {
	messages := make([]openai.ChatCompletionMessage, 0, 2*len(that.Context)+2)
	messages = append(
		messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: that.CNF.OpenAI.PromptStr,
		},
	)
	for _, c := range that.Context {
		messages = append(
			messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: c.Q,
			},
		)
		messages = append(
			messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: c.A,
			},
		)
	}
	if that.Current != nil {
		messages = append(
			messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: that.Current.Q,
			},
		)
	}
	return messages
}

func (that *Conversation) GetTokens() int {
	if that.Tokens == 0 {
		that.Tokens = NumTokensFromMessages(that.GetMessages(), that.CNF.OpenAI.Model)
	}
	return that.Tokens
}

func (that *Conversation) ClearContext() {
	that.History = append(that.History, that.Context...)
	that.Context = nil
	that.Tokens = 0
	that.ResetCursor()
}

func (that *Conversation) Len() int {
	l := len(that.History) + len(that.Context)
	if that.Current != nil {
		l++
	}
	return l
}

func (that *Conversation) ResetCursor() {
	that.Cursor = that.Len() - 1
}

func (that *Conversation) GetQAByCursor() QuesAnsw {
	if that.Cursor < len(that.History) {
		return that.History[that.Cursor]
	}

	if that.Current == nil {
		return that.Context[that.Cursor-len(that.History)]
	} else if that.Cursor > 0 {
		return that.Context[that.Cursor-len(that.History)-1]
	} else {
		return *that.Current
	}
}

func (that *Conversation) GetPrevQA() QuesAnsw {
	that.Cursor--
	if that.Cursor < 0 {
		that.ResetCursor()
	}
	return that.GetQAByCursor()
}

func (that *Conversation) GetNextQA() QuesAnsw {
	that.Cursor++
	if that.Cursor > that.Len()-1 {
		that.Cursor = 0
	}
	return that.GetQAByCursor()
}

func (that *Conversation) Save() {
	that.Saver.QAList = append(that.Saver.QAList, that.History...)
	that.Saver.QAList = append(that.Saver.QAList, that.Context...)
	that.Saver.Prompt = that.CNF.OpenAI.PromptStr
	if k, err := koanfer.NewKoanfer(that.path); err == nil {
		k.Save(that.Saver)
	}
}

func (that *Conversation) Load() {
	if k, err := koanfer.NewKoanfer(that.path); err == nil {
		err = k.Load(that.Saver)
		if err != nil {
			return
		}
		that.CNF.OpenAI.PromptStr = that.Saver.Prompt
		total := len(that.Saver.QAList)
		if total > that.CNF.OpenAI.ContextLen {
			that.Context = that.Saver.QAList[total-that.CNF.OpenAI.ContextLen:]
			that.History = that.Saver.QAList[:total-that.CNF.OpenAI.ContextLen]
		} else {
			that.Context = that.Saver.QAList
		}
	}
}
