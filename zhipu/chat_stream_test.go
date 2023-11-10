package zhipu_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/liudding/go-llm-api/internal/test/checks"
	. "github.com/liudding/go-llm-api/zhipu"
	"io"
	"testing"
	"time"
)

func TestCreateChatCompletionStreamOnRealServer(t *testing.T) {
	token, _ := GenerateAuthToken("xxx", time.Minute*60)

	client := NewClient(token)
	stream, err := client.CreateChatCompletionStream(context.Background(), ChatCompletionRequest{
		Prompt: []ChatCompletionMessage{
			{
				Role:    ChatMessageRoleUser,
				Content: "hi",
			},
		},
		//RequestId:   "",
		Incremental: true,
		//ReturnType:  "",
		//Ref:         ChatCompletionRef{},
	}, ModelChatGLMStd)
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
		fmt.Printf("resp: %s\n", response.Data)
	}
}
