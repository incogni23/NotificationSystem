package payments

import (
	"time"

	"github.com/auth"
	"github.com/google/uuid"
)

type Payment struct {
	PaymentID        string         `json:"payment_id" gorm:"primaryKey"`
	OrderID          string         `json:"order_id"`
	RecievedPID      string         `json:"recieved_pid"`
	UserID           uuid.UUID      `json:"userID" gorm:"type:uuid;not null;"`
	User             auth.User      `gorm:"references:UserID"`
	PaymentMethodID  string         `json:"method_id" gorm:"not null;"`
	PaymentMethod    PaymentMethod  `gorm:"references:MethodID"`
	PaymentGatewayID string         `json:"gateway_id" gorm:"not null;"`
	PaymentGateway   PaymentGateway `gorm:"references:GatewayID"`
	Amount           float64        `json:"amount"`
	Status           string         `json:"status"`
	CreatedAt        time.Time      `json:"created_at"`
	CompletedAt      time.Time      `json:"completed_at"`
}
