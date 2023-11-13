package tencent

import (
	"context"
	"net/http"
	"time"
)

// ChatCompletionStream
// Note: Perhaps it is more elegant to abstract Stream using generics.
type ChatCompletionStream struct {
	*streamReader
}

// CreateChatCompletionStream — API call to create a chat completion w/ streaming support.
func (c *Client) CreateChatCompletionStream(
	ctx context.Context,
	request ChatCompletionRequest,
) (stream *ChatCompletionStream, err error) {
	url := getFullURL()

	fullReq := ChatCompletionFullRequest{
		AppId:                 c.config.appId,
		SecretId:              c.config.secretId,
		Timestamp:             time.Now().Unix(),
		Expired:               time.Now().Unix() + 24*60*60,
		ChatCompletionRequest: request,
	}
	fullReq.Stream = 1

	req, err := c.newRequest(ctx, http.MethodPost, url, withBody(fullReq))
	if err != nil {
		return nil, err
	}

	resp, err := sendRequestStream[ChatCompletionResponse](c, req)
	if err != nil {
		return
	}
	stream = &ChatCompletionStream{
		streamReader: resp,
	}
	return
}
