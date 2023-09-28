package baidu

import (
	"net/http"
)

const (
	baiduaiAPIURLv1                = "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop"
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

func DefaultConfig(clientId string, clientSecret string, auto bool) ClientConfig {
	return ClientConfig{
		authToken:    "",
		ClientId:     clientId,
		ClientSecret: clientSecret,
		GrantType:    "client_credentials",

		BaseURL: baiduaiAPIURLv1,

		HTTPClient: &http.Client{},

		EmptyMessagesLimit: defaultEmptyMessagesLimit,

		AutoAuthToken: auto,
	}
}

func DefaultConfigWithAuth(authToken string) ClientConfig {
	return ClientConfig{
		authToken: authToken,

		GrantType: "client_credentials",

		BaseURL: baiduaiAPIURLv1,

		HTTPClient: &http.Client{},

		EmptyMessagesLimit: defaultEmptyMessagesLimit,
	}
}

func (ClientConfig) String() string {
	return "<Baidu AI API ClientConfig>"
}
