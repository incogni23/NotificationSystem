package order

import (
	"github.com/auth"
	"github.com/database"
	"github.com/google/uuid"
)

type OrderDao interface {
	InsertOrder(order *Order, user *auth.User) error
	GetOrder(orderID uuid.UUID, user *auth.User) (*Order, error)
}

type orderDao struct {
}

func NewOrderDao() OrderDao {
	return &orderDao{}
}

func (dao *orderDao) InsertOrder(order *Order, user *auth.User) error {
	db, err := database.GetDB(*user)
	if err != nil {
		return err
	}

	return db.Create(order).Error
}

func (dao *orderDao) GetOrder(orderId uuid.UUID, user *auth.User) (*Order, error) {
	db, err := database.GetDB(*user)
	if err != nil {
		return nil, err
	}

	var order Order

	result := db.First(&order, "order_id = ?", orderId)
	if result.Error != nil {
		return nil, result.Error
	}
	return &order, nil
}
