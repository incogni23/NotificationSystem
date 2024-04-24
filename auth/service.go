package auth

import (
	"errors"
	"time"

	"github.com/create"
	"github.com/validate"
)

type User struct {
	Username string
	Password string
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
	if existingUser.Username != "" {
		return errors.New("User already exists")
	}
	if err != nil {
		return err
	}

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
	if user.Password != password {
		return "", errors.New("wrong credentials")
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

// signup:
//u1 User:
// u1-> db
// succes -> if user is suceesfully created in db as well
// db X - fail

//Login
