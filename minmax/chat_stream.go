package minmax

import (
	"context"
	"net/http"
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
	model string,
) (stream *ChatCompletionStream, err error) {
	if request.Model == "" && model != "" {
		request.Model = model
	}
	req, err := c.newRequest(ctx, http.MethodPost, c.config.BaseURL, withQuery(map[string]string{
		"GroupId": c.config.groupId,
	}), withBody(request))
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
