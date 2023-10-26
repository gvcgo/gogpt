package main

import (
	"os"
	"path/filepath"

	"github.com/moqsien/gogpt/pkgs/config"
	"github.com/moqsien/gogpt/pkgs/gpt"
	"github.com/moqsien/gogpt/pkgs/tui"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/gutils"
	"github.com/postfinance/single"
)

func main() {
	lockFile, _ := single.New("chatgpt")
	if err := lockFile.Lock(); err != nil {
		gprint.PrintError("Another gogpt program is running: %s", lockFile.Lockfile())
		os.Exit(1)
	}
	defer func() {
		lockFile.Unlock()
	}()

	cnf := tui.GetDefaultConfig()
	cnf.OpenAI.PromptMsgUrl = config.PromptUrl
	promptPath := filepath.Join(cnf.GetWorkDir(), gpt.PromptFileName)
	if ok, _ := gutils.PathIsExist(promptPath); !ok {
		prompt := gpt.NewGPTPrompt(cnf)
		prompt.DownloadPrompt()
	}
	ui := tui.NewGPTUI(cnf)
	ui.Run()
}
