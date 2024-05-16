package auth

import (
	"errors"
	"time"

	"github.com/create"
	"github.com/google/uuid"
	"github.com/validate"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID   uuid.UUID `json:"userID" gorm:"type:uuid;primaryKey"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}

type AuthServicer interface {
	SignUp(u User) error
	Login(username, password string) (string, error)
	LoginWithToken(tokenString string) (string, error)
}

type DBVar struct {
	db Dao
}

func NewDBVar(d Dao) AuthServicer {
	return &DBVar{
		db: d,
	}
}

func (dbv *DBVar) SignUp(incomingUser User) error {
	existingUser, err := dbv.db.GetUser(incomingUser.Username)
	if existingUser != nil && existingUser.Username != "" {
		return errors.New("User already exists")
	}

	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(incomingUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	incomingUser.Password = string(hashedPassword)

	err = dbv.db.InsertUser(incomingUser)
	if err != nil {
		return errors.New("User creation failed")

	}

	return nil

}
func (dbv *DBVar) Login(username, password string) (string, error) {
	existingUser, err := dbv.db.GetUser(username)
	if err != nil {
		return "", err
	}

	if existingUser == nil {
		return "", errors.New("record doesnt exist")
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password))
	if err != nil {
		return "", err
	}

	tokenString, err := create.CreateToken(time.Minute*1, "secretkey")
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (dbv *DBVar) LoginWithToken(tokenString string) (string, error) {
	err := validate.Validate(tokenString, "secretkey")
	if err != nil {
		return "", errors.New("invalid token")
	}

	return "valid token", nil
}
