package gpt

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/moqsien/gogpt/pkgs/config"
	"github.com/sashabaranov/go-openai"
	nproxy "golang.org/x/net/proxy"
)

const (
	ProxyEnv string = "CHATGPT_PROXY"
)

type GPT struct {
	OpenAIClient openai.Client
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

func (that *GPT) SendMsg() {
}
