package main

import (
	"github.com/auth"
	"github.com/database"
	"github.com/gin-gonic/gin"

	"github.com/monolith/configurations"
	"github.com/monolith/dependencies"
	"github.com/monolith/order"
	"github.com/monolith/routes"
)

func main() {
	r := gin.New()

	db, err := database.SetupEnvAndDB()
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&auth.User{},
	)

	authDao := auth.NewDatabase(db)
	authService := auth.NewService(authDao)

	configService := configurations.NewDBConfig(authService)

	orderDao := order.NewOrderDao()
	orderService:= order.NewOrderService(orderDao)
	orderEndpoint := order.NewOrderEndpoint(authService,orderService)
	
	unprotectedGroup := routes.UnprotectedGroup(r)
	authEndpoint := auth.NewEndpoint(authService)

	unprotectedGroup.AuthGroup.POST("/signup", authEndpoint.Signup)
	unprotectedGroup.AuthGroup.POST("/login", authEndpoint.Login)

	allusers, err := authService.GetAllUsers()
	if err != nil {
		panic(err)
	}

	for _, user := range allusers {
		userdatabase, err := database.GetDB(user)
		if err != nil {
			panic(err)
		}
		err = dependencies.MigrateUserTables(userdatabase)
		if err != nil {
			panic(err)
		}

		dependencies.SeedPaymentProperties(userdatabase)
	}

	protectedGroup := routes.ProtectedGroup(r)
	protectedGroup.OrderGroup.POST("/order", orderEndpoint.CreateOrder)
	protectedGroup.OrderGroup.GET("/order/:order_id", orderEndpoint.GetOrderByID)

	protectedGroup.PaymentGroup.GET("/payment_method/:method_id", configService.GetPaymentMethod)
	protectedGroup.PaymentGroup.GET("/payment_gateway/:gateway_id", configService.GetPaymentGateway)
	protectedGroup.PaymentGroup.POST("/payment_config", configService.PaymentConfiguration)

	r.Run(":9111")
}
