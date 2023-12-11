package tencent

import (
	"fmt"
	"net/http"
)

const (
	tencentAPIURLv1 = "https://hunyuan.cloud.tencent.com/hyllm/v1"
	protocol        = "https"
	host            = "hunyuan.cloud.tencent.com"
	path            = "/hyllm/v1/chat/completions?"
)

// ClientConfig is a configuration of a client.
type ClientConfig struct {
	appId      int64
	secretId   string
	secretKey  string
	BaseURL    string
	HTTPClient *http.Client

	EmptyMessagesLimit uint
}

func DefaultConfig(appId int64, secretId string, secretKey string) ClientConfig {
	return ClientConfig{
		appId:     appId,
		secretId:  secretId,
		secretKey: secretKey,

		BaseURL: tencentAPIURLv1,

		HTTPClient: &http.Client{},
	}
}

func getFullPath() string {
	return host + path
}

// getFullURL returns full URL for request.
func getFullURL() string {
	return fmt.Sprintf("%s://%s%s", protocol, host, path)
}
