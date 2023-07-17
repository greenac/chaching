package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
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

	_, err = client.DeleteTable(context.Background(), &dynamodb.DeleteTableInput{TableName: aws.String(config.MainTable)})
	if err != nil {
		log.Error("main:failed to delete dynamo table with error: " + err.Error())
		panic(err)
	}

	log.Info("main:deleted dynamo table: " + config.MainTable)
}
