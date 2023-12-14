package sense

import (
	"encoding/json"
	"fmt"
)

// APIError provides error information returned by the API.
type APIError struct {
	Err struct {
		Code    int    `json:"code,omitempty"`
		Message string `json:"message"`
		Details []any  `json:"details"`
	} `json:"error"`
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
	if e.Err.Code > 0 {
		return fmt.Sprintf("error, error code: %d, message: %s", e.Err.Code, e.Err.Message)
	}

	return e.Err.Message
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
