package auth_test

import (
	"errors"
	"testing"
	"time"

	"github.com/auth"
	mocks "github.com/auth/mocks"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
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
	password := "Ankita@2307"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mockUser := auth.User{
		Username: "Ankita",
		Password: string(hashedPassword),
	}

	t.Run("Test Login Success", func(t *testing.T) {
		dbmock.On("GetUser", "Ankita").Return(&mockUser, nil)

		authService := auth.NewDBVar(dbmock)
		token, err := authService.Login("Ankita", password)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Test Login Failed- Wrong credentials", func(t *testing.T) {
		dbmock.On("GetUser", "Ankita").Return(&mockUser, nil)

		authService := auth.NewDBVar(dbmock)
		token, err := authService.Login("Ankita", "wrongpassword")

		assert.Error(t, err)
		assert.Equal(t, "", token)
	})

	t.Run("Test Login Failed- User not found", func(t *testing.T) {
		dbmock.On("GetUser", "Ankita").Return(nil, nil)

		authService := auth.NewDBVar(dbmock)
		token, err := authService.Login("Ankita", password)

		assert.EqualError(t, err, "record doesnt exist")
		assert.Equal(t, "", token)
	})

	t.Run("Test Login Failed- Db failed", func(t *testing.T) {
		dbmock.On("GetUser", "Ankita").Return(nil, errors.New("database error"))

		authService := auth.NewDBVar(dbmock)
		token, err := authService.Login("Ankita", password)

		assert.EqualError(t, err, "database error")
		assert.Equal(t, "", token)
	})

	t.Run("Test Login Failed- Token Invalid", func(t *testing.T) {
		invalidToken := "invalidtoken"

		authService := auth.NewDBVar(dbmock)
		_, err := authService.LoginWithToken(invalidToken)

		assert.EqualError(t, err, "invalid token")
	})
	t.Run("Test Login Failed- Token expired", func(t *testing.T) {
		expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": "Ankita",
			"exp":      time.Now().Add(-time.Minute).Unix(),
		})
		tokenString, _ := expiredToken.SignedString([]byte("secretkey"))

		authService := auth.NewDBVar(dbmock)
		_, err := authService.LoginWithToken(tokenString)

		assert.EqualError(t, err, "invalid token")
	})
}
