package zhipu

import (
	"net/http"
)

const (
	zhipuAPIURLv1                  = "https://open.bigmodel.cn/api/paas/v3/model-api"
	defaultEmptyMessagesLimit uint = 300
)

// ClientConfig is a configuration of a client.
type ClientConfig struct {
	authToken string

	GrantType    string
	ClientId     string
	ClientSecret string

	BaseURL    string
	HTTPClient *http.Client

	EmptyMessagesLimit uint

	AutoAuthToken bool   // 是否自动刷新 auth token。如果 true，那最好使用一个全局的 client
	AuthAPI       string // 授权 api
}

func DefaultConfig(authToken string) ClientConfig {
	return ClientConfig{
		authToken: authToken,

		BaseURL: zhipuAPIURLv1,

		HTTPClient: &http.Client{},

		EmptyMessagesLimit: defaultEmptyMessagesLimit,
	}
}

func (ClientConfig) String() string {
	return "<Baidu AI API ClientConfig>"
}
