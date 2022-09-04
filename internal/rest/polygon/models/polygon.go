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