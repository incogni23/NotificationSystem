package clock

import (
	"consumer"
	"fmt"
	"time"

	database "github.com/pikapika/database"
	service "github.com/pikapika/notification_services"
)

func eventToConsumerEvent(event database.Event) consumer.Event {
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
		events, err := database.GetIncompleteEvents()
		if err != nil {
			fmt.Print("Error getting incomplete events", err)
			continue
		}

		for _, event := range events {
			if event.NextRetry > time.Now().Unix() {
				continue
			}

			err := service.Deliver(eventToConsumerEvent(event), event.NotificationType)

			if err != nil {
				fmt.Print("Error delivering events", err)
				event.Status = "Not Completed Yet"
				event.Attempts++
				event.NextRetry = time.Now().Unix() + calculateRetryTime(event.Attempts)
				err := database.UpdateEvent(&event)
				if err != nil {
					fmt.Print("Error in updating event", err)
				}
			} else {
				event.Status = "Completed"
				event.Attempts = 0
				event.NextRetry = 0
				err := database.UpdateEvent(&event)
				if err != nil {
					fmt.Print("Error in updating events")
				}
			}
		}
		fmt.Print("retry is working")
		time.Sleep(time.Minute)

	}
}

func calculateRetryTime(attempts int) int64 {
	switch attempts {
	case 0:
		return 0
	case 1:
		return 120
	default:
		return int64(2*attempts)*60 + 60
	}
}
