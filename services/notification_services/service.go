package service

import (
	"bytes"
	"consumer"
	"errors"
	"fmt"
	"net/http"
)

func Deliver(event consumer.Event, notificationType string) error {
	switch notificationType {
	case "webhook":
		client := &http.Client{}

		req, err := http.NewRequest("POST", event.DestinationAddress, bytes.NewBuffer(event.Message))
		if err != nil {
			fmt.Print("error creating http req", err)
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Print("error sending msg to webhook", err)
			return err
		}
		defer resp.Body.Close()
		fmt.Print("msg sent", resp.Status)

	default:
		return errors.New("notification type required")
		//fmt.Print("notification type required", notificationType)

	}
	return nil
}

func Automate() {
	brokeraddress := []string{"localhost:9092"}
	consumerRecieved := consumer.NewReader(brokeraddress, "quickstart-events")
	for {
		consumedEvent := consumerRecieved.Consume()
		err := Deliver(consumedEvent, consumedEvent.NotificationType)
		if err != nil {
			fmt.Print("err in producing msg", err)
			return
		}
	}

}
