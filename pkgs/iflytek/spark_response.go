package iflytek

import (
	"fmt"
	"io"

	"github.com/gogf/gf/encoding/gjson"
)

/*
响应：

	{
	    "header":{
	        "code":0,
	        "message":"Success",
	        "sid":"cht000cb087@dx18793cd421fb894542",
	        "status":2
	    },
	    "payload":{
	        "choices":{
	            "status":2,
	            "seq":0,
	            "text":[
	                {
	                    "content":"我可以帮助你的吗？",
	                    "role":"assistant",
	                    "index":0
	                }
	            ]
	        },
	        "usage":{
	            "text":{
	                "question_tokens":4,
	                "prompt_tokens":5,
	                "completion_tokens":9,
	                "total_tokens":14
	            }
	        }
	    }
	}

字段解释
code	int	    错误码，0表示正常，非0表示出错；详细释义可在接口说明文档最后的错误码说明了解
message	string	会话是否成功的描述信息
sid	    string	会话的唯一id，用于讯飞技术人员查询服务端会话日志使用,出现调用错误时建议留存该字段
status	int	    会话状态，取值为[0,1,2]；0代表首次结果；1代表中间结果；2代表最后一个结果

status	int	    文本响应状态，取值为[0,1,2]; 0代表首个文本结果；1代表中间文本结果；2代表最后一个文本结果
seq	    int	    返回的数据序号，取值为[0,9999999]
content	string	AI的回答内容
role	string	角色标识，固定为assistant，标识角色为AI
index	int	    结果序号，取值为[0,10]; 当前为保留字段，开发者可忽略

question_tokens	   int	保留字段，可忽略
prompt_tokens	   int	包含历史问题的总tokens大小
completion_tokens  int	回答的tokens大小
total_tokens	   int	prompt_tokens和completion_tokens的和，也是本次交互计费的tokens大小

错误码
10000	升级为ws出现错误
10001	通过ws读取用户的消息出错
10002	通过ws向用户发送消息 错
10003	用户的消息格式有错误
10004	用户数据的schema错误
10005	用户参数值有错误
10006	用户并发错误：当前用户已连接，同一用户不能多处同时连接。
10007	用户流量受限：服务正在处理用户当前的问题，需等待处理完成后再发送新的请求。（必须要等大模型完全回复之后，才能发送下一个问题）
10008	服务容量不足，联系工作人员
10009	和引擎建立连接失败
10010	接收引擎数据的错误
10011	发送数据给引擎的错误
10012	引擎内部错误
10013	输入内容审核不通过，涉嫌违规，请重新调整输入内容
10014	输出内容涉及敏感信息，审核不通过，后续结果无法展示给用户
10015	appid在黑名单中
10016	appid授权类的错误。比如：未开通此功能，未开通对应版本，token不足，并发超过授权 等等
10017	清除历史失败
10019	表示本次会话内容有涉及违规信息的倾向；建议开发者收到此错误码后给用户一个输入涉及违规的提示
10110	服务忙，请稍后再试
10163	请求引擎的参数异常 引擎的schema 检查不通过
10222	引擎网络异常
10907	token数量超过上限。对话历史+问题的字数太多，需要精简输入
11200	授权错误：该appId没有相关功能的授权 或者 业务量超过限制
11201	授权错误：日流控超限。超过当日最大访问量的限制
11202	授权错误：秒级流控超限。秒级并发超过授权路数限制
11203	授权错误：并发流控超限。并发路数超过授权路数限制
*/

type SparkAPIError struct {
	Code int
	Info string
}

func (that SparkAPIError) Error() string {
	return fmt.Sprintf("code: %d, info: %s", that.Code, that.Info)
}

func NewSparkError(code int, info string) (sae SparkAPIError) {
	return SparkAPIError{Code: code, Info: info}
}

var (
	ErrUpgradeToWebsocketFailed = NewSparkError(10000, "upgrade to websocket failed")
	ErrReadMessageFailed        = NewSparkError(10001, "read message from user failed")
	ErrSendMessageFailed        = NewSparkError(10002, "send message to user failed")
	ErrMessageFormtIncorrect    = NewSparkError(10003, "incorrect format user message")
	ErrSchemaIncorrect          = NewSparkError(10004, "incorrect schema from user message")
	ErrParamsIncorrect          = NewSparkError(10005, "incorrect params from user")
	ErrConcurrency              = NewSparkError(10006, "connect already exist")
	ErrNetworkFlow              = NewSparkError(10007, "network flow error")
	ErrServiceCapacity          = NewSparkError(10008, "capacity not enough")
	ErrConnectToEngineFailed    = NewSparkError(10009, "failed to connect to engine")
	ErrRecieveDataFromEngine    = NewSparkError(10010, "failed to recieve data from engine")
	ErrSendDataToEngineFailed   = NewSparkError(10011, "failed to send data to engine")
	ErrInEngine                 = NewSparkError(10012, "engine errored")
	ErrIllegalMessage           = NewSparkError(10013, "illegal message")
	ErrSensitiveMessage         = NewSparkError(10014, "sensitive message")
	ErrAppIDInBlacklist         = NewSparkError(10015, "appID in blacklist")
	ErrAppIDAuthentication      = NewSparkError(10016, "authentication error")
	ErrClearHistoryFailed       = NewSparkError(10017, "clear history failed")
	ErrIllegalMessageTendency   = NewSparkError(10019, "illegal message tendency")
	ErrServerBusy               = NewSparkError(10110, "server is busy")
	ErrIncorrectParamForEngine  = NewSparkError(10163, "incorrect param for engine")
	ErrEngineException          = NewSparkError(10222, "engine exceptions")
	ErrReachMaxTokens           = NewSparkError(10907, "exceed max tokens")
	ErrNoAuth                   = NewSparkError(11200, "no authentication")
	ErrExceedDailyReqLimit      = NewSparkError(11201, "reach daily request limit")
	ErrExceedSecondReqLimit     = NewSparkError(11202, "reach QPS limit")
	ErrExceedConcurrencyLimit   = NewSparkError(11203, "reach concurrency limit")
)

var SparkErrorMap map[int]error = map[int]error{
	10000: ErrUpgradeToWebsocketFailed,
	10001: ErrReadMessageFailed,
	10002: ErrSendMessageFailed,
	10003: ErrMessageFormtIncorrect,
	10004: ErrSchemaIncorrect,
	10005: ErrParamsIncorrect,
	10006: ErrConcurrency,
	10007: ErrNetworkFlow,
	10008: ErrServiceCapacity,
	10009: ErrConnectToEngineFailed,
	10010: ErrRecieveDataFromEngine,
	10011: ErrSendDataToEngineFailed,
	10012: ErrInEngine,
	10013: ErrIllegalMessage,
	10014: ErrSensitiveMessage,
	10015: ErrAppIDInBlacklist,
	10016: ErrAppIDAuthentication,
	10017: ErrClearHistoryFailed,
	10019: ErrIllegalMessageTendency,
	10110: ErrServerBusy,
	10163: ErrIncorrectParamForEngine,
	10222: ErrEngineException,
	10907: ErrReachMaxTokens,
	11200: ErrNoAuth,
	11201: ErrExceedDailyReqLimit,
	11202: ErrExceedSecondReqLimit,
	11203: ErrExceedConcurrencyLimit,
}

type ResponseMsg struct {
	Content string `json:"content"`
	Role    string `json:"role"`
	Index   int    `json:"index"`
}

type SparkResponse struct {
	Raw             map[string]interface{}
	ErrCode         int
	Error           error
	ChoiceStatus    int
	TotalTokens     int64
	ResponseMsgList []ResponseMsg
}

func NewSparkResponse(raw map[string]interface{}) (sr *SparkResponse) {
	sr = &SparkResponse{Raw: raw, ResponseMsgList: []ResponseMsg{}, Error: nil}
	return
}

func (that *SparkResponse) Parse() {
	// fmt.Println(string(that.Raw))
	j := gjson.New(that.Raw)
	that.ErrCode = j.GetInt("header.code")
	that.Error = SparkErrorMap[that.ErrCode]
	if that.ErrCode != 0 {
		return
	}
	that.ChoiceStatus = j.GetInt("payload.choices.status")
	if that.ChoiceStatus == 2 {
		that.Error = io.EOF
		that.TotalTokens = j.GetInt64("payload.usage.text.total_tokens")
	}
	text := j.GetArray("payload.choices.text")
	for _, m := range text {
		msg := m.(map[string]interface{})
		respMsg := ResponseMsg{}
		respMsg.Content = msg["content"].(string)
		respMsg.Role = msg["role"].(string)
		respMsg.Index = int(msg["index"].(float64))
		if respMsg.Content != "" {
			that.ResponseMsgList = append(that.ResponseMsgList, respMsg)
		}
	}
}
