package tui

import (
	"github.com/moqsien/gogpt/pkgs/config"
	"github.com/moqsien/gogpt/pkgs/gpt"
	"github.com/moqsien/goutils/pkgs/gtea/selector"
)

/*
Choose a pormpt
*/
func GetPromptModel(cnf *config.Config) ExtraModel {
	p := gpt.NewGPTPrompt(cnf)
	itemList := selector.NewItemList()
	for _, item := range *p.PromptList {
		itemList.Add(item.Title, item.Msg)
	}
	sel := selector.NewSelectorModel(
		itemList.Keys(),
		selector.WithTitle("Choose a prompt"),
		selector.WidthEnableMulti(false),
		selector.WithEnbleInfinite(true),
		selector.WithWidth(100),
		selector.WithHeight(40),
		selector.WithShowStatusBar(true),
	)
	return sel
}
