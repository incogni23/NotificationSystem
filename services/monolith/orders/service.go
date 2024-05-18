package orders

import "time"

type Orders struct {
	OrderID     string
	Amount      string
	Status      string
	CreatedAt   time.Time
	CompletedAt time.Time
}


