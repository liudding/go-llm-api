package minmax

import (
	"net/http"
)

const (
	minmaxAPIURLv1                 = "https://api.minimax.chat/v1/text/chatcompletion"
	defaultEmptyMessagesLimit uint = 10
)

// ClientConfig is a configuration of a client.
type ClientConfig struct {
	authToken  string
	groupId    string
	BaseURL    string
	HTTPClient *http.Client

	EmptyMessagesLimit uint
}

func DefaultConfig(groupId string, apiKey string) ClientConfig {
	return ClientConfig{
		authToken: apiKey,
		groupId:   groupId,

		BaseURL: minmaxAPIURLv1,

		HTTPClient:         &http.Client{},
		EmptyMessagesLimit: defaultEmptyMessagesLimit,
	}
}

func (ClientConfig) String() string {
	return ""
}
