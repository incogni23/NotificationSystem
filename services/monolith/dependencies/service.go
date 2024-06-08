package dependencies

import (
	"github.com/auth"
	"github.com/database"
	"gorm.io/gorm"
)

type Dependencies struct {
	db *gorm.DB
	AuthService auth.AuthServicer
}

func Init() (*Dependencies, error) {
	db, err := database.SetupEnvAndDB()
	if err != nil {
		return nil, err
	}

	authDao := auth.NewDatabase(db)

	authService := auth.NewService(authDao)

	return &Dependencies{
		db:          db,
		AuthService: authService,
	}, nil
}
