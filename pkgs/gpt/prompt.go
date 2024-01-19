package gpt

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/moqsien/gogpt/pkgs/config"
	"github.com/moqsien/goutils/pkgs/gutils"
	"github.com/moqsien/goutils/pkgs/request"
)

const (
	PromptFileName string = "prompt.json"
)

type PromptItem struct {
	Title string `json:"act"`
	Msg   string `json:"prompt"`
}

type GPTPrompt struct {
	PromptList *[]PromptItem
	CNF        *config.Config
	prompt     string
	path       string
}

func NewGPTPrompt(cnf *config.Config) (gp *GPTPrompt) {
	gp = &GPTPrompt{CNF: cnf, path: filepath.Join(cnf.GetWorkDir(), PromptFileName)}
	gp.PromptList = &([]PromptItem{})
	gp.initiate()
	return
}

func (that *GPTPrompt) initiate() {
	if ok, _ := gutils.PathIsExist(that.path); !ok {
		that.DownloadPrompt()
	}
	if ok, _ := gutils.PathIsExist(that.path); ok {
		content, _ := os.ReadFile(that.path)
		json.Unmarshal(content, that.PromptList)
	}
}

func (that *GPTPrompt) DownloadPrompt() {
	f := request.NewFetcher()
	f.SetUrl(that.CNF.OpenAI.PromptMsgUrl)
	f.Timeout = 10 * time.Second
	size := f.GetFile(that.path, true)
	if size == 0 {
		fmt.Println("download failed")
	}
}

func (that *GPTPrompt) PromptStr() string {
	if that.prompt == "" {
		that.prompt = "You are ChatGPT, a large language model trained by OpenAI. Answer as concisely as possible."
	}
	return that.prompt
}

func (that *GPTPrompt) SetPrompt(prompt string) {
	that.prompt = prompt
}

func (that *GPTPrompt) GetPromptByTile(title string) (p string) {
	for _, pItem := range *that.PromptList {
		if pItem.Title == title {
			return pItem.Msg
		}
	}
	that.prompt = "You are ChatGPT, a large language model trained by OpenAI. Answer as concisely as possible."
	return that.prompt
}

func (that *GPTPrompt) GetTitleByPrompt(prompt string) (t string) {
	for _, pItem := range *that.PromptList {
		if pItem.Msg == prompt {
			return pItem.Title
		}
	}
	return
}
