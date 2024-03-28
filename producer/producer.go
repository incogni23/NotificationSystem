package producer

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/models"

	"github.com/segmentio/kafka-go"
)

type Event = models.Event
type producer struct {
	writer *kafka.Writer
}

func Newwriter(addr string) *producer {
	w := &kafka.Writer{
		Addr: kafka.TCP(addr), //kafka's address here
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
