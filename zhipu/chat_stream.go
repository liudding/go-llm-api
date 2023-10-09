package zhipu

import (
	"context"
	"fmt"
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
	url := fmt.Sprintf("%s/%s/sse-invoke", c.config.BaseURL, model)

	req, err := c.newRequest(ctx, http.MethodPost, url, withBody(request))
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
