package validate

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

func Validate(tokenString, secretKey string) error {
	token, err := jwt.Parse(tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("Token is invalid")

	}
	return nil
}
