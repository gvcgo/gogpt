package conversation

import (
	"strings"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	tiktoken "github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
)

/*
https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb
*/

func NumTokensFromMessages(messages []openai.ChatCompletionMessage, model string) (numTokens int) {
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		gprint.PrintError("encoding for model: %+v", err)
		return
	}

	var tokensPerMessage, tokensPerName int
	switch model {
	case openai.GPT3Dot5Turbo0613, openai.GPT3Dot5Turbo16K0613, openai.GPT40314,
		openai.GPT432K0314, openai.GPT40613, openai.GPT432K0613:
		tokensPerMessage = 3
		tokensPerName = 1
	case openai.GPT3Dot5Turbo0301:
		// every message follows <|start|>{role/name}\n{content}<|end|>\n
		tokensPerMessage = 4
		// if there's a name, the role is omitted
		tokensPerName = -1
	case openai.GPT4:
		return NumTokensFromMessages(messages, openai.GPT40613)
	default:
		if strings.Contains(model, openai.GPT3Dot5Turbo) {
			// gprint.PrintWarning(
			// 	"%s may update over time. Returning num tokens assuming %s.",
			// 	openai.GPT3Dot5Turbo,
			// 	openai.GPT3Dot5Turbo0613,
			// )
			return NumTokensFromMessages(messages, openai.GPT3Dot5Turbo0613)
		} else if strings.Contains(model, openai.GPT4) {
			// gprint.PrintWarning(
			// 	"%s may update over time. Returning num tokens assuming %s.",
			// 	openai.GPT4,
			// 	openai.GPT40613,
			// )
			return NumTokensFromMessages(messages, openai.GPT40613)
		} else {
			// gprint.PrintError(
			// 	"num_tokens_from_messages() is not implemented for model %s.",
			// 	model,
			// )
			// gprint.PrintInfo(
			// 	"See %s for information on how messages are converted to tokens.",
			// 	"https://github.com/openai/openai-python/blob/main/chatml.md",
			// )
			return
		}
	}

	for _, message := range messages {
		numTokens += tokensPerMessage
		numTokens += len(tkm.Encode(message.Content, nil, nil))
		numTokens += len(tkm.Encode(message.Role, nil, nil))
		numTokens += len(tkm.Encode(message.Name, nil, nil))
		if message.Name != "" {
			numTokens += tokensPerName
		}
	}
	// every reply is primed with <|start|>assistant<|message|>
	numTokens += 3
	return numTokens
}
