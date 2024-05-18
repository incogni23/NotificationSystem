package payments

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	PaymentID        string
	OrderID          string
	RecievedPID      string
	UserID           uuid.UUID
	PaymentMethodID  string
	PaymentGatewayID string
	Amount           float64
	Status           string
	CreatedAt        time.Time
	CompletedAt      time.Time
}
