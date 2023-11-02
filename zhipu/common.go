package zhipu

const (
	ChatMessageRoleUser      = "user"
	ChatMessageRoleAssistant = "assistant"

	ModelChatGLMPro  = "chatglm_pro"
	ModelChatGLMStd  = "chatglm_std"
	ModelChatGLMLite = "chatglm_lite"

	ModelTurbo = "chatglm_turbo"
)

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRef struct {
	Enable      bool   `json:"enable"`
	SearchQuery string `json:"search_query"`
}

// ChatCompletionRequest represents a request structure for chat completion API.
type ChatCompletionRequest struct {
	Prompt      []ChatCompletionMessage `json:"prompt"`
	Temperature float32                 `json:"temperature,omitempty"`
	TopP        float32                 `json:"top_p,omitempty"`
	RequestId   string                  `json:"request_id"`
	Incremental bool                    `json:"incremental"`
	ReturnType  string                  `json:"return_type,omitempty"`
	Ref         ChatCompletionRef       `json:"ref"`
}

// ChatCompletionResponse represents a response structure for chat completion API.
type ChatCompletionResponse struct {
	Id    string                      `json:"id"`
	Event string                      `json:"event"`
	Data  string                      `json:"data"`
	Meta  *ChatCompletionResponseMeta `json:"meta"`

	ErrorCode int    `json:"code"`
	ErrorMsg  string `json:"msg"`
	Success   bool   `json:"success"`
}

type ChatCompletionResponseMeta struct {
	Usage Usage `json:"usage"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
