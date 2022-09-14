package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/greenac/chaching/internal/database/helpers"
	"github.com/greenac/chaching/internal/database/models"
	"github.com/greenac/chaching/internal/env"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"os"
)

func main() {
	var logger zerolog.Logger
	if os.Getenv("GoEnv") == string(env.GoEnvLocal) {
		logger = zerolog.New(os.Stdout).With().Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout})
	} else {
		logger = zerolog.New(os.Stdout).With().Logger()
	}

	logger.Info().Msg("Running create database...")

	envVars, err := env.NewEnv(".env", viper.New())
	if err != nil {
		logger.Error().Msg("main:failed to read env file with error: " + err.Error())
		panic(err)
	}

	config := models.DynamoConfig{
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

	var input dynamodb.CreateTableInput
	data, err := os.ReadFile("files/dynamo-table.json")
	if err != nil {
		logger.Error().Msg("main:failed to read dynamo table file with error: " + err.Error())
		panic(err)
	}

	err = json.Unmarshal(data, &input)
	if err != nil {
		logger.Error().Msg("main:failed to read unmarshal dynamo table data: " + err.Error())
		panic(err)
	}

	_, err = client.CreateTable(context.Background(), &input)
	if err != nil {
		logger.Error().Msg("main:failed to create dynamo table with error: " + err.Error())
		panic(err)
	}

	logger.Info().Msg("main:created dynamo table: " + config.MainTable)
}
