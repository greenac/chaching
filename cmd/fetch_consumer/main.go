package main

import (
	"github.com/greenac/chaching/internal/consts"
	"github.com/greenac/chaching/internal/env"
	"github.com/greenac/chaching/internal/service/chaching_kafka"
	"github.com/greenac/chaching/internal/service/logger"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func main() {
	log := logger.NewLogger(logger.LogLevelForLogLevelName(os.Getenv("LogLevel")), os.Getenv("GO_ENV") != string(env.GoEnvLocal))

	envVars, err := env.NewEnv(".env", viper.New())
	if err != nil {
		log.Error("main:failed to read env file with error: " + err.Error())
		panic(err)
	}

	consumer := chaching_kafka.NewKafakConsumer(chaching_kafka.KafkaConsumerConfig{
		Brokers:   strings.Split(envVars.GetString("KAFKA_BROKERS"), ","),
		Topic:     consts.TopicNameFetch.String(),
		Partition: envVars.GetInt("FETCH_CONSUMER_PARTITION"),
		GroupId:   envVars.GetString("FETCH_CONSUMER_GROUP_ID"),
		GetReader: kafka.NewReader,
	})

}
