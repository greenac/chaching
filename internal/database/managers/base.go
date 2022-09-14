package managers

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/greenac/chaching/internal/database/models"
)

const MaxBatchItemCount = 10

type PersistenceManager interface {
	GetKeys() models.ModelKeys
}

type DynamoPersistenceManager struct {
	Client         models.IDatabaseClient
	Ctx            context.Context
	Config         models.DynamoConfig
	AttrMarshaller func(in interface{}) (map[string]types.AttributeValue, error)
}
