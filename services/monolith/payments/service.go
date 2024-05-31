package payments

import (
	"time"
)

type Payment struct {
	PaymentID        string         `json:"payment_id" gorm:"primaryKey"`
	OrderID          string         `json:"order_id"`
	RecievedPID      string         `json:"recieved_pid"`
	PaymentMethodID  string         `json:"method_id" gorm:"not null;"`
	PaymentMethod    PaymentMethod  `gorm:"references:MethodID"`
	PaymentGatewayID string         `json:"gateway_id" gorm:"not null;"`
	PaymentGateway   PaymentGateway `gorm:"references:GatewayID"`
	Amount           float64        `json:"amount"`
	Status           string         `json:"status"`
	CreatedAt        time.Time      `json:"created_at"`
	CompletedAt      time.Time      `json:"completed_at"`
}
