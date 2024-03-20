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

const (
	StatusNotCompleted = "Not Completed"
	StatusCompleted    = "Completed"
	StatusIgnored      = "Skipped"
)

var MaxRetryAttempts = 5

type ConfigDatabase struct {
	Port     string `env:"PORT" evn-default:"5432"`
	Host     string `env:"HOST" env-default:"localhost"`
	Name     string `env:"NAME" env-default:"postgres"`
	User     string `env:"USER" env-default:"user"`
	Password string `env:"PASSWORD"`
}
