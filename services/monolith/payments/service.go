package payments

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	PaymentID   string    `json:"payment_id" gorm:"primaryKey"`
	OrderID     string    `json:"order_id"`
	RecievedPID string    `json:"recieved_pid"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null;constraint:OnDelete:CASCADE;foreignKey:UserID;references:UserID"`
	MethodID    string    `json:"method_id" gorm:"not null;constraint:OnDelete:CASCADE;foreignKey:MethodID;references:MethodID"`
	GatewayID   string    `json:"gateway_id" gorm:"not null;constraint:OnDelete:CASCADE;foreignKey:GatewayID;references:GatewayID"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at"`
}
