package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO: create a repo here, and pass the auth dependency here to call Signup function
func Signup(c *gin.Context) {
	var user User

	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"respopnse": "done"})

		return
	}

	// signup
	c.JSON(http.StatusOK, gin.H{"respopnse": "done"})
}
