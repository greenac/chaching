package chaching_kafka

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/greenac/chaching/internal/service/logger"
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

func NewKafkaConsumer(config KafkaConsumerConfig) IKafkaConsumer {
	return config.GetReader(kafka.ReaderConfig{
		Brokers:   config.Brokers,
		GroupID:   config.GroupId,
		Topic:     config.Topic,
		Partition: config.Partition,
	})
}

type KafkaHeaders struct {
	Nonce             uuid.UUID         `json:"nonce"`
	AdditionalHeaders map[string]string `json:"additionalHeaders"`
}

type KafkaMessage[T any] struct {
	Payload T            `json:"payload"`
	Headers KafkaHeaders `json:"headers"`
}

type Message[T any] struct {
	KafkaMessage KafkaMessage[T] `json:"kafkaMessage"`
	CreatedAt    time.Time       `json:"createdAt"`
}

type IConsumer[T any] interface {
	Process(context.Context, KafkaMessage[T]) ConsumerState
	HandleFailedMessage(ctx context.Context, message KafkaMessage[T])
}

func NewBaseConsumer[T any](kc IKafkaConsumer, c IConsumer[T], p IProducer, l logger.ILogger, d time.Duration) *BaseConsumer[T] {
	return &BaseConsumer[T]{
		kafkaConsumer: kc,
		consumer:      c,
		producer:      p,
		logger:        l,
		durTillDql:    d,
	}
}

type BaseConsumer[T any] struct {
	kafkaConsumer IKafkaConsumer
	consumer      IConsumer[T]
	producer      IProducer
	logger        logger.ILogger
	durTillDql    time.Duration
}

func (c *BaseConsumer[T]) Consume() {
	for {
		ctx := context.Background()
		msg, err := c.kafkaConsumer.ReadMessage(ctx)
		if err != nil {
			c.logger.Error("BaseConsumer->ReadMessage:failed to read message with error: " + err.Error())
			continue
		}

		var m Message[T]
		err = json.Unmarshal(msg.Value, &m)
		if err != nil {
			c.logger.Error("BaseConsumer->ReadMessage:failed to unmarshal message with error: " + err.Error())
			continue
		}

		if time.Now().After(m.CreatedAt.Add(c.durTillDql)) {
			// send to dlq
			c.consumer.HandleFailedMessage(ctx, m.KafkaMessage)
		} else {
			switch c.consumer.Process(ctx, m.KafkaMessage) {
			case ConsumerStateSuccess:
			case ConsumerStateRetry:
				err = c.producer.WriteMessages(ctx, msg)
				if err != nil {
					c.logger.Error("BaseConsumer->ReadMessage:failed to unmarshal message with error: " + err.Error())
				}
			case ConsumerStateFailed:
				c.consumer.HandleFailedMessage(ctx, m.KafkaMessage)
			}
		}

		err = c.kafkaConsumer.CommitMessages(ctx, msg)
		if err != nil {
			c.logger.Error("BaseConsumer->ReadMessage:failed to commit message with error: " + err.Error())
		}
	}
}

func (c *BaseConsumer[T]) commit(ctx context.Context, msg kafka.Message) {
	err := c.kafkaConsumer.CommitMessages(ctx, msg)
	if err != nil {
		c.logger.Warn("BaseConsumer->commit:failed with error: " + err.Error())
	}
}

type RetryMessage[T any] struct {
	Topic     string     `json:"topic"`
	Message   Message[T] `json:"message"`
	CreatedAt time.Time  `json:"createdAt"`
}

type RetryConsumer struct {
	kafkaConsumer      IKafkaConsumer
	producer           IProducer
	logger             logger.ILogger
	durTillNextProcess time.Duration
	retryTopic         string
}

func (c *RetryConsumer) Run() {
	for {
		ctx := context.Background()
		msg, err := c.kafkaConsumer.ReadMessage(ctx)
		if err != nil {
			c.logger.Error("RetryConsumer->Run:failed to read message with error: " + err.Error())
			continue
		}

		var retryMessage RetryMessage[any]
		err = json.Unmarshal(msg.Value, &retryMessage)
		if err != nil {
			c.logger.Error("RetryConsumer->Run:failed to unmarshal message with error: " + err.Error())
			continue
		}

		var value []byte
		var topic string
		if time.Now().After(retryMessage.CreatedAt.Add(c.durTillNextProcess)) {
			topic = retryMessage.Topic
			value, err = json.Marshal(retryMessage.Message)
			if err != nil {
				c.logger.Error("RetryConsumer->Run:failed to marshal message with error: " + err.Error())
				continue
			}
		} else {
			topic = c.retryTopic
			value = msg.Value
		}

		err = c.producer.WriteMessages(ctx, kafka.Message{Topic: topic, Value: value})
		if err != nil {
			c.logger.Error("RetryConsumer->Run:failed to write message with error: " + err.Error())
		}

		err = c.kafkaConsumer.CommitMessages(ctx, msg)
		if err != nil {
			c.logger.Error("RetryConsumer->Run:failed to commit message with error: " + err.Error())
		}
	}
}
