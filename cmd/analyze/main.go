package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/greenac/chaching/internal/controller"
	"github.com/greenac/chaching/internal/database/helpers"
	"github.com/greenac/chaching/internal/database/models"
	"github.com/greenac/chaching/internal/database/service"
	"github.com/greenac/chaching/internal/env"
	"github.com/greenac/chaching/internal/service/analysis"
	"github.com/greenac/chaching/internal/service/database"
	"github.com/greenac/chaching/internal/service/logger"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"os"
	"time"
)

func main() {
	var zeroLogger zerolog.Logger
	if os.Getenv("GO_ENV") == string(env.GoEnvLocal) {
		zeroLogger = zerolog.New(os.Stdout).With().Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout})
	} else {
		zeroLogger = zerolog.New(os.Stdout).With().Logger()
	}

	log := logger.NewZeroLogWrapper(zeroLogger, logger.LogLevelDebug)
	log.Info("Running create analyze")

	envVars, err := env.NewEnv(".env", viper.New())
	if err != nil {
		log.Error("main:failed to read env file with error: " + err.Error())
		panic(err)
	}

	config := models.DynamoConfig{
		MainTable: envVars.GetString("DYNAMO_MAIN_TABLE_NAME"),
		Env:       env.GoEnv(envVars.GetString("GO_ENV")),
		Region:    envVars.GetString("AWS_REGION"),
		Url:       envVars.GetString("DYNAMO_URL"),
		Profile:   os.Getenv("AWS_PROFILE"),
		Index1:    "ChachingIndex1",
		Index2:    "ChachingIndex2",
	}

	startDate, err := time.Parse(time.RFC3339, "2023-01-01T06:29:00-07:00")
	if err != nil {
		panic(err)
	}

	endDate, err := time.Parse(time.RFC3339, "2023-12-31T15:01:00-07:00")
	if err != nil {
		panic(err)
	}

	client, ge := helpers.DynamoClient(context.Background(), config)
	if ge != nil {
		panic(ge)
	}

	db := database.NewDatabase[models.DbDataPoint](client, 25, envVars.GetString("DYNAMO_MAIN_TABLE_NAME"), attributevalue.MarshalMap, attributevalue.UnmarshalMap)

	analysisController := controller.NewAnalysisController(log, analysis.NewAnalysisService(), service.NewDatabaseService(db))
	tippingPoint, err := analysisController.BuySellInflectionPoint(startDate, endDate)
	if err != nil {
		panic(err)
	}

	log.InfoFmt("tipping point = %f", tippingPoint)

	log.Info("main:finished analysis")
}
