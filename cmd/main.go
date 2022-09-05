package main

import (
	"encoding/json"
	"github.com/greenac/chaching/internal/controller"
	"github.com/greenac/chaching/internal/env"
	rest "github.com/greenac/chaching/internal/rest/client"
	"github.com/greenac/chaching/internal/rest/models"
	model "github.com/greenac/chaching/internal/rest/polygon/models"
	fetch "github.com/greenac/chaching/internal/service"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func main() {
	envVars := env.Env{BaseEnv: viper.New()}

	fc := controller.FetchController{
		Targets:     []string{"AAPL", "AMZN"},
		Start:       time.Now().Add(-4 * 24 * time.Hour),
		Delta:       1 * 24 * time.Hour,
		Unmarshaler: json.Unmarshal,
		FetchService: fetch.FetchService{
			Url: envVars.GetString("PolygonBaseUrl"),
			RestClient: &rest.Client{
				BaseHeaders: &models.Headers{"Authorization": models.HeaderValue{envVars.GetString("PolygonApiKey")}},
				HttpClient:  &http.Client{Timeout: 30 * time.Second},
				BodyReader:  ioutil.ReadAll,
				GetRequest:  http.NewRequest,
			},
			PathJoiner: url.JoinPath,
		},
	}

	fc.RunFetch(controller.FetchParams{TimespanMultiplier: 1, Limit: 1000, Timespan: model.PolygonAggregateTimespanMinute})
}
