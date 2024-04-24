package main

import (
	"github.com/auth"
	"github.com/gin-gonic/gin"
	"github.com/monolith/routes"
)

func main() {
	// _, err := dependencies.Init()
	// if err != nil {
	// 	panic(err)
	// }

	r := gin.New()

	unprotectedGroup := routes.UnprotectedGroup(r)

	unprotectedGroup.AuthGroup.POST("/signup", auth.Signup)

	r.Run(":9111")
}
