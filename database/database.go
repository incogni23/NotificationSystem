package database

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

var db *gorm.DB

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
	Status             string         `json:"status"`
	Attempts           int            `json:"attempts"`
	NextRetry          int64          `json:"nextRetry"`
}

func CreateEvent(event *Event) error {
	return db.Create(&event).Error
}

func UpdateEvent(event *Event) error {
	return db.Save(event).Error
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
