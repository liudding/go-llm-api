package zhipu

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"strings"
	"time"
)

const defaultExpiresIn = 60 * time.Minute

type JwtClaims struct {
	APIKey    string `json:"api_key"`
	Exp       int64  `json:"exp"`
	Timestamp int64  `json:"timestamp"`
}

type AuthToken struct {
	token     string
	expiresAt int64
	expiresIn int64
}

// GenerateAuthToken
func GenerateAuthToken(apiKey string, expiresIn time.Duration) (string, error) {
	if apiKey == "" {
		return "", errors.New("api key 不能为空")
	}
	if !strings.Contains(apiKey, ".") {
		return "", errors.New("api key 格式不正确")
	}

	apiKeyInfo := strings.Split(apiKey, ".")
	key, secret := apiKeyInfo[0], apiKeyInfo[1]

	if expiresIn == 0 {
		expiresIn = defaultExpiresIn
	}

	return createToken(JwtClaims{
		key,
		time.Now().Add(expiresIn).Unix(),
		time.Now().Unix(),
	}, secret)
}

func createToken(claims JwtClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"api_key":   claims.APIKey,
		"exp":       claims.Exp,
		"timestamp": claims.Timestamp,
	})

	token.Header["alg"] = "HS256"
	token.Header["sign_type"] = "SIGN"
	res, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return res, nil
}
