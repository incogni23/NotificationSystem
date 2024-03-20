package service

import (
	"bytes"
	"consumer"
	"fmt"
	"net/http"

	"github.com/labstack/gommon/log"
)

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
		isRetryable = false

	default:
		errorMsg := fmt.Sprintf("Unknown notification type: %s", notificationType)
		log.Error(errorMsg)
		isRetryable = false
		return fmt.Errorf(errorMsg), false
	}
	return nil, isRetryable
}
