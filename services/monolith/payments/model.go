package payments

import "github.com/google/uuid"

type PaymentMethod struct {
	MethodID   string
	MethodType string
}

type PaymentGateway struct {
	GatewayID   string
	GatewayName string
}

type ThirdPartyToken struct {
	TokenID  string
	UserID   string
	Token    string
	Provider string
}

type PaymentConfiguration struct {
	ConfigID   string
	UserID     uuid.UUID
	MethodID   string
	MethodType string
}
