package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB
var cfg *models.ConfigDatabase

func SetupEnvAndDB() (*gorm.DB, error) {
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return db.Debug(), nil
}

func GetEnv() (*models.ConfigDatabase, error) {
	if cfg != nil {
		return cfg, nil
	}

	return nil, errors.New("error getting config")
}

func CreateEvent(event *models.Event) error {
	return db.Create(&event).Error
}

func UpdateEvent(event *models.Event) error {
	updateFields := map[string]interface{}{
		"Status":    event.Status,
		"Attempts":  event.Attempts,
		"NextRetry": event.NextRetry,
	}
	return db.Model(event).Updates(updateFields).Error
}

func GetEventsForRetry() ([]*models.Event, error) {
	var events []*models.Event
	err := db.Where("status = ? AND next_retry <= ?", models.StatusNotCompleted, time.Now().Unix()).Find(&events).Error
	return events, err
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models.Event{})
	if err != nil {
		return err
	}
	return nil
}
