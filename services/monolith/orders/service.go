package orders

import "time"

type Orders struct {
	OrderID     string
	Amount      float64
	Status      string
	CreatedAt   time.Time
	CompletedAt time.Time
}


