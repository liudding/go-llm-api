package tencent

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type AuthToken struct {
	token     string
	expiresAt int64
	expiresIn int64
}

// generateSignature 签名
func generateSignature(req ChatCompletionFullRequest, secretKey string) string {
	signUrl := buildURL(req)

	mac := hmac.New(sha1.New, []byte(secretKey))
	mac.Write([]byte(signUrl))
	sign := mac.Sum([]byte(nil))
	return base64.StdEncoding.EncodeToString(sign)
}

// buildSignURL 构建签名url
func buildURL(req ChatCompletionFullRequest) string {
	params := make([]string, 0)
	params = append(params, fmt.Sprintf("app_id=%d", req.AppId))
	params = append(params, "secret_id="+req.SecretId)
	params = append(params, "timestamp="+strconv.Itoa(int(req.Timestamp)))
	if req.QueryId != "" {
		params = append(params, "query_id="+req.QueryId)
	}
	if req.Temperature > 0 {
		params = append(params, "temperature="+strconv.FormatFloat(float64(req.Temperature), 'f', -1, 64))
	}
	if req.TopP > 0 {
		params = append(params, "top_p="+strconv.FormatFloat(float64(req.TopP), 'f', -1, 64))
	}

	params = append(params, "stream="+strconv.Itoa(req.Stream))
	params = append(params, "expired="+strconv.Itoa(int(req.Expired)))

	var messageStr string
	for _, msg := range req.Messages {
		messageStr += fmt.Sprintf(`{"role":"%s","content":"%s"},`, msg.Role, msg.Content)
	}
	messageStr = strings.TrimSuffix(messageStr, ",")
	params = append(params, "messages=["+messageStr+"]")

	sort.Sort(sort.StringSlice(params))
	return getFullPath() + strings.Join(params, "&")
}
