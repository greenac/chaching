package model

import (
	"fmt"
	"time"
)

type PolygonAggregateTimespan string

const (
	PolygonAggregateTimespanMinute  PolygonAggregateTimespan = "minute"
	PolygonAggregateTimespanHour    PolygonAggregateTimespan = "hour"
	PolygonAggregateTimespanDay     PolygonAggregateTimespan = "day"
	PolygonAggregateTimespanWeek    PolygonAggregateTimespan = "week"
	PolygonAggregateTimespanMonth   PolygonAggregateTimespan = "month"
	PolygonAggregateTimespanQuarter PolygonAggregateTimespan = "quarter"
	PolygonAggregateTimespanYear    PolygonAggregateTimespan = "year"
)

type PolygonAggregateSortDirection string

const (
	PolygonAggregateSortDirectionAsc PolygonAggregateSortDirection = "asc"
	PolygonAggregateSortDirectionDsc PolygonAggregateSortDirection = "dsc"
)

type PolygonAggregateRequestParams struct {
	Name          string
	Multiplier    int
	Timespan      PolygonAggregateTimespan
	From          time.Time
	To            time.Time
	SortDirection PolygonAggregateSortDirection
	Limit         int
}

func (pr PolygonAggregateRequestParams) Url() string {
	return fmt.Sprintf(
		"%s/range/%d/%s/%d/%d?adjusted=false&sort=%s&limit=%d",
		pr.Name, pr.Multiplier, pr.Timespan, pr.From.UnixMilli(), pr.To.UnixMilli(), pr.SortDirection, pr.Limit,
	)
}

type PolygonDataPoint struct {
	ClosePrice          float64 `json:"c"`  // close price
	HighestPrice        float64 `json:"h"`  // highest price
	LowestPrice         float64 `json:"l"`  // lowest price
	NumOfTxs            int     `json:"n"`  // number of transactions
	OpenPrice           float64 `json:"o"`  // open price
	StartTime           int64   `json:"t"`  // window start unix time stamp (millis)
	Volume              float64 `json:"v"`  // volume
	VolumeWeightedPrice float64 `json:"vw"` // volume weighted ave price
}

type PolygonAggregateResponse struct {
	Adjusted     bool               `json:"adjusted"`
	QueryCount   int                `json:"queryCount"`
	RequestID    string             `json:"request_id"`
	ResultsCount int                `json:"resultsCount"`
	Status       string             `json:"status"`
	Ticker       string             `json:"ticker"`
	DataPoints   []PolygonDataPoint `json:"results"`
}
