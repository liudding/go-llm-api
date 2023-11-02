# go-llm-api


百度千帆、智谱 Ai 等的大语言模型非官方 API，支持以下模型：

* 百度千帆平台所有模型
* 智谱 AI 模型

## Installation

```
go get github.com/liudding/go-llm-api
```
要求 go 版本最低是 1.18.


## Usage

### 百度千帆 Ernie:

```go
package main

import (
	"context"
	"fmt"
	 "github.com/liudding/github.com/liudding/go-llm-api/baidu"
)

func main() {
	client := NewClient("xxxx", "yyyy", true)
	stream, err := client.CreateChatCompletionStream(context.Background(), ChatCompletionRequest{
		Messages: []ChatCompletionMessage{
			{
				Role:    ChatMessageRoleUser,
				Content: "Hello!",
			},
		},
		Stream: true,
	})

	defer stream.Close()

	fmt.Println("Stream response: ")
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Printf("\nStream finished: %d %s\n", response.ErrorCode, response.ErrorMsg)
			return
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return
		}

		fmt.Printf("error: %s\n", response.ErrorMsg)
		fmt.Printf("resp: %s\n", response.Result)
	}
}

```

### 智谱

```go
token, _ := zhipu.GenerateAuthToken(l.svcCtx.Config.ZhipuAi.Key, time.Minute*60)
client := zhipu.NewClient(token)

stream, err := client.CreateChatCompletionStream(context.Background(), zhipu.ChatCompletionRequest{
	Prompt: []zhipu.ChatCompletionMessage{
		{
			Role:    zhipu.ChatMessageRoleUser,
			Content: prompt,
		},
	},
	Temperature: 0.7,
	//RequestId:   "",
	Incremental: true,
	//ReturnType:  "",
	//Ref:         ChatCompletionRef{},
}, zhipu.ModelTurbo)

defer stream.Close()
```
