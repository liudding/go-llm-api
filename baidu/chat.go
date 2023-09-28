package baidu

import (
	"context"
	"errors"
	"net/http"
)

// Chat message role defined by the OpenAI API.
const (
	ChatMessageRoleUser      = "user"
	ChatMessageRoleAssistant = "assistant"
)

const chatCompletionsSuffix = "/chat/completions"

var (
	ErrChatCompletionInvalidModel       = errors.New("this model is not supported with this method, please use CreateCompletion client method instead") //nolint:lll
	ErrChatCompletionStreamNotSupported = errors.New("streaming is not supported with this method, please use CreateChatCompletionStream")              //nolint:lll
)

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest represents a request structure for chat completion API.
type ChatCompletionRequest struct {
	Messages     []ChatCompletionMessage `json:"messages"`
	Temperature  float32                 `json:"temperature,omitempty"`
	TopP         float32                 `json:"top_p,omitempty"`
	PenaltyScore float32                 `json:"penalty_score,omitempty"`
	Stream       bool                    `json:"stream,omitempty"`
	UserId       string                  `json:"user_id,omitempty"`
}

// ChatCompletionResponse represents a response structure for chat completion API.
type ChatCompletionResponse struct {
	ID               string `json:"id"`
	Object           string `json:"object"`
	Created          int64  `json:"created"`
	SentenceId       int    `json:"sentence_id"`
	IsEnd            bool   `json:"is_end"`
	IsTruncated      bool   `json:"is_truncated"`
	Result           string `json:"result"`
	NeedClearHistory bool   `json:"need_clear_history"` // 表示用户输入是否存在安全，是否关闭当前会话，清理历史会话信息
	BanRound         int    `json:"ban_round"`          // 当need_clear_history为true时，此字段会告知第几轮对话有敏感信息，如果是当前问题，ban_round=-1
	Usage            Usage  `json:"usage"`

	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

// CreateChatCompletion — API call to Create a completion for the chat message.
func (c *Client) CreateChatCompletion(
	ctx context.Context,
	request ChatCompletionRequest,
	args ...any,
) (response ChatCompletionResponse, err error) {
	if request.Stream {
		err = ErrChatCompletionStreamNotSupported
		return
	}
	model := ""
	if len(args) > 0 {
		m, ok := args[0].(string)
		if !ok {
			err = ErrChatCompletionInvalidModel
			return
		}
		model = m
	}

	var req *http.Request
	if c.config.AutoAuthToken {
		req, err = c.newRequestWithToken(ctx, http.MethodPost, c.fullURL(model), withBody(request))
	} else {
		req, err = c.newRequest(ctx, http.MethodPost, c.fullURL(model), withBody(request))
	}

	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}
