package xunfei

import (
	"context"
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
	args ...any,
) (stream *ChatCompletionStream, err error) {
	model := "v2"
	if len(args) > 0 {
		model = args[0].(string)
	}

	url := c.config.BaseURL + "/" + model + "/chat"

	domain := parameterChatDomainGeneral
	if model == "v1" || model == "v1.5" {
		domain = parameterChatDomainGeneral
		url = c.config.BaseURL + "/v1.1/chat"
	} else if model == "v2" {
		domain = parameterChatDomainGeneralV2
		url = c.config.BaseURL + "/v2.1/chat"
	} else if model == "v3" {
		domain = parameterChatDomainGeneralV3
		url = c.config.BaseURL + "/v3.1/chat"
	}

	fullReq := ChatCompletionStreamRequest{
		Header: ChatCompletionStreamRequestHeader{
			AppId: c.config.appId,
			Uid:   request.Uid,
		},

		Parameter: ChatCompletionStreamRequestParameter{
			Chat: ChatCompletionStreamRequestParameterChat{
				Domain:      domain,
				Temperature: request.Temperature,
				TopK:        request.TopK,
				MaxTokens:   request.MaxTokens,
				Auditing:    request.Auditing,
			}},
		Payload: ChatCompletionStreamRequestPayload{
			ChatCompletionStreamRequestPayloadMessage{
				Text: request.Messages,
			},
		},
	}

	conn, err := c.connect(url)
	if err != nil {
		return
	}

	err = conn.WriteJSON(fullReq)
	if err != nil {
		return
	}

	stream = &ChatCompletionStream{newStreamReader(conn)}
	return
}
