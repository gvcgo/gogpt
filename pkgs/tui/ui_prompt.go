package tui

import (
	"github.com/moqsien/gogpt/pkgs/gpt"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/gtea/selector"
)

/*
Choose a pormpt
*/
func GetPromptModel(prompt *gpt.GPTPrompt) ExtraModel {
	itemList := selector.NewItemList()
	for _, item := range *prompt.PromptList {
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

func SetGPTPrompt(prompt *gpt.GPTPrompt, values map[string]string) {
	if prompt == nil {
		gprint.PrintError("conf object is nil!")
		return
	}
	for title := range values {
		result := prompt.GetPromptByTile(title)
		prompt.SetPrompt(result)
		break
	}
}
