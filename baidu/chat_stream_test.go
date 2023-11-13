package baidu_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/liudding/go-llm-api/baidu"
	"github.com/liudding/go-llm-api/internal/test/checks"
	"io"
	"net/http"
	"testing"
)

func TestCreateChatCompletionRealServer(t *testing.T) {
	client := NewClient("xxxx", "yyyy", true)
	resp, err := client.CreateChatCompletion(context.Background(), ChatCompletionRequest{
		Messages: []ChatCompletionMessage{
			{
				Role:    ChatMessageRoleUser,
				Content: "Hello!",
			},
		},
		Stream: false,
	}, "eb-instant")
	checks.NoError(t, err, "CreateCompletionStream returned error")

	println(resp.ErrorMsg)
	println(resp.Result)
}

func TestCreateChatCompletionStreamOnRealServer(t *testing.T) {
	client := NewClient("xxxx", "yyyy", true)
	stream, err := client.CreateChatCompletionStream(context.Background(), ChatCompletionRequest{
		Messages: []ChatCompletionMessage{
			{
				Role:    ChatMessageRoleUser,
				Content: "Hello!",
			},
		},
		Temperature: 2,
		Stream:      true,
	})
	checks.NoError(t, err, "CreateCompletionStream returned error")
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

func TestCreateChatCompletionStream(t *testing.T) {
	client, server, teardown := setupBaiduAITestServer()
	defer teardown()
	server.RegisterHandler("/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")

		// Send test responses
		dataBytes := []byte{}
		dataBytes = append(dataBytes, []byte("event: message\n")...)
		//nolint:lll
		data := `{"id":"1","object":"chat.completion","created":1598069254,"result":"response1"}`
		dataBytes = append(dataBytes, []byte("data: "+data+"\n\n")...)

		dataBytes = append(dataBytes, []byte("event: message\n")...)
		//nolint:lll
		data = `{"id":"2","object":"chat.completion","created":1598069255,"model":"gpt-3.5-turbo","result":"response2"}`
		dataBytes = append(dataBytes, []byte("data: "+data+"\n\n")...)

		dataBytes = append(dataBytes, []byte("event: done\n")...)
		dataBytes = append(dataBytes, []byte("data: [DONE]\n\n")...)

		_, err := w.Write(dataBytes)
		checks.NoError(t, err, "Write error")
	})

	stream, err := client.CreateChatCompletionStream(context.Background(), ChatCompletionRequest{
		Messages: []ChatCompletionMessage{
			{
				Role:    ChatMessageRoleUser,
				Content: "Hello!",
			},
		},
		Stream: true,
	})
	checks.NoError(t, err, "CreateCompletionStream returned error")
	defer stream.Close()

	expectedResponses := []ChatCompletionResponse{
		{
			ID:      "1",
			Object:  "chat.completion",
			Created: 1598069254,
			Result:  "response1",
		},
		{
			ID:      "2",
			Object:  "chat.completion",
			Created: 1598069255,
			Result:  "response2",
		},
	}

	for ix, expectedResponse := range expectedResponses {
		b, _ := json.Marshal(expectedResponse)
		t.Logf("%d: %s", ix, string(b))

		receivedResponse, streamErr := stream.Recv()
		checks.NoError(t, streamErr, "stream.Recv() failed")
		if !compareChatResponses(expectedResponse, receivedResponse) {
			t.Errorf("Stream response %v is %v, expected %v", ix, receivedResponse, expectedResponse)
		}
	}

	_, streamErr := stream.Recv()
	if !errors.Is(streamErr, io.EOF) {
		t.Errorf("stream.Recv() did not return EOF in the end: %v", streamErr)
	}

	_, streamErr = stream.Recv()

	checks.ErrorIs(t, streamErr, io.EOF, "stream.Recv() did not return EOF when the stream is finished")
	if !errors.Is(streamErr, io.EOF) {
		t.Errorf("stream.Recv() did not return EOF when the stream is finished: %v", streamErr)
	}
}

func TestCreateChatCompletionStreamError(t *testing.T) {
	client, server, teardown := setupBaiduAITestServer()
	defer teardown()
	server.RegisterHandler("/v1/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")

		// Send test responses
		dataBytes := []byte{}
		dataStr := []string{
			`{"error_code": 1, "error_msg": "Unknown error"}`,
		}
		for _, str := range dataStr {
			dataBytes = append(dataBytes, []byte(str+"\n")...)
		}

		_, err := w.Write(dataBytes)
		checks.NoError(t, err, "Write error")
	})

	stream, err := client.CreateChatCompletionStream(context.Background(), ChatCompletionRequest{
		Messages: []ChatCompletionMessage{
			{
				Role:    ChatMessageRoleUser,
				Content: "Hello!",
			},
		},
		Stream: true,
	})
	checks.NoError(t, err, "CreateCompletionStream returned error")
	defer stream.Close()

	_, streamErr := stream.Recv()
	checks.HasError(t, streamErr, "stream.Recv() did not return error")

	var apiErr *APIError
	if !errors.As(streamErr, &apiErr) {
		t.Errorf("stream.Recv() did not return APIError")
	}
	t.Logf("%+v\n", apiErr)
}

// Helper funcs.
func compareChatResponses(r1, r2 ChatCompletionResponse) bool {
	if r1.ID != r2.ID || r1.Object != r2.Object || r1.Created != r2.Created || r1.Result != r2.Result {
		return false
	}

	return true
}
