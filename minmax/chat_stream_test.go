package minmax_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/liudding/go-llm-api/internal/test/checks"
	. "github.com/liudding/go-llm-api/minmax"
	"io"
	"testing"
)

func TestCreateChatCompletionStreamOnRealServer(t *testing.T) {

	client := NewClient("", "")
	stream, err := client.CreateChatCompletionStream(context.Background(), ChatCompletionRequest{
		Prompt: "你是一个擅长发现故事中蕴含道理的专家，你很善于基于我给定的故事发现其中蕴含的道理",
		Messages: []ChatCompletionMessage{{
			SenderType: "USER",
			Text:       "hi",
		}},
		RoleMeta: ChatCompletionRequestRoleMeta{
			UserName: "我",
			BotName:  "AI助手",
		},
		Stream:         true,
		UseStandardSse: true,
		Temperature:    0.7,
	}, ModelAbab5Chat)
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
		fmt.Printf("resp: %s\n", response.Choices[0].Delta)
	}
}
