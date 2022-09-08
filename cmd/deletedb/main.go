package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/greenac/chaching/internal/database/helpers"
	"github.com/greenac/chaching/internal/env"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"os"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout})
	logger.Info().Msg("Running create database...")

	envVars, err := env.NewEnv(".env", viper.New())
	if err != nil {
		logger.Error().Msg("main:failed to read env file with error: " + err.Error())
		panic(err)
	}

	config := helpers.DynamoConfig{
		MainTable: envVars.GetString("DynamoMainTableName"),
		Env:       env.GoEnv(envVars.GetString("GoEnv")),
		Region:    envVars.GetString("AwsRegion"),
		Url:       envVars.GetString("DynamoUrl"),
		Profile:   os.Getenv("AwsProfile"),
		Index1:    "ChachingIndex1",
		Index2:    "ChachingIndex2",
	}

	client, ge := helpers.DynamoClient(context.Background(), config)
	if ge != nil {
		logger.Error().Msg("main:failed to create dynamo table with error: " + ge.Error())
		panic(ge)
	}

	_, err = client.DeleteTable(context.Background(), &dynamodb.DeleteTableInput{TableName: aws.String(config.MainTable)})
	if err != nil {
		logger.Error().Msg("main:failed to delete dynamo table with error: " + err.Error())
		panic(err)
	}

	logger.Info().Msg("main:deleted dynamo table: " + config.MainTable)
}
