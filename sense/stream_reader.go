package sense

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	utils "github.com/liudding/go-llm-api/internal"
	"github.com/liudding/go-llm-api/internal/sse"
	"io"
	"net/http"
)

type streamable interface {
	ChatCompletionResponse
}

type streamReader struct {
	isFinished bool

	reader         *sse.EventStreamReader
	response       *http.Response
	errAccumulator utils.ErrorAccumulator
	unmarshaler    utils.Unmarshaler
}

func newStreamReader(response *http.Response, emptyMessagesLimit uint) *streamReader {
	reader := sse.NewEventStreamReader(bufio.NewReader(response.Body), 1024, emptyMessagesLimit)

	return &streamReader{
		reader:         reader,
		response:       response,
		errAccumulator: utils.NewErrorAccumulator(),
		unmarshaler:    &utils.JSONUnmarshaler{},
	}
}

func (stream *streamReader) Recv() (response ChatCompletionResponse, err error) {
	if stream.isFinished {
		err = io.EOF
		return
	}

	event, err := stream.reader.Recv()
	if err != nil {
		return
	}

	if event.Data == nil {
		err = json.Unmarshal(event.Raw, &response)
		if err != nil {
			return
		}

		if response.Error.Code > 0 {
			err = fmt.Errorf("[%d] %s", response.Error.Code, response.Error.Message)
			return
		}

		err = errors.New("empty content")
		return
	}

	err = json.Unmarshal(event.Data, &response)
	if err != nil {
		return
	}

	if len(response.Data.Choices) == 0 {
		err = errors.New("empty content")
		return
	}

	if response.Data.Choices[0].FinishReason == "stop" {
		stream.isFinished = true
		err = io.EOF
		return
	}

	if response.Data.Choices[0].FinishReason == "length" {
		err = errors.New("too long content")
		return
	}

	if response.Data.Choices[0].FinishReason == "context" {
		err = errors.New("too long context")
		return
	}

	if response.Data.Choices[0].FinishReason == "sensitive" {
		err = errors.New("sensitive")
		return
	}

	return
}

func (stream *streamReader) Close() {
	err := stream.response.Body.Close()
	if err != nil {
		return
	}
}
