package baidu_test

import (
	"context"
	. "github.com/liudding/go-llm-api/baidu"
	"testing"
)

func TestChatCompletion(t *testing.T) {
	ctx := context.Background()
	client := NewClient("xxxx", "yyyy", true)
	resp, err := client.CreateChatCompletion(ctx, ChatCompletionRequest{
		Messages: []ChatCompletionMessage{
			{
				Role:    ChatMessageRoleUser,
				Content: "Hello!",
			},
		},
		Temperature: 0.7,
		Stream:      false,
		UserId:      "",
	})

	if err != nil {
		println(err.Error())
	}

	println(resp.ErrorMsg)
	if resp.ErrorCode != 0 {
		println(resp.ErrorMsg)
	}

	println(resp.Result)
}
