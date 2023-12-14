package sense

import (
	"context"
	"net/http"
)

// ChatCompletionStream
type ChatCompletionStream struct {
	*streamReader
}

// CreateChatCompletionStream — API call to create a chat completion w/ streaming support.
func (c *Client) CreateChatCompletionStream(
	ctx context.Context,
	request ChatCompletionRequest,
	args ...any,
) (stream *ChatCompletionStream, err error) {
	request.Stream = true
	if request.Model == "" && len(args) > 0 {
		m, ok := args[0].(string)
		if ok {
			request.Model = m
		}
	}

	req, err := c.newRequest(ctx, http.MethodPost, c.config.BaseURL+chatCompletionsSuffix, withBody(request))
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
