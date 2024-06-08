package auth

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Dao interface {
	GetUser(Username string) (*User, error)
	InsertUser(*User) error
	GetAllUsers() ([]User, error)
	GetUserByID(UserId uuid.UUID)(*User,error)
}

type dao struct {
	database *gorm.DB
}

func NewDatabase(db *gorm.DB) Dao {
	return &dao{
		database: db,
	}
}
func (db *dao) InsertUser(incomingUser *User) error {
	incomingUser.UserID = uuid.New()
	newUser := db.database.Create(incomingUser)

	if newUser.Error != nil {
		return newUser.Error
	}

	return nil
}
func (db *dao) GetUser(username string) (*User, error) {
	var user User
	getUser := db.database.First(&user, "username = ?", username)

	if getUser.Error != nil && getUser.Error.Error() == "record not found" {
		return nil, nil
	}

	if getUser.Error != nil {
		return nil, getUser.Error
	}

	return &user, nil
}

func (db *dao) GetUserByID(userID uuid.UUID) (*User, error) {
	var user User
	getUser := db.database.First(&user, "user_id = ?", userID)

	if getUser.Error != nil && getUser.Error.Error() == "record not found" {
		return nil, nil
	}

	if getUser.Error != nil {
		return nil, getUser.Error
	}

	return &user, nil
}

func (db *dao) GetAllUsers() ([]User, error) {
	var users []User

	result := db.database.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil

}
