package zhipu

import (
	"net/http"
)

const (
	zhipuAPIURLv3                  = "https://open.bigmodel.cn/api/paas/v3/model-api"
	defaultEmptyMessagesLimit uint = 300
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

		BaseURL: zhipuAPIURLv3,

		HTTPClient: &http.Client{},
	}
}

func (ClientConfig) String() string {
	return "<Zhipu AI API ClientConfig>"
}
