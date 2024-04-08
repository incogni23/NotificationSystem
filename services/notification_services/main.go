package main

import (
	"fmt"
	"sync"

	database "github.com/database"
	clock "github.com/notification_services/timer"
	service "github.com/notification_services/webhook"
	log "github.com/rs/zerolog/log"
)

func main() {
	db, err := database.SetupEnvAndDB()
	if err != nil {
		log.Fatal().Msg(fmt.Errorf("Failed to connect to database: %s", err).Error())
	}

	err = database.Migrate(db)
	if err != nil {
		log.Fatal().Msg(fmt.Errorf("failed to automigrate database: %s", err).Error())
	}

	cfg, err := database.GetEnv()
	if err != nil {
		log.Fatal().Msg(fmt.Errorf("error getting env: %s", err).Error())
	}

	go service.Consuming(cfg.BrokerAddress)

	go clock.RetryClock()

	// not exiting the server
	var wg sync.WaitGroup

	wg.Wait()
}
