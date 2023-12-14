package sense

import (
	"net/http"
)

const (
	apiURLv1                       = "https://api.sensenova.cn/v1/llm"
	defaultEmptyMessagesLimit uint = 10
)

// ClientConfig is a configuration of a client.
type ClientConfig struct {
	authToken  string
	BaseURL    string
	HTTPClient *http.Client

	EmptyMessagesLimit uint
}

func DefaultConfig(authToken string) ClientConfig {
	return ClientConfig{
		authToken: authToken,

		BaseURL: apiURLv1,

		HTTPClient:         &http.Client{},
		EmptyMessagesLimit: defaultEmptyMessagesLimit,
	}
}

func (ClientConfig) String() string {
	return ""
}
