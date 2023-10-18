package gpt

import (
	"github.com/moqsien/gogpt/pkgs/config"
)

type PromptItem struct {
	Title string `json:"act"`
	Msg   string `json:"prompt"`
}

type GPTPrompt struct {
	PromptList []PromptItem
	CNF        *config.Config
}

func NewGPTPrompt(cnf *config.Config) (gp *GPTPrompt) {
	gp = &GPTPrompt{CNF: cnf}
	gp.PromptList = []PromptItem{}
	gp.initiate()
	return
}

func (that *GPTPrompt) initiate() {

}

func (that *GPTPrompt) ChoosePrompt() {

}
