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
	NextRetry          int64          `json:"nextretry"`
}

const (
	StatusNotCompleted = "Not Completed"
	StatusCompleted    = "Completed"
	StatusIgnored      = "Skipped"
)

var MaxRetryAttempts = 5

var Topicname string

type ConfigDatabase struct {
	Port      string `env:"PORT" env-default:"5432"`
	Host      string `env:"HOST" env-default:"localhost"`
	Name      string `env:"NAME" env-default:"db1"`
	User      string `env:"USER" env-default:"pikapika"`
	Password  string `env:"PASSWORD" env-default:"Ankita@2307"`
	Topicname string `env:"TOPICNAME" env-default:"quickstart-events"`
}
