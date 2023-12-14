package sense_test

import (
	"context"
	. "github.com/liudding/go-llm-api/sense"
	"testing"
	"time"
)

func TestChatCompletion(t *testing.T) {
	ctx := context.Background()
	token, _ := GenerateAuthToken("", "xxx", time.Minute*60)
	client := NewClient(token)
	resp, err := client.CreateChatCompletion(ctx, ChatCompletionRequest{
		Messages: []ChatCompletionMessage{
			{
				Role:    ChatMessageRoleUser,
				Content: "Hello!",
			},
		},
		Temperature: 0.7,
		Stream:      false,
		User:        "",
	})

	if err != nil {
		println(err.Error())
	}

	println(resp.Data.Choices[0].Message)
}
