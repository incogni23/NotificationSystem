package service

import (
	"bytes"
	"consumer"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/pikapika/database"
	"github.com/pikapika/models"
)

func consumerAdapter(p *consumer.Event) *models.Event {
	d := models.Event{
		Message:            p.Message,
		Source:             p.Source,
		DestinationAddress: p.DestinationAddress,
		Topics:             p.Topics,
		NotificationType:   p.NotificationType,
		Status:             "",
		Attempts:           0,
		NextRetry:          0,
		IdempotencyKey:     p.IdempotencyKey,
	}

	return &d
}
func DeliverEvent(event consumer.Event, notificationType string) (error, bool) {
	var isRetryable bool

	switch notificationType {
	case "webhook":
		client := &http.Client{}

		req, err := http.NewRequest("POST", event.DestinationAddress, bytes.NewBuffer(event.Message))
		if err != nil {
			log.Error("error creating http req", err)

			// in case of non-recoverable errors we are marking retry as false
			isRetryable = false

			return err, isRetryable
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err == nil {
			log.Info("msg sent", resp.Status)

			// if the acknowledgement from the client fails, we will retry again
			if resp.StatusCode >= 400 {
				isRetryable = true
			}

			return err, isRetryable
		}

		defer resp.Body.Close()

		log.Info("msg sent", resp.Status)

	default:
		errorMsg := fmt.Sprintf("Unknown notification type: %s", notificationType)
		log.Error(errorMsg)
		isRetryable = false
		return fmt.Errorf(errorMsg), false
	}
	return nil, isRetryable
}

func Consuming(brokerAddress []string) {
	consumerRecieved := consumer.NewReader(brokerAddress, models.Topicname)
	for {
		consumedEvent := consumerRecieved.Consume()

		// try to diliver event based on its notification type
		err, isRetryable := DeliverEvent(consumedEvent, consumedEvent.NotificationType)
		if err != nil {
			log.Error("Failed to deliver the event", err)

			if isRetryable {
				dbEvent := consumerAdapter(&consumedEvent)

				dbEvent.Status = models.StatusNotCompleted
				dbEvent.Attempts = 1
				dbEvent.NextRetry = time.Now().Unix()

				err := database.CreateEvent(dbEvent)
				if err != nil {
					log.Error("Failed to insert data,err")
					return
				}

				log.Info("data inserted in db successsfully")
			} else {
				log.Info("event not retryable")
			}
		}
	}
}
