package chaching_kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type IProducer interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
}
