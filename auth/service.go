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
	UserID     uuid.UUID `json:"userID" gorm:"type:uuid;primaryKey"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	DbHost     string    `json:"-"`
	DbUser     string    `json:"-"`
	DbPassword string    `json:"-"`
}

func NewUserWithDefaults() *User {
	return &User{
		DbHost:     "localhost",
		DbUser:     "pikapika",
		DbPassword: "Ankita@2307",
	}
}

type AuthServicer interface {
	SignUp(u *User) (*User, error)
	Login(username, password string) (string, error)
	LoginWithToken(tokenString string) (string, error)
	GetAllUsers() ([]User, error)
	GetUserByID(userID uuid.UUID) (*User, error)
}

type Service struct {
	db Dao
}

func NewService(d Dao) AuthServicer {
	return &Service{
		db: d,
	}
}

func (dbv *Service) SignUp(incomingUser *User) (*User, error) {
	if incomingUser.Username == "" || incomingUser.Password == "" || incomingUser.Email == "" {
		return nil, errors.New("username, password, and email are required")
	}

	existingUser, err := dbv.db.GetUser(incomingUser.Username)

	if err != nil {
		return nil, err
	}

	if existingUser != nil && existingUser.Username != "" {
		return nil, errors.New("User already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(incomingUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	incomingUser.Password = string(hashedPassword)

	err = dbv.db.InsertUser(incomingUser)
	if err != nil {
		return nil, errors.New("User creation failed")
	}

	return incomingUser, nil
}

func (dbv *Service) Login(username, password string) (string, error) {
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

	tokenString, err := create.CreateToken(time.Hour*1, "secretKey", existingUser.Email, existingUser.UserID)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (dbv *Service) LoginWithToken(tokenString string) (string, error) {
	err := validate.Validate(tokenString, "secretKey")
	if err != nil {
		return "", errors.New("invalid token")
	}

	return "valid token", nil
}

func (dbv *Service) GetAllUsers() ([]User, error) {

	alluser, err := dbv.db.GetAllUsers()
	if err != nil {
		return nil, errors.New("cannot fetch all Users")
	}
	return alluser, nil
}

func (dbc *Service) GetUserByID(userID uuid.UUID) (*User, error) {
	return dbc.db.GetUserByID(userID)
}
