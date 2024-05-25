package configurations

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/monolith/payments"
	"gorm.io/gorm"
)

type DBConfig struct {
	DB *gorm.DB
}

func NewDBConfig(db *gorm.DB) *DBConfig {
	return &DBConfig{
		DB: db,
	}
}

func (dbc *DBConfig) PaymentMethod(c *gin.Context) {

	methodID := c.Param("method_id")

	var paymentMethod payments.PaymentMethod

	if err := dbc.DB.First(&paymentMethod, "method_id = ?", methodID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payement method not found"})
		return
	}

	c.JSON(http.StatusOK, paymentMethod)
}

func (dbc *DBConfig) PaymentGateway(c *gin.Context) {

	gatewayID := c.Param("gateway_id")

	var paymentGateway payments.PaymentGateway

	if err := dbc.DB.First(&paymentGateway, "gateway_id = ?", gatewayID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment gateway not found"})
		return
	}

	c.JSON(http.StatusOK, paymentGateway)
}

func (dbc *DBConfig) PaymentConfiguration(c *gin.Context) {

	var paymentConfig payments.PaymentConfiguration

	err := c.ShouldBindJSON(&paymentConfig)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = dbc.DB.Create(&paymentConfig).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create payment configuration"})
		return
	}

	if err := dbc.DB.Preload("PaymentGateway").Preload("PaymentMethod").Preload("User").First(&paymentConfig, "config_id = ?", paymentConfig.ConfigID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load related data"})
		return
	}

	c.JSON(http.StatusOK, paymentConfig)
}
