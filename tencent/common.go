package tencent

const (
	ChatMessageRoleUser      = "user"
	ChatMessageRoleAssistant = "assistant"
)

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest represents a request structure for chat completion API.
type ChatCompletionRequest struct {
	Messages    []ChatCompletionMessage `json:"messages"`
	Temperature float64                 `json:"temperature,omitempty"`
	TopP        float64                 `json:"top_p,omitempty"`
	QueryId     string                  `json:"query_id,omitempty"`
	Stream      int                     `json:"stream"`
}

// ChatCompletionFullRequest represents a request structure for chat completion API.
type ChatCompletionFullRequest struct {
	AppId     int64  `json:"app_id"`
	SecretId  string `json:"secret_id"`
	Timestamp int64  `json:"timestamp"`
	Expired   int64  `json:"expired"`

	ChatCompletionRequest
}

// ChatCompletionResponse represents a response structure for chat completion API.
type ChatCompletionResponse struct {
	Choices []ChatCompletionChoices `json:"choices,omitempty"` // 结果
	Created string                  `json:"created,omitempty"` //unix 时间戳的字符串
	ID      string                  `json:"id,omitempty"`      //会话 id
	Usage   Usage                   `json:"usage,omitempty"`   //token 数量
	Error   ResponseError           `json:"error,omitempty"`   //错误信息 注意：此字段可能返回 null，表示取不到有效值
	Note    string                  `json:"note,omitempty"`    //注释
	ReqID   string                  `json:"req_id,omitempty"`  //唯一请求 ID，每次请求都会返回。用于反馈接口入参
}

type ChatCompletionStreamChoiceDelta struct {
	Content string `json:"content,omitempty"`
}

type ChatCompletionChoices struct {
	Messages     ChatCompletionMessage           `json:"messages,omitempty"`      // 内容，同步模式返回内容，流模式为 null 输出 content 内容总数最多支持 1024token。
	Delta        ChatCompletionStreamChoiceDelta `json:"delta,omitempty"`         // 内容，流模式返回内容，同步模式为 null 输出 content 内容总数最多支持 1024token。
	FinishReason string                          `json:"finish_reason,omitempty"` // 流式结束标志位，为 stop 则表示尾包
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
