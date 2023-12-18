package minmax

const (
	ChatMessageRoleUser      = "USER"
	ChatMessageRoleAssistant = "BOT"
	ChatMessageRoleBot       = "BOT"

	ModelAbab55Chat = "abab5.5-chat"
	ModelAbab5Chat  = "abab5-chat"
)

type ChatCompletionMessage struct {
	SenderType string `json:"sender_type"`
	Text       string `json:"text"`
}

// ChatCompletionRequest represents a request structure for chat completion API.
type ChatCompletionRequest struct {
	Model               string                        `json:"model"`
	Prompt              string                        `json:"prompt,omitempty"`
	Messages            []ChatCompletionMessage       `json:"messages"`
	RoleMeta            ChatCompletionRequestRoleMeta `json:"role_meta"`
	ContinueLastMessage bool                          `json:"continue_last_message"`
	Temperature         float32                       `json:"temperature,omitempty"`
	TokensToGenerate    int                           `json:"tokens_to_generate,omitempty"`
	TopP                float32                       `json:"top_p,omitempty"`
	BeamWidth           int                           `json:"beam_width,omitempty"`
	Stream              bool                          `json:"stream,omitempty"`
	UseStandardSse      bool                          `json:"use_standard_sse"`
}
type ChatCompletionRequestRoleMeta struct {
	UserName string `json:"user_name"`
	BotName  string `json:"bot_name"`
}

// ChatCompletionResponse represents a response structure for chat completion API.
type ChatCompletionResponse struct {
	Id                  string                         `json:"id"`
	Created             int64                          `json:"created"`
	Model               string                         `json:"model"`
	Reply               string                         `json:"reply"`
	InputSensitive      bool                           `json:"input_sensitive"`
	InputSensitiveType  string                         `json:"input_sensitive_type"`
	OutputSensitive     bool                           `json:"output_sensitive"`
	OutputSensitiveType string                         `json:"output_sensitive_type"`
	Choices             []ChatCompletionResponseChoice `json:"choices"`
	Usage               Usage                          `json:"usage"`
	BaseResp            ChatCompletionResponseBaseResp `json:"base_resp"`
}

type ChatCompletionResponseChoice struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	FinishReason string `json:"finish_reason"`
	Delta        string `json:"delta"`
}

type Usage struct {
	//PromptTokens     int `json:"prompt_tokens"`
	//CompletionTokens int `json:"completion_tokens"`
	TotalTokens int `json:"total_tokens"`
}

type ChatCompletionResponseBaseResp struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}
