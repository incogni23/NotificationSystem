package auth

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

func GetUserIDFromToken(tokenString, secretKey string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("incorrect signing method %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return uuid.Nil, err
	}
	if !token.Valid {
		return uuid.Nil, fmt.Errorf("token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("failed to parse claims")
	}

	userIDstr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("user_id not found in claims")
	}
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user id")
	}

	return userID, nil
}
