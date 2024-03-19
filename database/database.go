package database

import (
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func SetupDB() (*gorm.DB, error) {
	dsn := "host=localhost user=pikapika password=Ankita@2307 dbname=db1 port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
func SetDB(database *gorm.DB) {
	db = database
}

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

func CreateEvent(event *Event) error {
	return db.Create(&event).Error
}

func UpdateEvent(event *Event) error {
	return db.Model(&event).Update("Status", event.Status).Error
}

func GetIncompleteEvents() ([]Event, error) {
	var events []Event
	err := db.Where("status = ?", "Not Completed Yet").Find(&events).Error
	return events, err
}
func Initialize() error {
	err := db.AutoMigrate(&Event{})
	if err != nil {
		return err
	}
	return nil
}
