package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Message            []byte         `json:"message"`
	Source             string         `json:"source"`
	DestinationAddress string         `json:"destinationAddress"`
	Topics             pq.StringArray `gorm:"type:text[]" json:"topics"`
	NotificationType   string         `json:"notificationtype"`
	IdempotencyKey     string         `json:"idempotencykey"`
	Status             string         `json:"status"`
	Attempts           int            `json:"attempts"`
	NextRetry          int64          `json:"nextRetry"`
}
