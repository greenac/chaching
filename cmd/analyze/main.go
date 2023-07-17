package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/greenac/chaching/internal/consts"
	"github.com/greenac/chaching/internal/database/helpers"
	"github.com/greenac/chaching/internal/database/managers"
	"github.com/greenac/chaching/internal/database/models"
	"github.com/greenac/chaching/internal/database/service"
	"github.com/greenac/chaching/internal/env"
	"github.com/greenac/chaching/internal/service/analysis"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"math"
	"os"
	"sort"
	"time"
)

const AppleStartStockAmount = 10

func main() {
	var logger zerolog.Logger
	if os.Getenv("GO_ENV") == string(env.GoEnvLocal) {
		logger = zerolog.New(os.Stdout).With().Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout})
	} else {
		logger = zerolog.New(os.Stdout).With().Logger()
	}

	logger.Info().Msg("Running create analyze")

	envVars, err := env.NewEnv(".env", viper.New())
	if err != nil {
		logger.Error().Msg("main:failed to read env file with error: " + err.Error())
		panic(err)
	}

	config := models.DynamoConfig{
		MainTable: envVars.GetString("DynamoMainTableName"),
		Env:       env.GoEnv(envVars.GetString("GO_ENV")),
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

	dbService := service.NewDatabaseService(
		managers.NewDataPointPersistenceManager(
			&managers.DynamoPersistenceManager{
				Client:         client,
				Ctx:            context.Background(),
				Config:         config,
				AttrMarshaller: attributevalue.MarshalMap,
			},
		),
	)

	startDate, err := time.Parse(time.RFC3339, "2022-09-07T06:29:00-07:00")
	if err != nil {
		panic(err)
	}

	endDate, err := time.Parse(time.RFC3339, "2022-09-07T15:01:00-07:00")
	if err != nil {
		panic(err)
	}

	dps, err := dbService.GetDataPointsInTimeRange(context.Background(), consts.Apple, startDate, endDate)
	if err != nil {
		logger.Error().Msg("main:failed to retrieve data with error: " + err.Error())
		panic(err)
	}

	dataPts := make([]models.DataPoint, len(dps))
	for i, dp := range dps {
		dataPts[i] = dp
	}

	logger.Info().Msgf("main:got data points %d", len(dataPts))

	analysisService := analysis.NewAnalysisService()

	slopeChanges := analysisService.FindSlopeChanges(dataPts)

	normFactor := analysisService.CalcSlopeNormalizationFactor(slopeChanges)

	for i, sc := range slopeChanges {
		sc.NormalizeSlope(normFactor)
		slopeChanges[i] = sc
	}

	sort.Slice(slopeChanges, func(i, j int) bool {
		return slopeChanges[i].Time.Before(slopeChanges[j].Time)
	})

	mp := analysisService.SlopeAbsMidpoint(slopeChanges)

	logger.Info().Msgf("midpoint: %f", mp)

	if len(slopeChanges) == 0 {
		logger.Error().Msg("main:no slope changes")
		return
	}

	sc1 := slopeChanges[0]

	logger.Info().Msgf("main: open price: %f, end price: %f, day diff %f", sc1.OpenPrice, sc1.ClosePrice, sc1.ClosePrice-sc1.OpenPrice)

	var sType models.StockSaleType
	if sc1.NormalizedSlope > 0 {
		sType = models.StockSaleTypeBuy
	} else {
		sType = models.StockSaleTypeSell
	}

	sales := []models.StockSale{}
	for i := 1; i < len(slopeChanges); i += 1 {
		sc := slopeChanges[i]
		if math.Abs(sc.NormalizedSlope) >= mp {
			t := models.StockSaleTypeBuy
			if sc.NormalizedSlope < 0 {
				t = models.StockSaleTypeSell
			}

			if sType != t {
				sType = t
				sales = append(sales, models.StockSale{
					Name:   consts.Amazon,
					Amount: AppleStartStockAmount,
					Price:  sc.HighestPrice,
					Type:   t,
				})
			}
		}
	}

	amount := sc1.OpenPrice
	for _, s := range sales {
		if s.Type == models.StockSaleTypeBuy {
			amount -= s.Price
			logger.Info().Msgf("buying at: %f amount: %f", s.Price, amount)
		} else {
			amount += s.Price
			logger.Info().Msgf("selling at: %f amount: %f", s.Price, amount)
		}
	}

	logger.Info().Msgf("main:amount is %f", amount)

	//plotService := analysis.PlotService{}
	//err = plotService.Plot(dps)
	//if err != nil {
	//	panic(err)
	//}

	logger.Info().Msg("main:finished analysis")
}
