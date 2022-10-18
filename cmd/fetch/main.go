package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/greenac/chaching/internal/consts"
	"github.com/greenac/chaching/internal/controller"
	"github.com/greenac/chaching/internal/database/helpers"
	"github.com/greenac/chaching/internal/database/managers"
	dynamoModels "github.com/greenac/chaching/internal/database/models"
	"github.com/greenac/chaching/internal/database/service"
	"github.com/greenac/chaching/internal/env"
	rest "github.com/greenac/chaching/internal/rest/client"
	"github.com/greenac/chaching/internal/rest/models"
	model "github.com/greenac/chaching/internal/rest/polygon/models"
	"github.com/greenac/chaching/internal/service/fetch"
	"github.com/greenac/chaching/internal/utils"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	var logger zerolog.Logger
	if os.Getenv("GoEnv") == string(env.GoEnvLocal) {
		logger = zerolog.New(os.Stdout).With().Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout})
	} else {
		logger = zerolog.New(os.Stdout).With().Logger()
	}

	logger.Info().Msg("Running fetch...")

	envVars, err := env.NewEnv(".env", viper.New())
	if err != nil {
		logger.Error().Msg("main:failed to read env file with error: " + err.Error())
		panic(err)
	}

	config := dynamoModels.DynamoConfig{
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

	start, err := time.Parse(time.RFC3339, "2022-09-05T09:30:00-04:00")
	if err != nil {
		logger.Error().Msg("main:failed to parse start time with error: " + err.Error())
		panic(err)
	}

	end, err := time.Parse(time.RFC3339, "2022-09-13T16:00:00-04:00")
	if err != nil {
		logger.Error().Msg("main:failed to parse end time with error: " + err.Error())
		panic(err)
	}

	endOfDay, err := time.Parse(time.RFC3339, "2022-09-05T16:00:00-04:00")
	if err != nil {
		logger.Error().Msg("main:failed to parse end time with error: " + err.Error())
		panic(err)
	}

	fc := controller.FetchController{
		Targets:        []string{consts.Apple, consts.Amazon},
		StartDate:      start,
		EndDate:        end,
		StartOfDay:     start,
		EndOfDay:       endOfDay,
		PartitionValue: time.Minute,
		DatabaseService: service.DatabaseService{
			DataPointPM: &managers.DataPointPersistenceManager{DynamoPersistenceManager: &managers.DynamoPersistenceManager{
				Client:         client,
				Ctx:            context.Background(),
				Config:         config,
				AttrMarshaller: attributevalue.MarshalMap,
			}},
		},
		Logger:      &logger,
		Unmarshaler: json.Unmarshal,
		FetchService: fetch.FetchService{
			Url: envVars.GetString("PolygonBaseUrl"),
			RestClient: &rest.Client{
				BaseHeaders: &models.Headers{"Authorization": models.HeaderValue{"Bearer " + envVars.GetString("PolygonApiKey")}},
				HttpClient:  &http.Client{Timeout: 30 * time.Second},
				BodyReader:  io.ReadAll,
				GetRequest:  http.NewRequest,
			},
			PathJoiner: utils.JoinUrl,
		},
	}

	errs := fc.RunFetch(controller.FetchParams{TimespanMultiplier: 1, Limit: 100, Timespan: model.PolygonAggregateTimespanMinute})
	if errs != nil {
		for _, e := range errs {
			logger.Error().Msg("main:fetching datapoint got error: " + e.Error())
		}
	}

	logger.Info().Msg("main:all done!!!")
}
