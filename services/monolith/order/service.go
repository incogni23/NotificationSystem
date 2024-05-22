package order

import "time"

type Order struct {
    OrderID     string    `json:"order_id" gorm:"primaryKey"`
    Amount      float64   `json:"amount"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    CompletedAt time.Time `json:"completed_at"`
}