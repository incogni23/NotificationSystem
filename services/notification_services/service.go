package service

import (
	"bytes"
	"consumer"
	"errors"
	"net/http"

	"github.com/labstack/gommon/log"
)

func Deliver(event consumer.Event, notificationType string) error {
	switch notificationType {
	case "webhook":
		client := &http.Client{}

		req, err := http.NewRequest("POST", event.DestinationAddress, bytes.NewBuffer(event.Message))
		if err != nil {
			log.Error("error creating http req", err)
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			log.Error("error sending msg to webhook", err)
			return err
		}
		defer resp.Body.Close()
		log.Info("msg sent", resp.Status)

	default:
		return errors.New("notification type required")
		//fmt.Print("notification type required", notificationType)

	}
	return nil
}
