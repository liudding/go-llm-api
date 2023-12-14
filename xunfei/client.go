package xunfei

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	config ClientConfig

	WebSocketDialer *websocket.Dialer
}

// NewClient creates new Xunfei Xinghuo AI API client.
func NewClient(appId string, apiKey string, apiSecret string) *Client {
	config := DefaultConfig(appId, apiKey, apiSecret)
	return NewClientWithConfig(config)
}

// NewClientWithConfig creates new API client for specified config.
func NewClientWithConfig(config ClientConfig) *Client {
	return &Client{
		config: config,
		WebSocketDialer: &websocket.Dialer{
			HandshakeTimeout: 5 * time.Second,
		},
	}
}

func (c *Client) connect(url string) (*websocket.Conn, error) {
	conn, resp, err := c.WebSocketDialer.Dial(assembleAuthUrl(url, c.config.apiKey, c.config.apiSecret), nil)
	if err != nil {
		return nil, errors.New(readResp(resp) + err.Error())
	} else if resp.StatusCode != 101 {
		return nil, errors.New("error to connect websocket: " + readResp(resp))
	}

	return conn, nil
}

func readResp(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("code=%d,body=%s", resp.StatusCode, string(b))
}

// 创建鉴权url  apikey 即 hmac username
func assembleAuthUrl(hosturl string, apiKey, apiSecret string) string {
	ul, err := url.Parse(hosturl)
	if err != nil {
		fmt.Println(err)
	}

	date := time.Now().UTC().Format(time.RFC1123)
	signParts := []string{"host: " + ul.Host, "date: " + date, "GET " + ul.Path + " HTTP/1.1"}
	signStr := strings.Join(signParts, "\n")

	sign := HmacWithShaTobase64("hmac-sha256", signStr, apiSecret)

	//构建请求参数 此时不需要urlencoding
	auth := fmt.Sprintf(`api_key="%s", algorithm="%s", headers="%s", signature="%s"`,
		apiKey, "hmac-sha256", "host date request-line", sign)
	authorization := base64.StdEncoding.EncodeToString([]byte(auth))

	v := url.Values{}
	v.Add("host", ul.Host)
	v.Add("date", date)
	v.Add("authorization", authorization)
	return hosturl + "?" + v.Encode()
}

func HmacWithShaTobase64(algorithm, data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}
