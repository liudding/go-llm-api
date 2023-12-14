package xunfei_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/liudding/go-llm-api/internal/test/checks"
	. "github.com/liudding/go-llm-api/xunfei"
	"io"
	"testing"
)

func TestCreateChatCompletionStreamOnRealServer(t *testing.T) {
	client := NewClient("", "", "")
	stream, err := client.CreateChatCompletionStream(context.Background(), ChatCompletionRequest{
		Messages: []ChatCompletionMessage{
			{
				Role:    ChatMessageRoleUser,
				Content: "hi",
			},
		},
		Temperature: 0.7,
	})
	checks.NoError(t, err, "CreateCompletionStream returned error")
	defer stream.Close()

	fmt.Println("Stream response: ")
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Printf("\nStream finished:\n")
			return
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return
		}

		//fmt.Printf("error: \n")
		if len(response.Payload.Choices.Text) > 0 {
			fmt.Printf(response.Payload.Choices.Text[0].Content)
		} else {
			fmt.Printf("resp: %s\n", response.Payload.Choices.Content)
		}

	}
}
