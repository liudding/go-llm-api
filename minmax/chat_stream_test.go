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

	client := NewClient("1699935883452478", "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJHcm91cE5hbWUiOiLkuIrmtbfnn6XlkKbnn6XlkKbnp5HmioAiLCJVc2VyTmFtZSI6IuS4iua1t-efpeWQpuefpeWQpuenkeaKgCIsIkFjY291bnQiOiIiLCJTdWJqZWN0SUQiOiIxNjk5OTM1ODgzMDkwNDMzIiwiUGhvbmUiOiIiLCJHcm91cElEIjoiMTY5OTkzNTg4MzQ1MjQ3OCIsIlBhZ2VOYW1lIjoiIiwiTWFpbCI6InJkQGltaWFvYmFuLmNvbSIsIkNyZWF0ZVRpbWUiOiIyMDIzLTEyLTE1IDE1OjEyOjAxIiwiaXNzIjoibWluaW1heCJ9.EI2MNK5W4b6Sqx_a0j9aqdSCSuxMoCnd6qoOoeoEHD01nb4tCrDHvHiF7kMkv9Wylb-sLbJTCXABTC8jViTBuBfuhrVkNhKXFa2qV0qLrva2dXkrguEkRcP5OcDhxljzPEts8LTbDtcCMkcWzNLH1asJNrnPUBQYYugVqiVkVZMBjZm7OxoDz_veE3vMFSOdak77tPCMAhnnbBN-6qYuaVogqyFOC6jHuty2tYitn-HPQ1Eo6lRihUm_Ev7QFF8QGHN7v1wYwftvZ_4xg5SJEYrhn-cBelVo31lc--JCb2G7Ots5p0n_OUf67g_oIkTLZKhNd-JhBJhT_SLprB4Y-A")
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
