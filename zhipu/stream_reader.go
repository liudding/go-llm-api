package zhipu

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	utils "llm-clients/internal"
	"net/http"
)

var (
	headerData  = []byte("data: ")
	errorPrefix = []byte(`{"error`)
)

var (
	ErrTooManyEmptyStreamMessages = errors.New("stream has sent too many empty messages")
)

type streamable interface {
	ChatCompletionResponse
}

type streamReader struct {
	emptyMessagesLimit uint
	isFinished         bool

	reader         *bufio.Reader
	response       *http.Response
	errAccumulator utils.ErrorAccumulator
	unmarshaler    utils.Unmarshaler
}

func (stream *streamReader) Recv() (response ChatCompletionResponse, err error) {
	if stream.isFinished {
		err = io.EOF
		return
	}

	response, err = stream.processLines()
	return
}

//nolint:gocognit
func (stream *streamReader) processLines() (ChatCompletionResponse, error) {
	var (
		emptyMessagesCount uint
		hasErrorPrefix     bool
	)

	for {
		rawLine, readErr := stream.reader.ReadBytes('\n')
		noSpaceLine := bytes.TrimSpace(rawLine)
		if stream.hasError(rawLine) {
			hasErrorPrefix = true
		}

		if readErr != nil || hasErrorPrefix {
			if hasErrorPrefix {
				var response ChatCompletionResponse
				unmarshalErr := stream.unmarshaler.Unmarshal(noSpaceLine, &response)
				if unmarshalErr != nil {
					return *new(ChatCompletionResponse), unmarshalErr
				}

				if hasErrorPrefix {
					noSpaceLine = bytes.TrimPrefix(noSpaceLine, headerData)
				}
				writeErr := stream.errAccumulator.Write(noSpaceLine)
				if writeErr != nil {
					return *new(ChatCompletionResponse), writeErr
				}

				return response, readErr
			}

			respErr := stream.unmarshalError()
			if respErr != nil {
				return *new(ChatCompletionResponse), fmt.Errorf("error, %w", respErr.Error)
			}
			return *new(ChatCompletionResponse), readErr
		}

		if !bytes.HasPrefix(noSpaceLine, headerData) {
			emptyMessagesCount++
			if stream.emptyMessagesLimit > 0 && emptyMessagesCount > stream.emptyMessagesLimit {
				return *new(ChatCompletionResponse), ErrTooManyEmptyStreamMessages
			}

			continue
		}

		noPrefixLine := bytes.TrimPrefix(noSpaceLine, headerData)
		if string(noPrefixLine) == "[DONE]" {
			stream.isFinished = true
			return *new(ChatCompletionResponse), io.EOF
		}

		var response ChatCompletionResponse
		unmarshalErr := stream.unmarshaler.Unmarshal(noPrefixLine, &response)
		if unmarshalErr != nil {
			return *new(ChatCompletionResponse), unmarshalErr
		}

		if response.IsEnd {
			stream.isFinished = true
		}

		return response, nil
	}
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
	stream.response.Body.Close()
}
