package auth

import (
	"gorm.io/gorm"
)

// dao -> dumb way to interac w db

// connect w db
// credentials

// query: update table
// existing user check - where usrn = " given",

//if fail : add a new record (user))- insert

//else : err

type Dao interface {
	GetUser(Username string) (*User, error)
	InsertUser(User) error
}

type dao struct {
	database gorm.DB
}

func NewDatabase(db gorm.DB) Dao {
	return &dao{
		database: db,
	}
}

func (db *dao) InsertUser(u User) error {
	newUser := db.database.Create(u)
	if newUser.Error != nil {
		return newUser.Error
	}
	return nil
}

func (db *dao) GetUser(username string) (*User, error) {
	var user User
	getUser := db.database.First(&user, "username = ?", username)
	if getUser.Error != nil {
		return nil, getUser.Error
	}
	return &user, nil

}

//struct- database.go-> gorm.db
//constructor sort of
//new fun1//
