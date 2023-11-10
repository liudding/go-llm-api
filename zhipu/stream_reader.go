package zhipu

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	utils "github.com/liudding/go-llm-api/internal"
	"github.com/liudding/go-llm-api/internal/sse"
	"io"
	"net/http"
)

var (
	errorPrefix = []byte(`{"error`)
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

func newStreamReader(response *http.Response) *streamReader {
	reader := sse.NewEventStreamReader(bufio.NewReader(response.Body), 1024)

	return &streamReader{
		reader:         reader,
		response:       response,
		errAccumulator: utils.NewErrorAccumulator(),
		unmarshaler:    &utils.JSONUnmarshaler{},
	}
}

func (stream *streamReader) Recv() (response ChatCompletionResponse, err error) {
	event, err := stream.reader.Recv()
	if err != nil {
		return
	}

	if stream.isFinished {
		err = io.EOF
		return
	}

	e := string(event.Event)

	response.Id = string(event.Id)
	response.Data = string(event.Data)
	response.Event = e

	if e == "finish" {
		err = io.EOF

		if m, ok := event.Extra["meta"]; ok {
			var meta ChatCompletionResponseMeta
			err := json.Unmarshal(m, &meta)
			if err != nil {
				return ChatCompletionResponse{}, err
			}

			response.Meta = &meta
		}

		return
	}

	if e == "error" || e == "interrupted" {
		err = errors.New(string(event.Data))
		return
	}

	return
}

func (stream *streamReader) hasError(rawLine []byte) bool {
	noSpaceLine := bytes.TrimSpace(rawLine)
	return bytes.HasPrefix(noSpaceLine, errorPrefix)
}

func (stream *streamReader) unmarshalError() (errResp *ErrorResponse) {
	errBytes := stream.errAccumulator.Bytes()
	if len(errBytes) == 0 {
		return
	}

	err := stream.unmarshaler.Unmarshal(errBytes, &errResp)
	if err != nil {
		errResp = nil
	}

	return
}

func (stream *streamReader) Close() {
	err := stream.response.Body.Close()
	if err != nil {
		return
	}
}
