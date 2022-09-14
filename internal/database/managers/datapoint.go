package managers

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/greenac/chaching/internal/database/helpers"
	"github.com/greenac/chaching/internal/database/models"
	genErr "github.com/greenac/chaching/internal/error"
	"strconv"
	"time"
)

type IDataPointPersistenceManager interface {
	PersistenceManager
	SaveNewDataPoints(ctx context.Context, dps []models.DataPoint) *[]genErr.IGenError
	GetDataPoints(ctx context.Context, name string) ([]models.DataPoint, genErr.IGenError)
}

type DataPointPersistenceManager struct {
	*DynamoPersistenceManager
}

var _ IDataPointPersistenceManager = (*DataPointPersistenceManager)(nil)

func (pm *DataPointPersistenceManager) GetKeys() models.ModelKeys {
	return models.GetModelKeys(models.ModelTypeDataPoint)
}

func (pm *DataPointPersistenceManager) SaveNewDataPoints(ctx context.Context, dps []models.DataPoint) *[]genErr.IGenError {
	errs := []genErr.IGenError{}
	keys := pm.GetKeys()

	for i := 0; i < len(dps); i += MaxBatchItemCount {
		ri := i + MaxBatchItemCount
		if ri > len(dps) {
			ri = len(dps)
		}

		points := dps[i:ri]
		wrs := make([]types.WriteRequest, len(points))
		for i, m := range points {
			m.Pk = helpers.CreateCompositeKey(keys.Pk, m.Name)
			m.Sk = helpers.CreateCompositeKey(keys.Sk, strconv.FormatInt(m.StartTime, 10))
			m.CreatedAt = time.Now()
			mdps, err := pm.AttrMarshaller(m)
			if err != nil {
				errs = append(errs, &genErr.GenError{Messages: []string{"DataPointPersistenceManager:SaveNewDataPoints:failed to marshal data points with error: " + err.Error()}})
				continue
			}

			wrs[i] = types.WriteRequest{PutRequest: &types.PutRequest{Item: mdps}}
		}

		dbInput := dynamodb.BatchWriteItemInput{RequestItems: map[string][]types.WriteRequest{pm.Config.MainTable: wrs}}
		_, err := pm.Client.BatchWriteItem(ctx, &dbInput)
		if err != nil {
			errs = append(errs, &genErr.GenError{Messages: []string{"DataPointPersistenceManager:SaveNewDataPoints:failed to write data points with error: " + err.Error()}})
		}
	}

	if len(errs) > 0 {
		return &errs
	}

	return nil
}

func (pm *DataPointPersistenceManager) GetDataPoints(ctx context.Context, name string) ([]models.DataPoint, genErr.IGenError) {
	keys := pm.GetKeys()
	qi := dynamodb.QueryInput{
		TableName: aws.String(pm.Config.MainTable),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":partitionKeyValue": &types.AttributeValueMemberS{Value: helpers.CreateCompositeKey(keys.Pk, name)},
			":sortKeyValue":      &types.AttributeValueMemberS{Value: keys.Sk},
		},
		KeyConditionExpression: aws.String("pk = :partitionKeyValue and begins_with (sk, :sortKeyValue)"),
	}

	out, err := pm.Client.Query(ctx, &qi)
	if err != nil {
		return []models.DataPoint{}, &genErr.GenError{Messages: []string{"DataPointPersistenceManager:GetDataPoints:query failed with error: " + err.Error()}}
	}

	var dps []models.DataPoint
	err = attributevalue.UnmarshalListOfMaps(out.Items, &dps)
	if err != nil {
		return []models.DataPoint{}, &genErr.GenError{Messages: []string{"DataPointPersistenceManager:GetDataPoints:unmarshal failed with error: " + err.Error()}}
	}

	return dps, nil
}
