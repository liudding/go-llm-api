package sense

const (
	ChatMessageRoleUser      = "user"
	ChatMessageRoleAssistant = "assistant"

	ModelNovaPtcXlV1 = "nova-ptc-xl-v1"
	ModelNovaPtcSV2  = "nova-ptc-s-v2"
	ModelNovaPtcXsV1 = "nova-ptc-xs-v1"
)

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest represents a request structure for chat completion API.
type ChatCompletionRequest struct {
	Model             string                  `json:"model"`
	Messages          []ChatCompletionMessage `json:"messages"`
	N                 int                     `json:"n,omitempty"` // 生成回复数量，响应参数中的index即为回复序号（在使用某些模型时不支持传入该参数，详情请参考模型清单） https://platform.sensenova.cn/#/doc?path=/model.md
	MaxNewTokens      int                     `json:"max_new_tokens,omitempty"`
	RepetitionPenalty float32                 `json:"repetition_penalty,omitempty"`
	Temperature       float32                 `json:"temperature,omitempty"`
	TopP              float32                 `json:"top_p,omitempty"`
	Stream            bool                    `json:"stream,omitempty"`
	User              string                  `json:"user,omitempty"`
}

// ChatCompletionResponse represents a response structure for chat completion API.
type ChatCompletionResponse struct {
	Data  ChatCompletionResponseData  `json:"data"`
	Error ChatCompletionResponseError `json:"error"`
}

type ChatCompletionResponseData struct {
	Id      string                          `json:"id"`
	Choices []ChatCompletionResponseChoices `json:"choices"`
	Usage   Usage                           `json:"usage"`

	ErrorCode int    `json:"code"`
	ErrorMsg  string `json:"msg"`
	Success   bool   `json:"success"`
}

type ChatCompletionResponseChoices struct {
	Message      string `json:"message"`
	Delta        string `json:"delta"`
	FinishReason string `json:"finish_reason"`
	Index        int    `json:"index"`
	Role         string `json:"role"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatCompletionResponseError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
	Details []any  `json:"details"`
}
