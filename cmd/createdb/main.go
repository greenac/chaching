package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/greenac/chaching/internal/database/helpers"
	"github.com/greenac/chaching/internal/env"
	"github.com/greenac/chaching/internal/service/logger"
	"github.com/spf13/viper"
	"os"
)

func main() {
	log := logger.NewLogger(logger.LogLevelForLogLevelName(os.Getenv("LOG_LEVEL")), os.Getenv("GO_ENV") != string(env.GoEnvLocal))

	log.Info("Running create database...")

	envVars, err := env.NewEnv(".env", viper.New())
	if err != nil {
		log.Error("main:failed to read env file with error: " + err.Error())
		panic(err)
	}

	config := helpers.GetDynamoConfig(helpers.GetDynamoConfigInput{
		MainTable:  envVars.GetString("DYNAMO_MAIN_TABLE_NAME"),
		Env:        env.GoEnv(envVars.GetString("GO_ENV")),
		AwsRegion:  envVars.GetString("AWS_REGION"),
		DynamoUrl:  envVars.GetString("DYNAMO_URL"),
		AwsProfile: os.Getenv("AWS_PROFILE"),
	})

	client, ge := helpers.DynamoClient(context.Background(), config)
	if ge != nil {
		log.Error("main:failed to create dynamo table with error: " + ge.Error())
		panic(ge)
	}

	var input dynamodb.CreateTableInput
	data, err := os.ReadFile("files/dynamo-table.json")
	if err != nil {
		log.Error("main:failed to read dynamo table file with error: " + err.Error())
		panic(err)
	}

	err = json.Unmarshal(data, &input)
	if err != nil {
		log.Error("main:failed to read unmarshal dynamo table data: " + err.Error())
		panic(err)
	}

	_, err = client.CreateTable(context.Background(), &input)
	if err != nil {
		log.Error("main:failed to create dynamo table with error: " + err.Error())
		panic(err)
	}

	log.Info("main:created dynamo table: " + config.MainTable)
}
