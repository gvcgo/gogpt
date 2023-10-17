package gpt

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	retry "github.com/avast/retry-go"
	"github.com/moqsien/gogpt/pkgs/config"
	"github.com/sashabaranov/go-openai"
	nproxy "golang.org/x/net/proxy"
)

const (
	ProxyEnv string = "CHATGPT_PROXY"
)

type GPT struct {
	OpenAIClient openai.Client
	Stream       *openai.ChatCompletionStream
	CNF          *config.Config
}

func NewGPT(cnf *config.Config) (g *GPT) {
	g = &GPT{
		CNF: cnf,
	}
	g.initiate()
	return
}

func (that *GPT) initiate() {
	var openaiConf openai.ClientConfig
	if that.CNF.OpenAI.ApiType == openai.APITypeOpenAI {
		openaiConf = openai.DefaultConfig(that.CNF.OpenAI.ApiKey)
		if that.CNF.OpenAI.BaseUrl != "" {
			openaiConf.BaseURL = that.CNF.OpenAI.BaseUrl
		}
		if that.CNF.OpenAI.EmptyMessagesLimit != 0 {
			openaiConf.EmptyMessagesLimit = that.CNF.OpenAI.EmptyMessagesLimit
		}
	} else {
		openaiConf = openai.DefaultAzureConfig(that.CNF.OpenAI.ApiKey, that.CNF.OpenAI.BaseUrl)
		if that.CNF.OpenAI.OrgID != "" {
			openaiConf.OrgID = that.CNF.OpenAI.OrgID
		}
		if that.CNF.OpenAI.ApiVersion != "" {
			openaiConf.APIVersion = that.CNF.OpenAI.ApiVersion
		}
	}
	if that.CNF.OpenAI.EmptyMessagesLimit != 0 {
		openaiConf.EmptyMessagesLimit = that.CNF.OpenAI.EmptyMessagesLimit
	}
	openaiConf.HTTPClient = that.getHttpClient()
}

func (that *GPT) getHttpClient() (httpClient *http.Client) {
	tt := time.Duration(that.CNF.OpenAI.TimeOut) * time.Second
	scheme, host, port := that.parseProxy()
	httpClient = &http.Client{Timeout: tt}
	switch scheme {
	case "http", "https":
		pUrl, err := url.Parse(that.CNF.OpenAI.Proxy)
		if err != nil {
			return
		}
		httpClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(pUrl),
		}
	case "socks5":
		if dialer, err := nproxy.SOCKS5("tcp", fmt.Sprintf("%s:%d", host, port), nil, nproxy.Direct); err == nil {
			httpClient.Transport = &http.Transport{Dial: dialer.Dial}
		}
	default:
	}
	return
}

func (that *GPT) parseProxy() (scheme, host string, port int) {
	p := that.CNF.OpenAI.Proxy
	if p == "" {
		p = os.Getenv(ProxyEnv)
	}
	if p == "" {
		return
	}
	if u, err := url.Parse(p); err == nil {
		scheme = u.Scheme
		host = u.Hostname()
		port, _ = strconv.Atoi(u.Port())
		if port == 0 {
			port = 80
		}
	}
	return
}

func (that *GPT) SendMsg(msgs []openai.ChatCompletionMessage) (m string, err error) {
	err = retry.Do(
		func() error {
			req := openai.ChatCompletionRequest{
				Model:       that.CNF.OpenAI.Model,
				Messages:    msgs,
				MaxTokens:   that.CNF.OpenAI.MaxTokens,
				Temperature: that.CNF.OpenAI.Temperature,
				N:           1,
			}
			that.Stream, err = that.OpenAIClient.CreateChatCompletionStream(context.Background(), req)
			if err != nil {
				that.Stream = nil
				return err
			}
			resp, err := that.Stream.Recv()
			if err != nil {
				return err
			}
			m = resp.Choices[0].Delta.Content
			return nil
		},
		retry.Attempts(3),
		retry.LastErrorOnly(true),
	)
	if err != nil {
		return "", err
	}
	return
}

func (that *GPT) RecvMsg() (m string, err error) {
	if that.Stream == nil {
		return "", fmt.Errorf("no stream found")
	}
	resp, err := that.Stream.Recv()
	if err != nil {
		return "", err
	}
	m = resp.Choices[0].Delta.Content
	return
}

func (that *GPT) Close() {
	if that.Stream != nil {
		that.Stream.Close()
	}
	that.Stream = nil
}
