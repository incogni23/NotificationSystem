package dependencies

import (
	"log"

	"github.com/monolith/payments"
	"gorm.io/gorm"
)

//func Initconfig(db *gorm.DB, config payments.PaymentConfiguration) error {
//res := db.Where("user_id = ?", config.UserID).FirstOrCreate(&config)
//
//if res.Error != nil {
//	log.Printf("error creating payment configuration")
//	return res.Error
//}
//if res.RowsAffected == 0 {
//	log.Printf("Payment configuration already exists")
//} else {
//	log.Printf("Payment configuration successful")
//}
//
//return nil
//}

func SeedPaymentProperties(db *gorm.DB) {

	paymentMethod := payments.PaymentMethod{
		MethodID:   "1",
		MethodType: "CREDIT_CARD",
	}

	if err := db.FirstOrCreate(&paymentMethod, payments.PaymentMethod{MethodID: paymentMethod.MethodID}).Error; err != nil {
		log.Printf("Failed to add payment method: %v\n", err)
	} else {
		log.Printf("Seeded payment method: %v\n", paymentMethod.MethodType)
	}

	paymentGateway := payments.PaymentGateway{
		GatewayID:   "2",
		GatewayName: "JUSPAY",
	}

	if err := db.FirstOrCreate(&paymentGateway, payments.PaymentGateway{GatewayID: paymentGateway.GatewayID}).Error; err != nil {
		log.Printf("Failed to add payment gateway: %v\n", err)
	} else {
		log.Printf("Seeded payment gateway: %v\n", paymentGateway.GatewayName)
	}
}
