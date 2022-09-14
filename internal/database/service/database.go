package service

import (
	"context"
	"github.com/greenac/chaching/internal/database/managers"
	"github.com/greenac/chaching/internal/database/models"
	genErr "github.com/greenac/chaching/internal/error"
)

type IDatabaseService interface {
	SaveDataPoints(ctx context.Context, dps []models.DataPoint) *[]genErr.IGenError
}

var _ IDatabaseService = (*DatabaseService)(nil)

type DatabaseService struct {
	Client      models.IDatabaseClient
	DataPointPM managers.IDataPointPersistenceManager
}

func (dbs *DatabaseService) SaveDataPoints(ctx context.Context, dps []models.DataPoint) *[]genErr.IGenError {
	return dbs.DataPointPM.SaveNewDataPoints(ctx, dps)
}

func (dbs *DatabaseService) GetDataPoints(ctx context.Context, name string) ([]models.DataPoint, genErr.IGenError) {
	return dbs.DataPointPM.GetDataPoints(ctx, name)
}
