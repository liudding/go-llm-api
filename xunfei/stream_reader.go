package xunfei

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io"
)

type streamReader struct {
	isFinished bool
	conn       *websocket.Conn
}

func newStreamReader(conn *websocket.Conn) *streamReader {
	return &streamReader{
		conn: conn,
	}
}

func (stream *streamReader) Recv() (response ChatCompletionStreamResponse, err error) {
	if stream.isFinished {
		err = io.EOF
		return
	}

	err = stream.conn.ReadJSON(&response)
	if err != nil {
		return
	}

	if response.Header.Code != 0 {
		err = fmt.Errorf("[%d] %s", response.Header.Code, response.Header.Message)
		return
	}

	if response.Payload.Choices.Status == 2 {
		stream.isFinished = true
		return
	}

	return
}

func (stream *streamReader) Close() {
	err := stream.conn.Close()
	if err != nil {
		return
	}
}
