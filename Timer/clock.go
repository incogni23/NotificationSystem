package clock

import (
	"consumer"
	"time"

	"github.com/labstack/gommon/log"
	database "github.com/pikapika/database"
	"github.com/pikapika/models"
	service "github.com/pikapika/notification_services"
)

func eventToConsumerEvent(event models.Event) consumer.Event {
	return consumer.Event{
		Message:            event.Message,
		Source:             event.Source,
		DestinationAddress: event.DestinationAddress,
		Topics:             event.Topics,
		NotificationType:   event.NotificationType,
		Status:             event.Status,
		NextRetry:          event.NextRetry,
		Attempts:           event.Attempts,
	}
}

func RetryClock() {
	for {
		events, err := database.GetEventsForRetry()
		if err != nil {
			log.Error("Error getting incomplete events", err)
			continue
		}

		for _, event := range events {
			consumerEvent := eventToConsumerEvent(*event)

			err := service.Deliver(eventToConsumerEvent(consumerEvent), event.NotificationType)

			if err != nil {
				log.Error("Error delivering events", err)
				event.Status = "Not Completed Yet"
				event.Attempts++
				event.NextRetry = time.Now().Unix() + calculateRetryTime(event.Attempts)
				err := database.UpdateEvent(event)
				if err != nil {
					log.Error("Error in updating event", err)
				}
			} else {
				event.Status = "Completed"
				event.Attempts = 0
				event.NextRetry = 0
				err := database.UpdateEvent(event)
				if err != nil {
					log.Error("Error in updating events")
				}
			}
		}
		log.Info("retry is working")
		time.Sleep(time.Minute)

	}
}

func calculateRetryTime(attempts int) int64 {
	switch attempts {
	case 0:
		return 60
	case 1:
		return 120
	default:
		return int64(2*attempts)*60 + 60
	}
}
