package dependencies

import (
	"github.com/monolith/order"
	"github.com/monolith/payments"
	"gorm.io/gorm"
)

func MigrateUserTables(userDB *gorm.DB) error {

	err := userDB.AutoMigrate(&payments.PaymentConfiguration{},
		&payments.Payment{},
		&payments.ThirdPartyToken{},
		&order.Order{},
	)
	if err != nil {
		return err
	}
	return nil
}
