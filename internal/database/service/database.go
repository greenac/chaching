package service

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/greenac/chaching/internal/database/models"
	genErr "github.com/greenac/chaching/internal/error"
	"github.com/greenac/chaching/internal/service/database"
	"strconv"
	"time"
)

type IDatabaseService interface {
	SaveDataPoints(ctx context.Context, dps []models.DataPoint) *[]genErr.IGenError
	GetDataPointsInTimeRange(ctx context.Context, companyName string, startDate time.Time, endDate time.Time) ([]models.DataPoint, genErr.IGenError)
}

var _ IDatabaseService = (*DatabaseService)(nil)

func NewDatabaseService(database database.IDatabase[models.DbDataPoint]) IDatabaseService {
	return &DatabaseService{database: database}
}

type DatabaseService struct {
	database database.IDatabase[models.DbDataPoint]
}

func (dbs *DatabaseService) SaveDataPoints(ctx context.Context, dps []models.DataPoint) *[]genErr.IGenError {
	var errs *[]genErr.IGenError
	dbModels := make([]models.DbDataPoint, len(dps))
	for i, dp := range dps {
		dbModels[i] = dp.DatabaseModel()
	}

	err := dbs.database.BatchWrite(ctx, dbModels)
	if err != nil {
		var genErrs []genErr.IGenError
		switch err.(type) {
		case database.BatchWriteError[models.DataPoint]:
			bwe := err.(database.BatchWriteError[models.DataPoint])
			for _, e := range bwe.FailedWrites {
				genErrs = append(genErrs, &genErr.GenError{Messages: []string{e.Error.Error()}})
			}
		default:
			genErrs = []genErr.IGenError{&genErr.GenError{Messages: []string{err.Error()}}}
		}
		errs = &genErrs
	}

	return errs
}

func (dbs *DatabaseService) GetDataPointsInTimeRange(ctx context.Context, companyName string, startDate time.Time, endDate time.Time) ([]models.DataPoint, genErr.IGenError) {
	qry := map[string]types.Condition{
		models.DbPartitionKey: {
			ComparisonOperator: types.ComparisonOperatorEq,
			AttributeValueList: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: models.DataPointModelKeyPk + companyName},
			},
		},
		models.DbSearchKey: {
			ComparisonOperator: types.ComparisonOperatorBetween,
			AttributeValueList: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: models.DataPointModelKeySk + strconv.FormatInt(startDate.Unix(), 10)},
				&types.AttributeValueMemberS{Value: models.DataPointModelKeySk + strconv.FormatInt(endDate.Unix(), 10)},
			},
		},
	}

	dbDataPoints, err := dbs.database.Query(ctx, qry, "")
	if err != nil {
		return []models.DataPoint{}, &genErr.GenError{Messages: []string{err.Error()}}
	}

	dps := make([]models.DataPoint, len(dbDataPoints))
	for i, m := range dbDataPoints {
		dps[i] = m.DataPoint
	}

	return dps, nil
}
