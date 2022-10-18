package stocks

import (
	"github.com/greenac/chaching/internal/database/models"
	genErr "github.com/greenac/chaching/internal/error"
)

type IStocksService interface {
	ActOnStock(s models.StockSale) genErr.IGenError
}

type StocksService struct {
	Sales []models.StockSale
}

func (ss *StocksService) ActOnStock(s models.StockSale) genErr.IGenError {
	var e genErr.IGenError
	switch s.Type {
	case models.StockSaleTypeSell:
		e = ss.sellStock(s)
	case models.StockSaleTypeBuy:
		e = ss.buyStock(s)
	}

	return e
}

func (ss *StocksService) sellStock(s models.StockSale) genErr.IGenError {
	ss.Sales = append(ss.Sales, s)
	return nil
}

func (ss *StocksService) buyStock(s models.StockSale) genErr.IGenError {
	ss.Sales = append(ss.Sales, s)
	return nil
}
