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

func consumertoDB(p *consumer.Event) *models.Event {

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
func Deliver(event consumer.Event, notificationType string) (error, bool) {

	var isRetryable bool
	switch notificationType {
	case "webhook":
		client := &http.Client{}

		req, err := http.NewRequest("POST", event.DestinationAddress, bytes.NewBuffer(event.Message))
		if err != nil {
			log.Error("error creating http req", err)
			isRetryable = false
			return err, false
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			log.Error("error sending msg to webhook", err)
			isRetryable = true
			return err, true
		}
		defer resp.Body.Close()
		log.Info("msg sent", resp.Status)
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			isRetryable = true
		} else {
			isRetryable = false
		}

	default:
		errorMsg := fmt.Sprintf("Unknown notification type: %s", notificationType)
		log.Error(errorMsg)
		isRetryable = false
		return fmt.Errorf(errorMsg), false
	}
	return nil, isRetryable
}

//func Automate() {
//	brokeraddress := []string{"localhost:9092"}
//	consumerRecieved := consumer.NewReader(brokeraddress, "quickstart-events")
//	for {
//		consumedEvent := consumerRecieved.Consume()
//		err, _ := Deliver(consumedEvent, consumedEvent.NotificationType)
//		if err != nil {
//			fmt.Print("err in producing msg", err)
//			return
//		}
//	}
//
//}

func Consuming() {
	brokeraddress := []string{"localhost:9092"}
	consumerRecieved := consumer.NewReader(brokeraddress, models.Topicname)
	for {
		consumedEvent := consumerRecieved.Consume()
		//time.Sleep(time.Minute)
		err, isRetryable := Deliver(consumedEvent, consumedEvent.NotificationType)
		if err != nil {
			log.Error("Failed to deliver the event", err)
			if isRetryable {
				dbEvent := consumertoDB(&consumedEvent)
				dbEvent.Status = models.StatusNotCompleted
				dbEvent.Attempts = 1
				dbEvent.NextRetry = time.Now().Unix()
				err := database.CreateEvent(dbEvent)
				if err != nil {
					log.Error("Failed to insert data,err")
					return
				}
				log.Info("Data inserted in DB successsfully")
			} else {
				log.Info("Event not retryable")
			}
		} else {
			log.Info("Event delivered !")
		}
	}
}
