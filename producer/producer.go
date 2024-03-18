package producer

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/segmentio/kafka-go"
)

type Event struct {
	Message            []byte         `json:"message"`
	Source             string         `json:"source"`
	DestinationAddress string         `json:"destinationAddress"`
	Topics             pq.StringArray `gorm:"type:text[]" json:"topics"`
	NotificationType   string         `json:"notificationtype"`
	IdempotencyKey     string         `json:"idempotencykey"`
	Status             string         `json:"status"`
	Attempts           int            `json:"attempts"`
	NextRetry          int64          `json:"nextretry"`
}

type producer struct {
	writer *kafka.Writer
}

func Newwriter() *producer {
	w := &kafka.Writer{
		Addr: kafka.TCP("localhost:9092"), //kafka's address here
		//Topic:    Topic,
		//Balancer: &kafka.LeastBytes{}, //use round robin which is by default
	}
	return &producer{
		writer: w,
	}

}

func conversion(events Event) []byte {
	message, _ := json.Marshal(events)
	return message
}

func (producer *producer) Produce(ctx context.Context, event *Event) (err error) {

	for _, topic := range event.Topics {
		event.IdempotencyKey = uuid.New().String()
		finalmessage := kafka.Message{
			Topic: topic,
			Value: conversion(*event),
		}
		err := producer.writer.WriteMessages(ctx, finalmessage)
		if err != nil {
			return err
		}
	}
	return err
}
