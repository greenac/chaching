package models

import (
	model "github.com/greenac/chaching/internal/rest/polygon/models"
)

type DataPoint struct {
	BaseDbModel
	model.PolygonDataPoint
	Name string `json:"name" dynamodbav:"name"`
}
