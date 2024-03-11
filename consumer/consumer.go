package consumer

import (
	"context"
	"encoding/json"
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
	r.SetOffset(11)

	return &Consumer{
		reader: r,
	}
}

type Events struct {
	Msg any
}

func conversion(msgs []byte) (interface{}, error) {
	var msg interface{}
	err := json.Unmarshal(msgs, &msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (c *Consumer) Consume() {
	for {
		m, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Print("error is", err)
			continue
		}
		messagerecieved, err := conversion(m.Value)
		if err != nil {
			fmt.Print("error after getting the msg is ", err)
			continue
		}
		fmt.Print("so far msg is", messagerecieved)
	}

}
