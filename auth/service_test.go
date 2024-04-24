package auth_test

import (
	"errors"
	"testing"
	"time"

	"github.com/auth"
	mocks "github.com/auth/mocks"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestSignup(t *testing.T) {

	dbmock := mocks.NewDao(t)
	mockUser := auth.User{
		Username: "Ankita",
		Password: "Ankita@2307",
	}
	t.Run("Test Signup Sucess", func(t *testing.T) {

		blankuser := auth.User{}

		dbmock.On("InsertUser", mockUser).Return(nil)
		dbmock.On("GetUser", "Ankita").Return((&blankuser), nil)

		sp := auth.NewDBVar(dbmock)
		signupsuceess := sp.SignUp(mockUser)

		assert.NoError(t, signupsuceess)

	})

	t.Run("Test Signup Failure - db failed", func(t *testing.T) {
		blankuser := auth.User{}
		dbmock.On("GetUser", "Ankita").Return((&blankuser), errors.New("some error"))

		sp := auth.NewDBVar(dbmock)
		signupfailed := sp.SignUp(mockUser)

		assert.EqualError(t, signupfailed, "some error")

	})

	t.Run("Test Signup Failure - Duplicate Records", func(t *testing.T) {
		//	dbmock.On("InsertUser", mockUser).Return(errors.New("User creation failed"))
		dbmock.On("GetUser", "Ankita").Return((&mockUser), nil)

		sp := auth.NewDBVar(dbmock)
		signupfailed := sp.SignUp(mockUser)

		assert.EqualError(t, signupfailed, "User already exists")

	})
}

func TestLogin(t *testing.T) {
	dbmock := mocks.NewDao(t)
	mockUser := auth.User{
		Username: "Ankita",
		Password: "Ankita@2307",
	}

	t.Run("Test Login Success", func(t *testing.T) {
		blankUser := auth.User{}

		dbmock.On("InsertUser", mockUser).Return(nil)
		dbmock.On("GetUser", "Ankita").Return(&blankUser, nil)

		authService := auth.NewDBVar(dbmock)
		err := authService.SignUp(mockUser)

		assert.NoError(t, err)

	})

	t.Run("Test Login Failed- Wrong credentials", func(t *testing.T) {
		dbmock.On("GetUser", "Ankita").Return(&mockUser, nil)

		authService := auth.NewDBVar(dbmock)
		err := authService.SignUp(mockUser)

		assert.EqualError(t, err, "User already exists")

	})

	t.Run("Test Login Failed- User not found", func(t *testing.T) {

	})
	t.Run("Test Login Failed- Db failed", func(t *testing.T) {
		dbmock.On("GetUser", "Ankita").Return(&auth.User{}, errors.New("database error"))

		authService := auth.NewDBVar(dbmock)
		err := authService.SignUp(mockUser)

		assert.EqualError(t, err, "database error")
	})

	t.Run("Test Login Failed- Token Invalid", func(t *testing.T) {
		invalidtoken := "invalidtoken"

		authService := auth.NewDBVar(dbmock)
		_, err := authService.LoginWithToken(invalidtoken)

		assert.EqualError(t, err, "invalid token")

	})
	t.Run("Test Login Failed- Token expired", func(t *testing.T) {
		expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": "Ankita",
			"exp":      time.Now().Add(-time.Minute).Unix(),
		})
		tokenString, _ := expiredToken.SignedString([]byte("secretkey"))

		authservice := auth.NewDBVar(dbmock)
		_, err := authservice.LoginWithToken(tokenString)

		assert.EqualError(t, err, "invalid token")

	})

}
