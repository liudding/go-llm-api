package baidu

import (
	"context"
	"net/http"
)

// ChatCompletionStream
// Note: Perhaps it is more elegant to abstract Stream using generics.
type ChatCompletionStream struct {
	*streamReader
}

// CreateChatCompletionStream — API call to create a chat completion streaming
func (c *Client) CreateChatCompletionStream(
	ctx context.Context,
	request ChatCompletionRequest,
	args ...any,
) (stream *ChatCompletionStream, err error) {

	model := ""
	if len(args) > 0 {
		m, ok := args[0].(string)
		if !ok {
			err = ErrChatCompletionInvalidModel
			return
		}
		model = m
	}

	request.Stream = true

	var req *http.Request
	if c.config.AutoAuthToken {
		req, err = c.newRequestWithToken(ctx, http.MethodPost, c.fullURL(model), withBody(request))
	} else {
		req, err = c.newRequest(ctx, http.MethodPost, c.fullURL(model), withBody(request))
	}

	if err != nil {
		return nil, err
	}

	resp, err := sendRequestStream(c, req)
	if err != nil {
		return
	}
	stream = &ChatCompletionStream{
		streamReader: resp,
	}
	return
}
