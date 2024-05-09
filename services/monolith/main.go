package main

import (
	"github.com/auth"
	"github.com/database"
	"github.com/gin-gonic/gin"

	"github.com/monolith/routes"
)

func main() {
	//dep, err := dependencies.Init()
	//if err != nil {
	//	panic(err)
	//	}

	r := gin.New()

	db, err := database.SetupEnvAndDB()
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&auth.User{})

	authDao := auth.NewDatabase(db)

	authService := auth.NewDBVar(authDao)

	unprotectedGroup := routes.UnprotectedGroup(r)

	authEndpoint := auth.NewEndpoint(authService)

	unprotectedGroup.AuthGroup.POST("/signup", authEndpoint.Signup)

	unprotectedGroup.AuthGroup.POST("/login", authEndpoint.Login)

	r.Run(":9111")

}
