package sense

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

const defaultExpiresIn = 60 * time.Minute

type JwtClaims struct {
	Iss string
	Exp int64
	Nbf int64
}

type AuthToken struct {
	token     string
	expiresAt int64
	expiresIn int64
}

// GenerateAuthToken
func GenerateAuthToken(accessKey string, secretKey string, expiresIn time.Duration) (string, error) {
	if accessKey == "" {
		return "", errors.New("api key 不能为空")
	}
	if secretKey == "" {
		return "", errors.New("secret key 不能为空")
	}

	if expiresIn == 0 {
		expiresIn = defaultExpiresIn
	}

	return createToken(JwtClaims{
		accessKey,
		time.Now().Add(expiresIn).Unix(),
		time.Now().Add(-2 * time.Second).Unix(),
	}, secretKey)
}

func createToken(claims JwtClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": claims.Iss,
		"exp": claims.Exp,
		"nbf": claims.Nbf,
	})
	res, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return res, nil
}
