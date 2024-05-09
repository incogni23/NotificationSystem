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

// new signup

func (dbv *DBVar) SignUp(u User) error {
	existingUser, err := dbv.db.GetUser(u.Username)
	if existingUser != nil {
		if existingUser.Username != "" {
			return errors.New("User already exists")
		}
	}
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	err = dbv.db.InsertUser(u)
	if err != nil {
		return errors.New("User creation failed")

	}
	return nil

}
func (dbv *DBVar) Login(username, password string) (string, error) {
	user, err := dbv.db.GetUser(username)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", errors.New("record doesnt exist")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
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
