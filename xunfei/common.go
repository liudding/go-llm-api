package xunfei

const (
	ChatMessageRoleUser      = "user"
	ChatMessageRoleAssistant = "assistant"

	parameterChatDomainGeneral   = "general"
	parameterChatDomainGeneralV2 = "generalv2"
	parameterChatDomainGeneralV3 = "generalv3"
)

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionStreamRequest struct {
	Header    ChatCompletionStreamRequestHeader    `json:"header"`
	Parameter ChatCompletionStreamRequestParameter `json:"parameter"`
	Payload   ChatCompletionStreamRequestPayload   `json:"payload"`
}

type ChatCompletionStreamRequestHeader struct {
	AppId string `json:"app_id"`
	Uid   string `json:"uid,omitempty"`
}

type ChatCompletionStreamRequestParameter struct {
	Chat ChatCompletionStreamRequestParameterChat `json:"chat"`
}

type ChatCompletionStreamRequestPayload struct {
	Message ChatCompletionStreamRequestPayloadMessage `json:"message"`
}

type ChatCompletionStreamRequestPayloadMessage struct {
	Text []ChatCompletionMessage `json:"text"`
}

type ChatCompletionStreamRequestParameterChat struct {
	Domain      string  `json:"domain"`
	Temperature float64 `json:"temperature,omitempty"`
	TopK        int     `json:"top_k,omitempty"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
	Auditing    string  `json:"auditing,omitempty"`
}

type ChatCompletionStreamResponse struct {
	Header  ChatCompletionStreamResponseHeader  `json:"header"`
	Payload ChatCompletionStreamResponsePayload `json:"payload"`
}

type ChatCompletionStreamResponseHeader struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Sid     string `json:"sid"`
	Status  int    `json:"status"` // 会话状态，取值为[0,1,2]；0代表首次结果；1代表中间结果；2代表最后一个结果
}

type ChatCompletionStreamResponsePayload struct {
	Choices ChatCompletionStreamResponsePayloadChoices `json:"choices"`
	Usage   Usage                                      `json:"usage"`
}

type ChatCompletionStreamResponsePayloadChoices struct {
	Status int `json:"status"` // 文本响应状态，取值为[0,1,2]; 0代表首个文本结果；1代表中间文本结果；2代表最后一个文本结果
	Seq    int `json:"seq"`    // 返回的数据序号，取值为[0,9999999]

	ChatCompletionMessage
	Text []ChatCompletionMessage `json:"text"` // 讯飞文档中有两种格式，这里兼容处理
}

// ChatCompletionRequest represents a request structure for chat completion API.
type ChatCompletionRequest struct {
	Messages    []ChatCompletionMessage `json:"messages"`
	Temperature float64                 `json:"temperature,omitempty"`
	TopK        int                     `json:"top_k,omitempty"`
	MaxTokens   int                     `json:"max_tokens,omitempty"`
	Auditing    string                  `json:"auditing,omitempty"`
	Uid         string                  `json:"uid,omitempty"`
}

// ChatCompletionResponse represents a response structure for chat completion API.
//type ChatCompletionResponse struct {
//	Choices []ChatCompletionChoices `json:"choices,omitempty"` // 结果
//	Usage   Usage                   `json:"usage,omitempty"`   // token 数量
//	Error   ResponseError           `json:"error,omitempty"`
//}

type ChatCompletionStreamChoiceDelta struct {
	Content string `json:"content,omitempty"`
}

type ChatCompletionChoices struct {
	Messages ChatCompletionMessage `json:"messages,omitempty"` // 内容，同步模式返回内容，流模式为 null 输出 content 内容总数最多支持 1024token。
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ResponseError :错误信息
type ResponseError struct {
	Message string `json:"message,omitempty"` // 错误提示信息
	Code    int    `json:"code,omitempty"`    // Code 错误码
}
