package main

import (
	"github.com/auth"
	"github.com/database"
	"github.com/gin-gonic/gin"

	"github.com/monolith/configurations"
	"github.com/monolith/dependencies"
	"github.com/monolith/order"
	"github.com/monolith/payments"
	"github.com/monolith/routes"
)

func main() {

	r := gin.New()

	db, err := database.SetupEnvAndDB()
	if err != nil {
		panic(err)
	}
	
	db.AutoMigrate(&auth.User{},
		&payments.PaymentMethod{},
		&payments.PaymentGateway{},
		&order.Order{},
		&payments.Payment{},
		&payments.ThirdPartyToken{},
		&payments.PaymentConfiguration{},
	)

	dependencies.SeedPayment(db)

	authDao := auth.NewDatabase(db)

	authService := auth.NewDBVar(authDao)

	unprotectedGroup := routes.UnprotectedGroup(r)

	authEndpoint := auth.NewEndpoint(authService)

	unprotectedGroup.AuthGroup.POST("/signup", authEndpoint.Signup)

	unprotectedGroup.AuthGroup.POST("/login", authEndpoint.Login)

	configService := configurations.NewDBConfig(db)

	r.GET("/payment_method/:method_id", configService.PaymentMethod)
	r.GET("/payment_gateway/:gateway_id", configService.PaymentGateway)
	r.POST("/payment_config", configService.PaymentConfiguration)

	r.Run(":9111")

}
