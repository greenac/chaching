package models

import (
	"github.com/greenac/chaching/internal/rest/polygon/models"
	"strconv"
	"time"
)

const (
	DataPointModelKeyPk string = "type#dataPoint#name#"
	DataPointModelKeySk string = "timeStamp#"
)

type DbDataPoint struct {
	DataPoint
	DataBaseModel
}

type DataPoint struct {
	model.PolygonDataPoint
	CompanyName string    `json:"companyName" dynamodbav:"companyName"`
	CreatedAt   time.Time `json:"createdAt,omitempty" dynamodbav:"omitempty,createdAt"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty" dynamodbav:"omitempty,updatedAt"`
}

func (dp *DataPoint) XVal() float64 {
	return float64(dp.StartTime)
}

func (dp *DataPoint) YVal() float64 {
	return dp.HighestPrice
}

func (dp *DataPoint) Time() time.Time {
	return time.Unix(dp.StartTime/1000, 0)
}

func (dp *DataPoint) DatabaseModel() DbDataPoint {
	return DbDataPoint{
		DataPoint: *dp,
		DataBaseModel: DataBaseModel{
			BaseDbModelWith2GlobalKeys{
				BaseDbModelWith1GlobalKeys: BaseDbModelWith1GlobalKeys{
					BaseDbModel: BaseDbModel{
						Pk: DataPointModelKeyPk + dp.CompanyName,
						Sk: DataPointModelKeySk + strconv.FormatInt(dp.StartTime, 10),
					},
				},
			},
		},
	}
}
