package create

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(expireIn time.Duration, secretKey string) (string, error) {
	claims := jwt.MapClaims{}
	if expireIn > 0 {
		claims["exp"] = time.Now().Add(expireIn).Unix()

	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil

}
