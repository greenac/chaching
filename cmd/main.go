package main

import (
	"encoding/json"
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
	l := zerolog.New(os.Stdout).With().Logger()

	l.Info().Msg("Running fetch...")

	envVars, err := env.NewEnv(".env", viper.New())
	if err != nil {
		l.Error().Msg("main:failed to read env file with error: " + err.Error())
		panic(err)
	}

	fc := controller.FetchController{
		Targets:     []string{"AAPL", "AMZN"},
		Start:       time.Now().Add(-4 * 24 * time.Hour),
		Delta:       1 * 24 * time.Hour,
		Logger:      &l,
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

	fc.RunFetch(controller.FetchParams{TimespanMultiplier: 1, Limit: 10, Timespan: model.PolygonAggregateTimespanMinute})
}
