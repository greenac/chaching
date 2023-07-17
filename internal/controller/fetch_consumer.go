package controller

import (
	"context"
	"github.com/greenac/chaching/internal/service/chaching_kafka"
	"github.com/greenac/chaching/internal/service/database"
	"github.com/greenac/chaching/internal/service/logger"
	"github.com/segmentio/kafka-go"
)

type FetchMessage struct {
	KafkaMessage kafka.Message
}

var _ chaching_kafka.IConsumer[kafka.Message, FetchMessage] = (*FetchConsumer)(nil)

type FetchConsumerProcessInput struct {
}

type FetchConsumer struct {
	consumer           chaching_kafka.IKafkaConsumer
	producer           chaching_kafka.IProducer
	deadLetterDatabase database.Database[]
	logger             logger.ILogger
}

func (fc *FetchConsumer) Read() {
	for {
		// TODO: use custom context
		ctx := context.Background()
		msg, err := fc.consumer.ReadMessage(ctx)
		if err != nil {
			fc.logger.Error("FetchConsumer->Read:failed to read message with error: " + err.Error())
			fc.HandleOutput(chaching_kafka.ConsumerStateFailed, FetchMessage{KafkaMessage: msg})
			continue
		}

		fm, err := fc.MessageToDomain(msg)
		if err != nil {
			fc.HandleOutput(chaching_kafka.ConsumerStateFailed, FetchMessage{KafkaMessage: msg})
			continue
		}

		cs, err := fc.Process(ctx, fm)
		if err != nil {
			fc.HandleOutput(cs, fm)
			continue
		}

		fc.HandleOutput(chaching_kafka.ConsumerStateSuccess, FetchMessage{KafkaMessage: msg})
	}
}

func (fc *FetchConsumer) MessageToDomain(msg kafka.Message) (FetchMessage, error) {
	return FetchMessage{}, nil
}

func (fc *FetchConsumer) Process(ctx context.Context, msg FetchMessage) (chaching_kafka.ConsumerState, error) {
	return chaching_kafka.ConsumerStateSuccess, nil
}

func (fc *FetchConsumer) HandleOutput(state chaching_kafka.ConsumerState, msg FetchMessage) {
	return nil
}
