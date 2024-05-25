package payments

import (
	"github.com/auth"
	"github.com/google/uuid"
)

type PaymentMethod struct {
	MethodID   string `json:"method_id" gorm:"primaryKey"`
	MethodType string `json:"method_type"`
}

type PaymentGateway struct {
	GatewayID   string `json:"gateway_id" gorm:"primaryKey"`
	GatewayName string `json:"gateway_name"`
}

type ThirdPartyToken struct {
	TokenID          string         `json:"token_id" gorm:"primaryKey"`
	UserID           uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;"`
	User             auth.User      `gorm:"references:UserID"`
	Token            string         `json:"token"`
	PaymentGatewayID string         `json:"paymentgatewayid" gorm:"not null;"`
	PaymentGateway   PaymentGateway `gorm:"references:GatewayID"`
}

type PaymentConfiguration struct {
	ConfigID         string         `json:"config_id" gorm:"primaryKey"`
	UserID           uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;"`
	User             auth.User      `gorm:"references:UserID"`
	PaymentMethodID  string         `json:"method_id"`
	PaymentGatewayID string         `json:"gateway_id"`
	PaymentGateway   PaymentGateway `gorm:"references:GatewayID"`
	PaymentMethod    PaymentMethod  `gorm:"references:MethodID"`
}
