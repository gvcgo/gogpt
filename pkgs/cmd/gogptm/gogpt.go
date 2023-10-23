package main

import (
	"os"

	"github.com/moqsien/gogpt/pkgs/tui"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
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
	ui := tui.NewGPTUI(cnf)
	ui.Run()
}
