package main

import (
	"encoding/json"
	"fmt"
	"github.com/greenac/chaching/internal/controller"
	"github.com/greenac/chaching/internal/env"
	rest "github.com/greenac/chaching/internal/rest/client"
	"github.com/greenac/chaching/internal/rest/models"
	model "github.com/greenac/chaching/internal/rest/polygon/models"
	fetch "github.com/greenac/chaching/internal/service"
	"github.com/greenac/chaching/internal/utils"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	fmt.Println("env is:", os.Getenv("GoEnv"))

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

	start, err := time.Parse(time.RFC3339, "2022-08-29T09:30:00-04:00")
	if err != nil {
		logger.Error().Msg("main:failed to parse start time with error: " + err.Error())
		panic(err)
	}

	end, err := time.Parse(time.RFC3339, "2022-08-29T16:00:00-04:00")
	if err != nil {
		logger.Error().Msg("main:failed to parse end time with error: " + err.Error())
		panic(err)
	}

	fc := controller.FetchController{
		Targets:     []string{"AAPL", "AMZN"},
		Start:       start,
		End:         end,
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

	data, errs := fc.RunFetch(controller.FetchParams{TimespanMultiplier: 1, Limit: 10, Timespan: model.PolygonAggregateTimespanMinute})
	if errs != nil {
		for _, e := range errs {
			logger.Error().Msg("main:fetching datapoint got error: " + e.Error())
		}
	}

	for _, d := range data {
		logger.Info().Msg(fmt.Sprintf("main:got data points %d", len(d)))
	}

	logger.Info().Msg("main:all done!!!")
}
