package configurations

import (
	"net/http"
	"strings"

	"github.com/auth"
	"github.com/database"
	"github.com/gin-gonic/gin"
	"github.com/monolith/payments"
	"gorm.io/gorm/clause"
)

type DBConfig struct {
	auth auth.AuthServicer
}

func NewDBConfig(auth auth.AuthServicer) *DBConfig {

	return &DBConfig{
		auth: auth,
	}
}

type PaymentConfigurationResponse struct {
	GatewayName string `json:"gateway_name"`
	MethodName  string `json:"method_name"`
}

func (dbc *DBConfig) GetPaymentMethod(c *gin.Context) {

	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is missing"})
		return
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	userID, err := auth.GetUserIDFromToken(tokenString, "secretKey")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	user := auth.User{UserID: userID}
	userDB, err := database.GetDB(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to user db"})
		return
	}
	methodID := c.Param("method_id")

	var paymentMethod payments.PaymentMethod

	if err := userDB.First(&paymentMethod, "method_id = ?", methodID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payement method not found"})
		return
	}

	c.JSON(http.StatusOK, paymentMethod)
}

func (dbc *DBConfig) GetPaymentGateway(c *gin.Context) {

	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is missing"})
		return
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	userID, err := auth.GetUserIDFromToken(tokenString, "secretKey")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	user := auth.User{UserID: userID}
	userDB, err := database.GetDB(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to user db"})
		return
	}

	gatewayID := c.Param("gateway_id")

	var paymentGateway payments.PaymentGateway

	if err := userDB.First(&paymentGateway, "gateway_id = ?", gatewayID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment gateway not found"})
		return
	}

	c.JSON(http.StatusOK, paymentGateway)
}

func (dbc *DBConfig) PaymentConfiguration(c *gin.Context) {

	var request struct {
		MethodType  string `json:"method_type"`
		GatewayName string `json:"gateway_name"`
	}

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is missing"})
		return
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	userID, err := auth.GetUserIDFromToken(tokenString, "secretKey")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	user, err := dbc.auth.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to user db"})
		return
	}

	userDB, err := database.GetDB(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to user db"})
		return
	}

	if request.MethodType == "" || request.GatewayName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment method and gateway cannot be empty"})
		return
	}

	var paymentMethod payments.PaymentMethod
	result := userDB.First(&paymentMethod, "method_type = ?", request.MethodType)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment method does not exist"})
		return
	}

	var paymentGateway payments.PaymentGateway
	result = userDB.First(&paymentGateway, "gateway_name = ?", request.GatewayName)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment gateway does not exist"})
		return
	}

	paymentConfig := payments.PaymentConfiguration{

		PaymentMethodID:  paymentMethod.MethodID,
		PaymentGatewayID: paymentGateway.GatewayID,
	}

	db := userDB.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{

				{Name: "payment_method_id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{

				"payment_gateway_id",
			}),
		},
	).Create(&paymentConfig)

	if db.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upsert"})
		return
	}

	response := PaymentConfigurationResponse{
		GatewayName: paymentGateway.GatewayName,
		MethodName:  paymentMethod.MethodType,
	}

	c.JSON(http.StatusOK, response)

}
