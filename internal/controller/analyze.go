package controller

import (
	"context"
	"errors"
	"github.com/greenac/chaching/internal/consts"
	"github.com/greenac/chaching/internal/database/models"
	"github.com/greenac/chaching/internal/database/service"
	"github.com/greenac/chaching/internal/service/analysis"
	"github.com/greenac/chaching/internal/service/logger"
	"math"
	"sort"
	"time"
)

const AppleStartStockAmount = 10
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

func (ctr *AnalysisController) BuySellInflectionPoint(startDate time.Time, endDate time.Time) (float64, error) {
	dps, err := ctr.databaseService.GetDataPointsInTimeRange(context.Background(), consts.Apple, startDate, endDate)
	if err != nil {
		ctr.logger.Error("main:failed to retrieve data with error: " + err.Error())
		panic(err)
	}

	slopeChanges := ctr.analysisService.FindSlopeChanges(dps)

	ctr.logger.InfoFmt("number of slope changes: %d", len(slopeChanges))

	normFactor := ctr.analysisService.CalcSlopeNormalizationFactor(slopeChanges)

	for i, sc := range slopeChanges {
		sc.NormalizeSlope(normFactor)
		slopeChanges[i] = sc
	}

	sort.Slice(slopeChanges, func(i, j int) bool {
		return slopeChanges[i].Time.Before(slopeChanges[j].Time)
	})

	if len(slopeChanges) == 0 {
		ctr.logger.Error("main:no slope changes")
		return 0, errors.New("no slope change")
	}

	var amounts []SaleAmount
	sellPoint := 0.01

	for sellPoint < 1 {
		sc1 := slopeChanges[0]

		ctr.logger.InfoFmt("main: open price: %f, end price: %f, day diff %f", sc1.OpenPrice, sc1.ClosePrice, sc1.ClosePrice-sc1.OpenPrice)

		var saleType models.StockSaleType
		if sc1.NormalizedSlope > 0 {
			saleType = models.StockSaleTypeBuy
		} else {
			saleType = models.StockSaleTypeSell
		}

		var sales []models.StockSale
		for i := 1; i < len(slopeChanges); i += 1 {
			sc := slopeChanges[i]
			if math.Abs(sc.NormalizedSlope) >= sellPoint {
				st := models.StockSaleTypeBuy
				if sc.NormalizedSlope < 0 {
					st = models.StockSaleTypeSell
				}

				if saleType != st {
					saleType = st
					sales = append(sales, models.StockSale{
						Name:   consts.Apple,
						Amount: AppleStartStockAmount,
						Price:  sc.HighestPrice,
						Type:   st,
					})
				}
			}
		}

		amount := sc1.OpenPrice
		for i, s := range sales {
			if s.Type == models.StockSaleTypeBuy {
				amount -= s.Price
				ctr.logger.InfoFmt("index: %d, buying at: %f amount: %f", i, s.Price, amount)
			} else {
				amount += s.Price
				ctr.logger.InfoFmt("index: %d, selling at: %f amount: %f", i, s.Price, amount)
			}
		}

		amounts = append(amounts, SaleAmount{Amount: amount, SellPoint: sellPoint})

		sellPoint += sellStep
	}

	var maxSellPoint float64 = 0
	var maxAmount float64 = 0
	for i, a := range amounts {
		ctr.logger.InfoFmt("Amount %d = %f, Sell Point: %f", i, a.Amount, a.SellPoint)
		if a.Amount > maxAmount {
			maxAmount = a.Amount
			maxSellPoint = a.SellPoint
		}
	}

	return maxSellPoint, nil
}
