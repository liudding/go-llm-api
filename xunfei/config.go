package xunfei

const (
	xunfeiAPIURL                   = "wss://spark-api.xf-yun.com"
	defaultEmptyMessagesLimit uint = 10
)

// ClientConfig is a configuration of a client.
type ClientConfig struct {
	appId     string
	apiKey    string
	apiSecret string
	BaseURL   string

	EmptyMessagesLimit uint
}

func DefaultConfig(appId string, apiKey string, apiSecret string) ClientConfig {
	return ClientConfig{
		appId:     appId,
		apiSecret: apiSecret,
		apiKey:    apiKey,
		BaseURL:   xunfeiAPIURL,
	}
}

func (ClientConfig) String() string {
	return ""
}
