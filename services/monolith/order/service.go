package order

import (
	"errors"
	"time"

	"github.com/auth"
	"github.com/google/uuid"
)

type Order struct {
	OrderID     uuid.UUID   `json:"order_id" gorm:"primaryKey"`
	Amount      float64     `json:"amount"`
	Status      OrderStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	CompletedAt time.Time   `json:"completed_at"`
}

type OrderService interface {
	CreateOrder(amount float64, user *auth.User) (*Order, error)
	GetOrderByID(OrderID uuid.UUID, user *auth.User) (*Order, error)
}

type orderService struct {
	dao OrderDao
}

func NewOrderService(dao OrderDao) OrderService {
	return &orderService{
		dao: dao,
	}
}

func (os *orderService) CreateOrder(amount float64, user *auth.User) (*Order, error) {
	if amount <= 0 {
		return nil, errors.New("amount should be more than 0")
	}

	order := &Order{
		OrderID:   uuid.New(),
		Amount:    amount,
		Status:    PENDING,
		CreatedAt: time.Now(),
	}

	err := os.dao.InsertOrder(order, user)
	if err != nil {
		return nil, err
	}

	return order, nil

}

func (os *orderService) GetOrderByID(orderID uuid.UUID, user *auth.User) (*Order, error) {
	return os.dao.GetOrder(orderID, user)
}
