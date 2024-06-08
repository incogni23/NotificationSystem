package create

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

func CreateToken(expireIn time.Duration, secretKey string, email string, userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{}

	claims["email"] = email
	claims["user_id"]= userID

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
