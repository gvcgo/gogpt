package gpt

import (
	"github.com/moqsien/gogpt/pkgs/config"
)

/*
Manage Chatgpt conversation
*/

type QuesAnsw struct {
	Q string // question
	A string // answer
}

type Conversation struct {
	Context []QuesAnsw
	History []QuesAnsw
	Current *QuesAnsw
	Tokens  int
	CNF     *config.Config
}

func NewConversation(cnf *config.Config) (conv *Conversation) {
	conv = &Conversation{
		Context: []QuesAnsw{},
		History: []QuesAnsw{},
		CNF:     cnf,
	}
	return
}

func (that *Conversation) AddQuestion(ques string) {
	that.Current = &QuesAnsw{
		Q: ques,
	}
	that.Tokens = 0
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
