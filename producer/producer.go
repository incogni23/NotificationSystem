package producer

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type producer struct {
	writer *kafka.Writer
}

func Newwriter(Topic string) *producer {
	w := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    Topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &producer{
		writer: w,
	}

}

type Events struct {
	Msg any
}

func conversion(Events Events) []byte {
	message, _ := json.Marshal(Events.Msg)
	return message

}
func (producer *producer) Produce(ctx context.Context, Events Events) (err error) {
	msgs := conversion(Events)
	var finalmsg kafka.Message
	finalmsg.Value = msgs

	err = producer.writer.WriteMessages(ctx, finalmsg)
	return err

}
