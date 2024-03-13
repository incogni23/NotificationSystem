package producer

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

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

type Event struct {
	Message            []byte   `json:"message"`
	Source             string   `json:"source"`
	DestinationAddress string   `json:"destinationAddress"`
	Topics             []string `json:"topics"`
	NotificationType   string   `json:"notificationtype"`
}

func conversion(events Event) []byte {
	message, _ := json.Marshal(events)
	return message
}

func (producer *producer) Produce(ctx context.Context, event Event) (err error) {
	//msgs := conversion(Events)
	//var finalmsg kafka.Message
	//finalmsg.Value = msgs
	//err = producer.writer.WriteMessages(ctx, finalmsg)
	//return err
	for _, topic := range event.Topics {
		finalmessage := kafka.Message{
			Topic: topic,
			Value: conversion(event),
		}
		err := producer.writer.WriteMessages(ctx, finalmessage)
		if err != nil {
			return err
		}
	}
	return err
}
