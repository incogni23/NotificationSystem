package clock

import (
	"consumer"
	"time"

	"github.com/labstack/gommon/log"
	database "github.com/pikapika/database"
	"github.com/pikapika/models"
	service "github.com/pikapika/notification_services/webhook"
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
		if len(events) == 0 {
			log.Info("No events found for retry.")
			continue
		}
		for _, event := range events {
			consumerEvent := eventToConsumerEvent(*event)
			if event.Attempts >= models.MaxRetryAttempts {

				log.Info("Max attempts reached.")

				event.Status = models.StatusIgnored
				event.NextRetry = 0

				err := database.UpdateEvent(event)
				if err != nil {
					log.Error("Error in updating event", err)
				}

				continue
			}

			err, isRetryable := service.DeliverEvent(eventToConsumerEvent(consumerEvent), event.NotificationType)
			if err != nil {
				log.Error("Error delivering events", err)
				if isRetryable {
					event.Status = models.StatusNotCompleted
					event.Attempts++
					event.NextRetry = time.Now().Unix() + calculateRetryTime(event.Attempts)

					err := database.UpdateEvent(event)
					if err != nil {
						log.Error("Error in updating event", err)
					}

				} else {
					log.Info("Event is not Retryable, terminate!")
					continue
				}
			} else {
				event.Status = models.StatusCompleted
				event.Attempts = 0
				event.NextRetry = 0

				err := database.UpdateEvent(event)
				if err != nil {
					log.Error("Error in updating events")
				}
			}
		}

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
