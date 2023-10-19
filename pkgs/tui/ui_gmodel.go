package tui

import (
	"github.com/moqsien/goutils/pkgs/gtea/selector"
	openai "github.com/sashabaranov/go-openai"
)

/*
Select ChaGPT Model
*/
func GetModelSelector() ExtraModel {
	models := []string{
		openai.GPT4,
		openai.GPT432K0613,
		openai.GPT432K0314,
		openai.GPT432K,
		openai.GPT40613,
		openai.GPT40314,
		openai.GPT3Dot5Turbo,
		openai.GPT3Dot5Turbo0613,
		openai.GPT3Dot5Turbo0301,
		openai.GPT3Dot5Turbo16K,
		openai.GPT3Dot5Turbo16K0613,
		openai.GPT3Dot5TurboInstruct,
		openai.GPT3Davinci,
		openai.GPT3Davinci002,
		openai.GPT3Curie,
		openai.GPT3Curie002,
		openai.GPT3Ada,
		openai.GPT3Ada002,
		openai.GPT3Babbage,
		openai.GPT3Babbage002,
	}
	selectorItems := selector.NewItemList()
	for _, model := range models {
		selectorItems.Add(model, model)
	}
	sel := selector.NewSelectorModel(
		selectorItems.Keys(),
		selector.WithTitle("Choose model"),
		selector.WithHeight(20),
		selector.WithWidth(100),
		selector.WidthEnableMulti(false),
		selector.WithEnbleInfinite(true),
		selector.WithFilteringEnabled(false),
	)
	return sel
}
