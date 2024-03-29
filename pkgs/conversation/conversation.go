package conversation

import (
	"path/filepath"

	"github.com/gvcgo/goutils/pkgs/koanfer"
	"github.com/gvcgo/gogpt/pkgs/config"
	"github.com/sashabaranov/go-openai"
)

/*
Manage Chatgpt conversation
*/
const (
	ConversationFileName string = "gpt_conversation.json"
	BotGPT               string = "ChatGPT"
	BotSpark             string = "Spark"
)

type QuesAnsw struct {
	Q string `koanf,json:"question"` // question
	A string `koanf,json:"answer"`   // answer
}

type ConversationSaver struct {
	QAList  []QuesAnsw `koanf,json:"qa_list"`
	Prompt  string     `koanf,json:"prompt"`
	BotType string     `koanf,json:"bot_type"`
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
	BotType string
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

func (that *Conversation) SetBotType(botType string) {
	that.BotType = botType
	that.ClearAll()
}

func (that *Conversation) ClearAll() {
	that.Context = []QuesAnsw{}
	that.History = []QuesAnsw{}
	that.Current = nil
	that.Saver = &ConversationSaver{
		QAList: []QuesAnsw{},
	}
	that.Tokens = 0
	that.Cursor = 0
}

func (that *Conversation) AddQuestion(ques string) {
	if that.Current == nil {
		that.Current = &QuesAnsw{
			Q: ques,
		}
	} else {
		that.Current.Q = ques
		that.Current.A = ""
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
	if that.BotType == BotGPT {
		return NumTokensFromMessages(that.GetMessages(), that.CNF.OpenAI.Model)
	}
	// tokens for Spark
	return that.Tokens
}

func (that *Conversation) AddTokens(tokens int64) int {
	that.Tokens += int(tokens)
	return that.Tokens
}

func (that *Conversation) ClearContext() {
	that.History = append(that.History, that.Context...)
	that.Context = []QuesAnsw{}
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
	if that.Cursor < 0 || that.Cursor > that.Len()-1 {
		that.Cursor = 0
	}
	if that.Len() == 0 {
		return QuesAnsw{}
	}
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
	return that.GetQAByCursor()
}

func (that *Conversation) GetNextQA() QuesAnsw {
	that.Cursor++
	return that.GetQAByCursor()
}

func (that *Conversation) Save() {
	that.Saver.QAList = append(that.Saver.QAList, that.History...)
	that.Saver.QAList = append(that.Saver.QAList, that.Context...)
	that.Saver.Prompt = that.CNF.OpenAI.PromptStr
	that.Saver.BotType = that.BotType
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
		that.BotType = that.Saver.BotType
	}
}

func (that *Conversation) ClearCurrentAnswer() {
	if that.Current != nil {
		that.Current.A = ""
	}
}
