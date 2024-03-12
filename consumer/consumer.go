package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"fmt"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewReader(brokers []string, Topic string) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokers,
		Topic:     Topic,
		Partition: 0,
	})
	r.SetOffset(6)

	return &Consumer{
		reader: r,
	}
}

type Event struct {
	Message            []byte   `json:"message"`
	Source             string   `json:"source"`
	DestinationAddress string   `json:"destinationAddress"`
	Topics             []string `json:"topics"`
}

func conversion(msgs []byte) (Event, error) {
	var event Event
	err := json.Unmarshal(msgs, &event)
	if err != nil {
		return Event{}, err
	}
	return event, nil
}

func (c *Consumer) Consume() {

	for {
		m, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Print("error is", err)
			continue
		}
		event, err := conversion(m.Value)
		if err != nil {
			fmt.Print("error after getting the msg is ", err)
			continue
		}
		fmt.Print("so far msg is", event.Message)
		deliver(event)
		if err != nil {
			fmt.Print("error in sending msg to webhook", err)
			continue
		}

	}
}

func deliver(event Event) (err error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", event.DestinationAddress, bytes.NewBuffer(event.Message))
	if err != nil {
		fmt.Print("error creating HTTP req", err)
		return err

	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Print("error in send msg to webhook", err)
		return err

	}
	defer resp.Body.Close()
	fmt.Print("message sent", resp.Status)
	return

}
