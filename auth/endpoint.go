package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type endpoint struct {
	authservices AuthServicer
}

func NewEndpoint(authService AuthServicer) *endpoint {
	return &endpoint{
		authservices: authService}

}

func (e *endpoint) Signup(c *gin.Context) {
	var user User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"response": err.Error()})

		return
	}

	err = e.authservices.SignUp(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"response": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"userId": user.UserID})
}

func (e *endpoint) Login(c *gin.Context) {
	var user User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"response": err.Error()})
		return
	}

	token, err := e.authservices.Login(user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"response": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}

func (e *endpoint) LoginWithToken(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"response": "token not provided"})
		return
	}

	_, err := e.authservices.LoginWithToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"response": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": "valid token"})
}
