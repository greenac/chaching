package controller

import (
	"context"
	"encoding/json"
	"github.com/greenac/chaching/internal/service/chaching_kafka"
	"github.com/greenac/chaching/internal/service/database"
	"github.com/greenac/chaching/internal/service/logger"
	"github.com/greenac/chaching/internal/utils"
	"time"
)

type DeadLetterRecord struct {
	Pk   string `json:"pk" dynamodbav:"pk"`
	Sk   string `json:"sk" dynamodbav:"sk"`
	Data string `json:"data" dynamodbav:"data"`
}

type FetchMessage struct {
	Company string    `json:"company"`
	From    time.Time `json:"from"`
	To      time.Time `json:"to"`
}

var _ chaching_kafka.IConsumer[FetchMessage] = (*FetchConsumer)(nil)

type FetchConsumerProcessInput struct {
}

type FetchConsumer struct {
	consumer           chaching_kafka.IKafkaConsumer
	producer           chaching_kafka.IProducer
	deadLetterDatabase database.Database[DeadLetterRecord]
	logger             logger.ILogger
}

func (fc *FetchConsumer) HandleFailedMessage(ctx context.Context, message chaching_kafka.KafkaMessage[FetchMessage]) {
	data, err := json.Marshal(message)
	if err != nil {
		fc.logger.Error("FetchConsumer->HandleFailedMessage:failed to marshal message with error: " + err.Error())
		return
	}

	err = fc.deadLetterDatabase.UpsertOne(ctx, DeadLetterRecord{
		Pk:   "company#" + message.Payload.Company,
		Sk:   "from#" + message.Payload.From.Format(time.RFC3339),
		Data: string(data),
	})
	if err == nil {
		fc.logger.Error("FetchConsumer->HandleFailedMessage:saved message successfully")
	} else {
		fc.logger.Error("FetchConsumer->HandleFailedMessage:failed to save message with error: " + err.Error())
	}
}

func (fc *FetchConsumer) Process(ctx context.Context, message chaching_kafka.KafkaMessage[FetchMessage]) chaching_kafka.ConsumerState {
	return fc.process(utils.AddLoggerToCtx(ctx, fc.logger, map[string]string{"nonce": message.Headers.Nonce.String()}), message.Payload)
}

func (fc *FetchConsumer) process(ctx context.Context, message FetchMessage) chaching_kafka.ConsumerState {
	return chaching_kafka.ConsumerStateSuccess
}
