package chaching_kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)

type ConsumerState string

const (
	ConsumerStateSuccess ConsumerState = "success"
	ConsumerStateFailed  ConsumerState = "failed"
	ConsumerStateRetry   ConsumerState = "retry"
)

type IKafkaConsumer interface {
	Config() kafka.ReaderConfig
	Close() error
	ReadMessage(ctx context.Context) (kafka.Message, error)
	FetchMessage(ctx context.Context) (kafka.Message, error)
	CommitMessages(ctx context.Context, msgs ...kafka.Message) error
	Offset() int64
	SetOffsetAt(ctx context.Context, t time.Time) error
}

type KafkaConsumerConfig struct {
	Brokers   []string
	Topic     string
	Partition int
	GroupId   string
	GetReader func(config kafka.ReaderConfig) *kafka.Reader
}

func NewKafakConsumer(config KafkaConsumerConfig) IKafkaConsumer {
	return config.GetReader(kafka.ReaderConfig{
		Brokers:   config.Brokers,
		GroupID:   config.GroupId,
		Topic:     config.Topic,
		Partition: config.Partition,
	})
}

type IConsumer[T any, V any] interface {
	Read()
	MessageToDomain(T) (V, error)
	Process(context.Context, V) (ConsumerState, error)
	HandleOutput(ConsumerState, V) error
}
