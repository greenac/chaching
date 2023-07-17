package main

import (
	"github.com/greenac/chaching/internal/consts"
	"github.com/greenac/chaching/internal/env"
	"github.com/greenac/chaching/internal/service/chaching_kafka"
	"github.com/greenac/chaching/internal/service/logger"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"net"
	"os"
	"strconv"
)

func main() {
	log := logger.NewLogger(logger.LogLevelForLogLevelName(os.Getenv("LogLevel")), os.Getenv("GO_ENV") != string(env.GoEnvLocal))

	envVars, err := env.NewEnv(".env", viper.New())
	if err != nil {
		log.Error("main:failed to read env file with error: " + err.Error())
		panic(err)
	}

	topicNames := consts.AllTopics()
	topics := make([]chaching_kafka.Topic, len(topicNames))
	topicConfigs := make([]kafka.TopicConfig, len(topicNames))

	for i, name := range topicNames {
		topics[i] = chaching_kafka.Topic{
			Name:              name,
			Partitions:        consts.NumberOfPartitions,
			ReplicationFactor: consts.ReplicationFactor,
		}

		topicConfigs[i] = kafka.TopicConfig{
			Topic:             name.String(),
			NumPartitions:     consts.NumberOfPartitions,
			ReplicationFactor: consts.ReplicationFactor,
		}
	}

	conn, err := kafka.Dial("tcp", envVars.GetString("KafkaBrokers"))
	if err != nil {
		log.Error("main:failed dial kafka " + err.Error())
		panic(err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		log.Error("main:failed to create controller: " + err.Error())
		panic(err)
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		log.Error("main:failed to read create kafka connection: " + err.Error())
		panic(err)
	}
	defer controllerConn.Close()

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		log.Error("main:failed to read create kafka topics: " + err.Error())
		panic(err.Error())
	}

	log.Info("all done creating topics!")
}
