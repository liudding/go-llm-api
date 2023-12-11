package sse

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	utils "github.com/liudding/go-llm-api/internal"
	"io"
	"regexp"
	"strings"
	"time"
)

var (
	headerId    = []byte("id:")
	headerEvent = []byte("event:")
	headerData  = []byte("data:")
	headerRetry = []byte("retry:")
)

var (
	ErrTooManyEmptyStreamMessages = errors.New("stream has sent too many empty messages")

	matchExtraKeyRegex, _ = regexp.Compile(`^([a-zA-z_-]+):`)
)

type Event struct {
	timestamp time.Time
	Id        []byte
	Data      []byte
	Event     []byte
	Retry     []byte
	Comment   []byte

	Extra map[string][]byte
	Other []byte

	Raw []byte
}

type EventStreamReader struct {
	emptyMessagesLimit uint
	isFinished         bool

	scanner        *bufio.Scanner
	encodingBase64 bool
	errAccumulator utils.ErrorAccumulator
	unmarshaler    utils.Unmarshaler
}

func NewEventStreamReader(stream *bufio.Reader, maxBufferSize int, emptyMessagesLimit uint) *EventStreamReader {
	scanner := bufio.NewScanner(stream)
	initBufferSize := minPosInt(4096, maxBufferSize)
	scanner.Buffer(make([]byte, initBufferSize), maxBufferSize)

	split := func(data []byte, atEOF bool) (int, []byte, error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		// We have a full event payload to parse.
		if i, nlen := containsDoubleNewline(data); i >= 0 {
			return i + nlen, data[0:i], nil
		}
		// If we're at EOF, we have all of the data.
		if atEOF {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	}
	// Set the split function for the scanning operation.
	scanner.Split(split)

	return &EventStreamReader{
		scanner:            scanner,
		emptyMessagesLimit: emptyMessagesLimit,
		errAccumulator:     utils.NewErrorAccumulator(),
		unmarshaler:        &utils.JSONUnmarshaler{},
	}
}

func (stream *EventStreamReader) Recv() (event *Event, err error) {
	if stream.isFinished {
		err = io.EOF
		return
	}

	event, err = stream.processEvent()
	return
}

func (stream *EventStreamReader) processEvent() (*Event, error) {
	var (
		emptyMessagesCount uint
	)

	msg, err := stream.ReadEvent()
	noSpaceMsg := bytes.TrimSpace(msg)

	if err != nil {
		if err == io.EOF {
			stream.isFinished = true
			return nil, io.EOF
		}

		return nil, err
	}

	if len(noSpaceMsg) < 1 {
		emptyMessagesCount++
		if emptyMessagesCount > stream.emptyMessagesLimit {
			return nil, ErrTooManyEmptyStreamMessages
		}
		return nil, nil
	}

	var event Event
	event.Raw = msg

	// Normalize the crlf to lf to make it easier to split the lines.
	// Split the line by "\n" or "\r", per the spec.
	for _, line := range bytes.FieldsFunc(msg, func(r rune) bool { return r == '\n' || r == '\r' }) {
		switch {
		case bytes.HasPrefix(line, headerId):
			event.Id = append([]byte(nil), trimHeader(len(headerId), line)...)
		case bytes.HasPrefix(line, headerData):
			// The spec allows for multiple data fields per event, concatenated them with "\n".
			event.Data = append(event.Data[:], append(trimHeader(len(headerData), line), byte('\n'))...)
		// The spec says that a line that simply contains the string "data" should be treated as a data field with an empty body.
		case bytes.Equal(line, bytes.TrimSuffix(headerData, []byte(":"))):
			event.Data = append(event.Data, byte('\n'))
		case bytes.HasPrefix(line, headerEvent):
			event.Event = append([]byte(nil), trimHeader(len(headerEvent), line)...)
		case bytes.HasPrefix(line, headerRetry):
			event.Retry = append([]byte(nil), trimHeader(len(headerRetry), line)...)
		default:
			matches := matchExtraKeyRegex.FindSubmatch(line)
			if len(matches) == 0 {
				event.Other = line
			} else {
				k := string(matches[1])
				v := strings.Replace(string(line), string(matches[0]), "", 1)

				if event.Extra == nil {
					event.Extra = make(map[string][]byte)
				}
				event.Extra[k] = []byte(v)
			}
		}
	}

	// Trim the last "\n" per the spec.
	event.Data = bytes.TrimSuffix(event.Data, []byte("\n"))

	if stream.encodingBase64 {
		buf := make([]byte, base64.StdEncoding.DecodedLen(len(event.Data)))

		n, err := base64.StdEncoding.Decode(buf, event.Data)
		if err != nil {
			err = fmt.Errorf("failed to decode event message: %s", err)
		}
		event.Data = buf[:n]
	}
	return &event, err
}

//func (stream *EventStreamReader) unmarshalError() (errResp *zhipu.ErrorResponse) {
//	errBytes := stream.errAccumulator.Bytes()
//	if len(errBytes) == 0 {
//		return
//	}
//
//	err := stream.unmarshaler.Unmarshal(errBytes, &errResp)
//	if err != nil {
//		errResp = nil
//	}
//
//	return
//}

// ReadEvent scans the EventStream for events.
func (stream *EventStreamReader) ReadEvent() ([]byte, error) {
	if stream.scanner.Scan() {
		event := stream.scanner.Bytes()
		return event, nil
	}
	if err := stream.scanner.Err(); err != nil {
		if err == context.Canceled {
			return nil, io.EOF
		}
		return nil, err
	}
	return nil, io.EOF
}

// Returns the minimum non-negative value out of the two values. If both
// are negative, a negative value is returned.
func minPosInt(a, b int) int {
	if a < 0 {
		return b
	}
	if b < 0 {
		return a
	}
	if a > b {
		return b
	}
	return a
}

// Returns a tuple containing the index of a double newline, and the number of bytes
// represented by that sequence. If no double newline is present, the first value
// will be negative.
func containsDoubleNewline(data []byte) (int, int) {
	// Search for each potentially valid sequence of newline characters
	crcr := bytes.Index(data, []byte("\r\r"))
	lflf := bytes.Index(data, []byte("\n\n"))
	crlflf := bytes.Index(data, []byte("\r\n\n"))
	lfcrlf := bytes.Index(data, []byte("\n\r\n"))
	crlfcrlf := bytes.Index(data, []byte("\r\n\r\n"))
	// Find the earliest position of a double newline combination
	minPos := minPosInt(crcr, minPosInt(lflf, minPosInt(crlflf, minPosInt(lfcrlf, crlfcrlf))))
	// Detemine the length of the sequence
	nlen := 2
	if minPos == crlfcrlf {
		nlen = 4
	} else if minPos == crlflf || minPos == lfcrlf {
		nlen = 3
	}
	return minPos, nlen
}

func trimHeader(size int, data []byte) []byte {
	if data == nil || len(data) < size {
		return data
	}

	data = data[size:]
	// Remove optional leading whitespace
	if len(data) > 0 && data[0] == 32 {
		data = data[1:]
	}
	// Remove trailing new line
	if len(data) > 0 && data[len(data)-1] == 10 {
		data = data[:len(data)-1]
	}
	return data
}
