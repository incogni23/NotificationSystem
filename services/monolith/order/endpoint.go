package order

import (
	"net/http"
	"strings"

	auth "github.com/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderEndpoint struct {
	authService  auth.AuthServicer
	orderService OrderService
}

func NewOrderEndpoint(authService auth.AuthServicer, ordorderService OrderService) *OrderEndpoint {
	return &OrderEndpoint{
		authService:  authService,
		orderService: ordorderService,
	}
}

func (oe *OrderEndpoint) CreateOrder(c *gin.Context) {
	var request CreateOrderRequest
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

	user, err := oe.authService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	orderDao := NewOrderDao()
	orderService := NewOrderService(orderDao)

	order, err := orderService.CreateOrder(request.Amount, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := CreateOrderResponse{
		Amount:      order.Amount,
		OrderID:     order.OrderID,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt,
		CompletedAt: order.CompletedAt,
	}

	c.JSON(http.StatusOK, response)
}

func (oe *OrderEndpoint) GetOrderByID(c *gin.Context) {
	orderIDStr := c.Param("order_id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order_id"})
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

	user, err := oe.authService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	order, err := oe.orderService.GetOrderByID(orderID, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := GetOrderResponse{
		Amount:      order.Amount,
		OrderID:     order.OrderID,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt,
		CompletedAt: order.CompletedAt,
	}

	c.JSON(http.StatusOK, response)

}
