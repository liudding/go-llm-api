package tencent

import (
	"encoding/json"
	"fmt"
)

// APIError provides error information returned API.
type APIError struct {
	ErrorCode int    `json:"error_code,omitempty"`
	ErrorMsg  string `json:"error_msg"`

	HTTPStatusCode int `json:"-"`
}

// RequestError provides informations about generic request errors.
type RequestError struct {
	HTTPStatusCode int
	Err            error
}

type ErrorResponse struct {
	Error *APIError `json:"error,omitempty"`
}

func (e *APIError) Error() string {
	if e.ErrorCode > 0 {
		return fmt.Sprintf("error, error code: %d, message: %s", e.ErrorCode, e.ErrorMsg)
	}

	return e.ErrorMsg
}

func (e *APIError) UnmarshalJSON(data []byte) (err error) {
	err = json.Unmarshal(data, &e)
	if err != nil {
		return
	}
	return
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("error, status code: %d, message: %s", e.HTTPStatusCode, e.Err)
}

func (e *RequestError) Unwrap() error {
	return e.Err
}
