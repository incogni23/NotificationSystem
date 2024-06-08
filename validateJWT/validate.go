package validate

import (
	"fmt"
	"time"

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
		return fmt.Errorf("token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("failed to parse claims")
	}

	expirationTimeClaim, ok := claims["exp"]
	if !ok {
		return fmt.Errorf("exp Time not found in claims")
	}

	expirationTime, ok := expirationTimeClaim.(float64)
	if !ok {
		return fmt.Errorf("failed to parse exp time")
	}

	if time.Now().Unix() > int64(expirationTime) {
		return fmt.Errorf("token has expired")
	}

	return nil
}
