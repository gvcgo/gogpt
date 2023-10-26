package iflytek

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/avast/retry-go"
	"github.com/moqsien/gogpt/pkgs/config"
	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/sashabaranov/go-openai"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

/*
	 请求：
		{
		        "header": {
		            "app_id": "12345",
		            "uid": "12345"
		        },
		        "parameter": {
		            "chat": {
		                "domain": "general",
		                "temperature": 0.5,
		                "max_tokens": 1024,
		            }
		        },
		        "payload": {
		            "message": {
		                # 如果想获取结合上下文的回答，需要开发者每次将历史问答信息一起传给服务端，如下示例
		                # 注意：text里面的所有content内容加一起的tokens需要控制在8192以内，开发者如有较长对话需求，需要适当裁剪历史信息
		                "text": [
		                    {"role": "user", "content": "你是谁"} # 用户的历史问题
		                    {"role": "assistant", "content": "....."}  # AI的历史回答结果
		                    # ....... 省略的历史对话
		                    {"role": "user", "content": "你会做什么"}  # 最新的一条问题，如无需上下文，可只传最新一条问题
		                ]
		        }
		    }
		}

# 参数
app_id	string	是		应用appid，从开放平台控制台创建的应用中获取
uid	    string	否	最大长度32	每个用户的id，用于区分不同用户

domain	    string	是	取值为[general,generalv2,generalv3]	指定访问的领域,general指向V1.5版本,generalv2指向V2版本,generalv3指向V3版本 。注意：不同的取值对应的url也不一样！
temperature	float	否	取值为[0,1],默认为0.5	核采样阈值。用于决定结果随机性，取值越高随机性越强即相同的问题得到的不同答案的可能性越高
max_tokens	int	    否	V1.5取值为[1,4096]，V2.0取值为[1,8192]。默认为2048	模型回答的tokens的最大长度
top_k	    int	    否	取值为[1，6],默认为4	从k个候选中随机选择⼀个（⾮等概率）
chat_id	    string	否	需要保障用户下的唯一性	用于关联用户会话

role	string	是	取值为[user,assistant]	user表示是用户的问题，assistant表示AI的回复
content	string	是	所有content的累计tokens需控制8192以内	用户和AI的对话内容
*/

var RoleMap map[string]string = map[string]string{
	openai.ChatMessageRoleAssistant: "assistant",
	openai.ChatMessageRoleUser:      "user",
}

type RequestData map[string]interface{}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Spark struct {
	CNF         *config.Config
	Conn        *websocket.Conn
	AuthUrl     string
	sparkDomain string
	hostUrl     string
	token       int64
}

func NewSpark(cnf *config.Config) (s *Spark) {
	s = &Spark{
		CNF:   cnf,
		token: 0,
	}
	s.AssembleAuthUrl()
	return
}

func (that *Spark) AssembleAuthUrl() {
	switch that.CNF.Spark.APIVersion {
	case config.SparkAPIV1:
		that.hostUrl = config.SparkAPIV1Dot1
		that.sparkDomain = config.SparkDomainV1
	case config.SparkAPIV2:
		that.hostUrl = config.SparkAPIV2Dot1
		that.sparkDomain = config.SparkDomainV2
	case config.SparkAPIV3:
		that.hostUrl = config.SparkAPIV3Dot1
		that.sparkDomain = config.SparkDomainV3
	default:
		that.hostUrl = config.SparkAPIV1Dot1
		that.sparkDomain = config.SparkDomainV1
	}

	ul, err := url.Parse(that.hostUrl)
	if err != nil {
		gprint.PrintError("parse spark url failed: %+v", err)
		os.Exit(1)
	}
	//签名时间 "Tue, 28 May 2019 09:10:42 MST"
	date := time.Now().UTC().Format(time.RFC1123)

	//参与签名的字段 host ,date, request-line
	signString := []string{"host: " + ul.Host, "date: " + date, "GET " + ul.Path + " HTTP/1.1"}

	//拼接签名字符串
	sgin := strings.Join(signString, "\n")

	//签名结果
	sha := that.HmacWithShaTobase64("hmac-sha256", sgin, that.CNF.Spark.APPSecrete)

	//构建请求参数 此时不需要urlencoding
	authUrl := fmt.Sprintf(
		"hmac username=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"",
		that.CNF.Spark.APPKey,
		"hmac-sha256", "host date request-line",
		sha,
	)

	//将请求参数使用base64编码
	authorization := base64.StdEncoding.EncodeToString([]byte(authUrl))

	v := url.Values{}
	v.Add("host", ul.Host)
	v.Add("date", date)
	v.Add("authorization", authorization)

	//将编码后的字符串url encode后添加到url后面
	that.AuthUrl = that.hostUrl + "?" + v.Encode()
}

func (that *Spark) readResp(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("code=%d,body=%s", resp.StatusCode, string(b))
}

func (that *Spark) HmacWithShaTobase64(algorithm, data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}

func (that *Spark) Connect() {
	if that.CNF.Spark.Timeout == 0 {
		that.CNF.Spark.Timeout = 60
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(that.CNF.Spark.Timeout)*time.Second)
	defer cancel()
	if that.Conn != nil {
		// Spark v1.1 一次回答之后会自动关闭会话，从而导致继续使用原有Conn读写会出错
		// 所以这里先关闭本地Conn，然后重新连接。
		that.Conn.CloseNow()
		time.Sleep(2 * time.Second)
	}
	var (
		resp *http.Response
		err  error
	)

	err = retry.Do(
		func() error {
			that.Conn, resp, err = websocket.Dial(ctx, that.AuthUrl, nil)
			return err
		},
		retry.Attempts(3),
		retry.LastErrorOnly(true),
	)

	if err != nil {
		panic(that.readResp(resp) + err.Error())
	} else if resp.StatusCode != 101 {
		panic(that.readResp(resp) + err.Error())
	}
}

func (that *Spark) generateRequestData(msgs []openai.ChatCompletionMessage) RequestData {
	messages := []Message{}
	for _, m := range msgs {
		if role := RoleMap[m.Role]; role != "" {
			messages = append(messages, Message{
				Role:    role,
				Content: m.Content,
			})
		}
	}
	var (
		temperature float64 = 0.5
		topK        int64   = 4
		maxTokens   int64   = 2048
	)
	if that.CNF.Spark.Temperature != 0.0 {
		temperature = that.CNF.Spark.Temperature
	}
	if that.CNF.Spark.TopK != 0 {
		topK = that.CNF.Spark.TopK
	}

	if that.CNF.Spark.MaxTokens != 0 {
		maxTokens = that.CNF.Spark.MaxTokens
	}

	data := map[string]interface{}{
		"header": map[string]interface{}{
			"app_id": that.CNF.Spark.APPID,
		},
		"parameter": map[string]interface{}{
			"chat": map[string]interface{}{
				"domain":      that.sparkDomain,
				"temperature": temperature,
				"top_k":       topK,
				"max_tokens":  maxTokens,
			},
		},
		"payload": map[string]interface{}{
			"message": map[string]interface{}{
				"text": messages,
			},
		},
	}
	return data
}

func (that *Spark) SendMsg(msgs []openai.ChatCompletionMessage) (m string, err error) {
	that.Connect()
	reqData := that.generateRequestData(msgs)
	if that.Conn == nil {
		return
	}
	err = retry.Do(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
		defer cancel()
		return wsjson.Write(ctx, that.Conn, reqData)
	},
		retry.Attempts(3),
		retry.LastErrorOnly(true),
	)
	return "", err
}

func (that *Spark) RecvMsg() (m string, err error) {
	if that.Conn == nil {
		return "", fmt.Errorf("no conn found")
	}
	var msg map[string]interface{}
	err = wsjson.Read(context.Background(), that.Conn, &msg)
	if err != nil {
		return "", err
	}
	resp := NewSparkResponse(msg)
	resp.Parse()
	err = resp.Error
	that.token += resp.TotalTokens
	for _, r := range resp.ResponseMsgList {
		if r.Role == RoleMap[openai.ChatMessageRoleAssistant] {
			m += r.Content
		}
	}
	return
}

func (that *Spark) Close() {
	if that.Conn != nil {
		that.Conn.CloseNow()
		that.Conn = nil
	}
}

func (that *Spark) GetTokens() int64 {
	return that.token
}
