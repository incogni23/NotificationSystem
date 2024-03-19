package database

import (
	"time"

	"github.com/pikapika/models"
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

func CreateEvent(event *models.Event) error {
	return db.Create(&event).Error
}

func UpdateEvent(event *models.Event) error {
	return db.Model(&event).Update("Status", event.Status).Error
}

func GetEventsForRetry() ([]*models.Event, error) {
	var events []*models.Event
	err := db.Where("status = ? AND NextRetry > ?", "Not Completed Yet", time.Now().Unix()).Find(&events).Error
	return events, err
}

func Initialize() error {
	err := db.AutoMigrate(&models.Event{})
	if err != nil {
		return err
	}
	return nil
}
