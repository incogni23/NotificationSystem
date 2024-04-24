package routes

import (
	"github.com/gin-gonic/gin"
)

type ProtectedGroups struct {
	PaymentGroup *gin.RouterGroup
}

type UnProtectedGroups struct {
	AuthGroup *gin.RouterGroup
}

func ProtectedGroup(r *gin.Engine) ProtectedGroups {
	paymentGroup := r.Group("/payments")

	paymentGroup.Use(Validate())

	return ProtectedGroups{
		PaymentGroup: paymentGroup,
	}
}

func UnprotectedGroup(r *gin.Engine) UnProtectedGroups {
	authGroup := r.Group("/")

	return UnProtectedGroups{
		AuthGroup: authGroup,
	}
}

func Validate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// call auth validate function
		// case -> success
		// error -> return
		// https://gin-gonic.com/docs/examples/custom-middleware
	}
}
