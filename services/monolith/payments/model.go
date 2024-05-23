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
	TokenID  string    `json:"token_id" gorm:"primaryKey"`
	UserID   uuid.UUID `json:"user_id" gorm:"type:uuid;not null;constraint:OnDelete:CASCADE;foreignKey:UserID;references:UserID"`
	User     auth.User
	Token    string `json:"token"`
	Provider string `json:"provider"`
}

type PaymentConfiguration struct {
	ConfigID      string    `json:"config_id" gorm:"primaryKey"`
	UserID        uuid.UUID `json:"user_id" gorm:"type:uuid;not null;constraint:OnDelete:CASCADE;foreignKey:UserID;references:UserID"`
	User          auth.User
	MethodID      string `json:"method_id"`
	MethodType    string `json:"method_type"`
	PaymentMethod PaymentMethod
}
