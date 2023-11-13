package tencent

import (
	"context"
	"errors"
	"net/http"
	"time"
)

const chatCompletionsSuffix = "/chat/completions"

var (
	ErrChatCompletionInvalidModel       = errors.New("this model is not supported with this method, please use CreateCompletion client method instead") //nolint:lll
	ErrChatCompletionStreamNotSupported = errors.New("streaming is not supported with this method, please use CreateChatCompletionStream")              //nolint:lll
)

// CreateChatCompletion — API call to Create a completion for the chat message.
func (c *Client) CreateChatCompletion(
	ctx context.Context,
	request ChatCompletionRequest,
) (response ChatCompletionResponse, err error) {

	fullReq := ChatCompletionFullRequest{
		AppId:                 c.config.appId,
		SecretId:              c.config.secretId,
		Timestamp:             time.Now().Unix(),
		Expired:               time.Now().Unix() + 24*60*60,
		ChatCompletionRequest: request,
	}
	fullReq.Stream = 0

	req, err := c.newRequest(ctx, http.MethodPost, getFullURL(), withBody(fullReq))
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}
