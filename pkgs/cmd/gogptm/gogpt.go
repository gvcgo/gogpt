package main

import (
	"os"
	"path/filepath"

	"github.com/gvcgo/goutils/pkgs/gtea/gprint"
	"github.com/gvcgo/goutils/pkgs/gutils"
	"github.com/gvcgo/gogpt/pkgs/config"
	"github.com/gvcgo/gogpt/pkgs/gpt"
	"github.com/gvcgo/gogpt/pkgs/tui"
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
