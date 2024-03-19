package consumer

import (
	"context"
	"encoding/json"

	"fmt"

	"github.com/pikapika/models"
	"github.com/segmentio/kafka-go"
)

type Event = models.Event
type Consumer struct {
	reader *kafka.Reader
}

func NewReader(brokers []string, Topic string) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		GroupID: "consumergc",
		Brokers: brokers,
		Topic:   Topic,
	})

	return &Consumer{
		reader: r,
	}
}

func conversion(msgs []byte) (Event, error) {
	var event Event
	err := json.Unmarshal(msgs, &event)
	if err != nil {
		return Event{}, err
	}
	return event, nil
}

func (c *Consumer) Consume() Event {

	m, err := c.reader.ReadMessage(context.Background())
	if err != nil {
		fmt.Print("error is", err)

	}
	event, err := conversion(m.Value)
	if err != nil {
		fmt.Print("error after getting the msg is ", err)

	}
	fmt.Print("so far msg is", event.Message)
	return event

}
