package models

import (
	model "github.com/greenac/chaching/internal/rest/polygon/models"
	"time"
)

type DataPoint struct {
	BaseDbModel
	model.PolygonDataPoint
	Name string `json:"name" dynamodbav:"name"`
}

func (dp DataPoint) XVal() float64 {
	return float64(dp.StartTime)
}

func (dp DataPoint) YVal() float64 {
	return dp.HighestPrice
}

func (dp DataPoint) Time() time.Time {
	return time.Unix(dp.StartTime/1000, 0)
}
