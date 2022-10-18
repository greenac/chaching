package models

type StockSaleType string

const (
	StockSaleTypeBuy  StockSaleType = "buy"
	StockSaleTypeSell StockSaleType = "sell"
)

type StockSale struct {
	BaseDbModel
	Name   string
	Amount int
	Price  float64
	Type   StockSaleType
}
