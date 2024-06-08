package order

import (
	"time"

	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	Amount float64 `json:"amount" binding:"required,min=1"`
}
type CreateOrderResponse struct {
	OrderID     uuid.UUID   `json:"order_id"`
	Amount      float64     `json:"amount"`
	Status      OrderStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	CompletedAt time.Time   `json:"complete_at"`
}

type GetOrderResponse struct {
	OrderID     uuid.UUID   `json:"order_id"`
	Amount      float64     `json:"amount"`
	Status      OrderStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	CompletedAt time.Time   `json:"complete_at"`
}

type OrderStatus string

const (
	PENDING   OrderStatus = "pending"
	COMPLETED OrderStatus = "completed"
)
