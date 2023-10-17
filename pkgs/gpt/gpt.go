package gpt

import (
	"github.com/sashabaranov/go-openai"
)

type GPT struct {
	OpenAIClient openai.Client
}
