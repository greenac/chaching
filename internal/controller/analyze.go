package controller

import (
	"context"
	"github.com/greenac/chaching/internal/database/service"
	error2 "github.com/greenac/chaching/internal/error"
	"github.com/greenac/chaching/internal/service/analysis"
	"github.com/greenac/chaching/internal/service/logger"
	"time"
)

const NumOfStocks = 1
const sellStep float64 = 0.01

type SaleAmount struct {
	Amount    float64
	SellPoint float64
}

func NewAnalysisController(l logger.ILogger, as analysis.IAnalysisService, dbs service.IDatabaseService) *AnalysisController {
	return &AnalysisController{
		analysisService: as,
		databaseService: dbs,
		logger:          l,
	}
}

type AnalysisController struct {
	analysisService analysis.IAnalysisService
	databaseService service.IDatabaseService
	logger          logger.ILogger
}

func (ctr *AnalysisController) BuySellInflectionPoint(company string, startDate time.Time, endDate time.Time) (float64, error) {
	dps, err := ctr.databaseService.GetDataPointsInTimeRange(context.Background(), company, startDate, endDate)
	if err != nil {
		ctr.logger.Error("main:failed to retrieve data with error: " + err.Error())
		return 0, err
	}

	slopeChanges := ctr.analysisService.FindSlopeChanges(dps)

	normFactor := ctr.analysisService.CalcSlopeNormalizationFactor(slopeChanges)

	for i, sc := range slopeChanges {
		sc.NormalizeSlope(normFactor)
		slopeChanges[i] = sc
	}

	var amounts []SaleAmount
	sellPoint := 0.01

	for sellPoint < 1 {
		sales, err := ctr.analysisService.CalcSales(company, NumOfStocks, slopeChanges, sellPoint)
		if err != nil {
			return 0, err
		}

		amounts = append(amounts, SaleAmount{Amount: ctr.analysisService.CalcAmount(slopeChanges[0].OpenPrice, sales), SellPoint: sellPoint})

		sellPoint += sellStep
	}

	var maxSellPoint float64 = 0
	var maxAmount float64 = 0
	for _, a := range amounts {
		if a.Amount > maxAmount {
			maxAmount = a.Amount
			maxSellPoint = a.SellPoint
		}
	}

	return maxSellPoint, nil
}

func (ctr *AnalysisController) InflectionPointsInRange(company string, startDate time.Time, endDate time.Time) ([]float64, error) {
	var genErr error2.IGenError
	var inflections []float64
	date := startDate
	for endDate.After(date) {
		saveSlope := true
		nextDay := date.AddDate(0, 0, 1)
		ip, err := ctr.BuySellInflectionPoint(company, date, nextDay)
		if err != nil {
			if err.Error() != "no slope changes" {
				ctr.logger.Error("AnalysisController->InflectionPointsInRange:failed to get inflection point with error: " + err.Error())
				_ = genErr.AddMsg(err.Error())
			}
			saveSlope = false
		}

		if saveSlope {
			inflections = append(inflections, ip)
		}

		date = nextDay
	}

	return inflections, nil
}

func (ctr *AnalysisController) AmountsByDay(company string, numOfStocks int, sellPoint float64, startDate time.Time, endDate time.Time) (map[time.Time]float64, error) {
	amounts := map[time.Time]float64{}
	date := startDate
	for endDate.After(date) {
		ed := time.Date(date.Year(), date.Month(), date.Day(), endDate.Hour(), endDate.Minute(), 0, 0, date.Location())
		dps, err := ctr.databaseService.GetDataPointsInTimeRange(context.Background(), company, date, ed)
		if err != nil {
			ctr.logger.Error("AnalysisController->AmountsByDay:failed to retrieve data with error: " + err.Error())
			return amounts, err
		}

		slopeChanges := ctr.analysisService.FindSlopeChanges(dps)

		normFactor := ctr.analysisService.CalcSlopeNormalizationFactor(slopeChanges)

		for i, sc := range slopeChanges {
			sc.NormalizeSlope(normFactor)
			slopeChanges[i] = sc
		}

		sales, gErr := ctr.analysisService.CalcSales(company, numOfStocks, slopeChanges, sellPoint)
		if err != nil {
			ctr.logger.Error("AnalysisController->AmountsByDay:failed to calculate sales with error: " + gErr.Error())
			return map[time.Time]float64{}, gErr
		}

		amounts[date] = ctr.analysisService.CalcAmount(slopeChanges[0].OpenPrice, sales)
	}

	return amounts, nil
}
