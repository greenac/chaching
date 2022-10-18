package service

import (
	"context"
	"github.com/greenac/chaching/internal/database/managers"
	"github.com/greenac/chaching/internal/database/models"
	genErr "github.com/greenac/chaching/internal/error"
	"time"
)

type IDatabaseService interface {
	SaveDataPoints(ctx context.Context, dps []models.DataPoint) *[]genErr.IGenError
	GetDataPointsInTimeRange(ctx context.Context, name string, startDate time.Time, endDate time.Time) ([]models.DataPoint, genErr.IGenError)
}

var _ IDatabaseService = (*DatabaseService)(nil)

func NewDatabaseService(dataPointPerstanceManager managers.IDataPointPersistenceManager) IDatabaseService {
	return &DatabaseService{DataPointPM: dataPointPerstanceManager}
}

type DatabaseService struct {
	DataPointPM managers.IDataPointPersistenceManager
}

func (dbs *DatabaseService) SaveDataPoints(ctx context.Context, dps []models.DataPoint) *[]genErr.IGenError {
	return dbs.DataPointPM.SaveNewDataPoints(ctx, dps)
}

func (dbs *DatabaseService) GetDataPointsInTimeRange(ctx context.Context, name string, startDate time.Time, endDate time.Time) ([]models.DataPoint, genErr.IGenError) {
	return dbs.DataPointPM.GetDataPointsInTimeRange(ctx, name, startDate, endDate)
}
