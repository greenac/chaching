package database

import (
	"cirello.io/dynamolock/v2"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"strings"
)

// This is the maximum allowed by dynamodb
const maxRecordsToInsertOnBatchOperation = 25
const batchWriteErrorMessage = "batch item failed to write"

type BatchWriteErrorResult[T any] struct {
	Error error
	Input T
}

type BatchWriteError[T any] struct {
	FailedWrites []BatchWriteErrorResult[T]
}

func (e BatchWriteError[T]) Error() string {
	msg := strings.Builder{}
	for i, fw := range e.FailedWrites {
		msg.WriteString(fw.Error.Error())
		if i < len(e.FailedWrites)-1 {
			msg.WriteString(",")
		}
	}

	return msg.String()
}

func (e BatchWriteError[T]) HasErrors() bool {
	return e.FailedWrites != nil && len(e.FailedWrites) > 0
}

type batchItemWriteRequest[T any] struct {
	WriteRequest types.WriteRequest
	Item         T
}

type DynamoConfig struct {
	MainTable string
	Env       string
	Region    string
	Url       string
	Profile   string
	Index1    string
	Index2    string
}

func NewDatabase[T any](
	c IDatabaseClient,
	bi int,
	tn string,
	am func(in interface{}) (map[string]types.AttributeValue, error),
	aum func(map[string]types.AttributeValue, interface{}) error,
) IDatabase[T] {
	return &Database[T]{client: c, numToBatchInsert: bi, tableName: tn, attributeMarshaller: am, attributeUnmarshaler: aum}
}

var _ IDatabase[any] = (*Database[any])(nil)

type Database[T any] struct {
	client               IDatabaseClient
	tableName            string
	numToBatchInsert     int
	attributeMarshaller  func(in interface{}) (map[string]types.AttributeValue, error)
	attributeUnmarshaler func(map[string]types.AttributeValue, interface{}) error
	getLock              func() (*dynamolock.Client, error)
}

func (db *Database[T]) UpsertOne(ctx context.Context, m T) error {
	md, err := db.attributeMarshaller(m)
	if err != nil {
		return err
	}

	// we can expand this to have conditionals as needed
	// we'll probably want to build some conditional functionality into
	// the single table model interface, or as params on the function
	_, err = db.client.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      md,
		TableName: aws.String(db.tableName),
	})

	return err
}

func (db *Database[T]) GetItem(ctx context.Context, key map[string]types.AttributeValue) (T, error) {
	i := &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(db.tableName),
	}

	var item T
	res, err := db.client.GetItem(ctx, i)
	if err != nil {
		return item, err
	}

	err = db.attributeUnmarshaler(res.Item, &item)
	if err != nil {
		return item, err
	}

	return item, nil
}

func (db *Database[T]) Query(ctx context.Context, key map[string]types.Condition, index string) ([]T, error) {
	var items []T
	var startKey map[string]types.AttributeValue
	for {
		itms, sk, err := db.QueryWithLimit(ctx, key, startKey, index, nil)
		if err != nil {
			return []T{}, err
		}

		items = append(items, itms...)

		if sk == nil || len(sk) == 0 {
			break
		} else {
			startKey = sk
		}
	}

	return items, nil
}

func (db *Database[T]) QueryWithLimit(ctx context.Context, key map[string]types.Condition, startKey map[string]types.AttributeValue, index string, limit *int32) ([]T, map[string]types.AttributeValue, error) {
	qi := dynamodb.QueryInput{
		TableName:         aws.String(db.tableName),
		KeyConditions:     key,
		ExclusiveStartKey: startKey,
		Limit:             limit,
	}

	if index != "" {
		qi.IndexName = aws.String(index)
	}

	var items []T
	res, err := db.client.Query(ctx, &qi)
	if err != nil {
		return items, nil, err
	}

	for _, i := range res.Items {
		item := new(T)
		err = db.attributeUnmarshaler(i, item)
		if err != nil {
			return items, nil, errors.New(fmt.Sprintf("Error unmarshaling values: %+v and error: %s", res.Items, err.Error()))
		}

		items = append(items, *item)
	}

	return items, res.LastEvaluatedKey, nil
}

// BatchWrite
// This function inserts items of the database type into the database in batches
//
// Errors:
//
// error # 1: a systematic failure error.
//
// error #2: a batch error. This can occur for two reasons.
// 1. There is an error inserting an entire batch
// 2. An item or items fails to be insert in the batch operation. This result would be returned in `BatchItemOutput`
// We can return an error that maps failed insertions to the original input item.
//
// It will be up to the caller to check what type of error is returned and to implement appropriate error handling
func (db *Database[T]) BatchWrite(ctx context.Context, items []T) error {
	batches, err := db.createBatches(items)
	if err != nil {
		return err
	}

	batchErr := BatchWriteError[T]{}
	// TODO: we can run this concurrently if we need to speed this up,
	// but I think we should only do this if we notice a bottleneck here
	for _, batch := range batches {
		requests := make([]types.WriteRequest, len(batch))
		for i, b := range batch {
			requests[i] = b.WriteRequest
		}

		input := dynamodb.BatchWriteItemInput{RequestItems: map[string][]types.WriteRequest{db.tableName: requests}}
		output, err := db.client.BatchWriteItem(ctx, &input)
		if err != nil {
			fmt.Printf("batch write error: %#v %s\n", err, err.Error())
			for _, b := range batch {
				batchErr.FailedWrites = append(batchErr.FailedWrites, BatchWriteErrorResult[T]{Error: err, Input: b.Item})
			}
			continue
		}

		failedWrites, has := output.UnprocessedItems[db.tableName]
		if has {
			for _, f := range failedWrites {
				var v T
				err := db.attributeUnmarshaler(f.PutRequest.Item, &v)
				if err != nil {
					// TODO: 🤷‍♀️
				}

				batchErr.FailedWrites = append(batchErr.FailedWrites, BatchWriteErrorResult[T]{Error: errors.New(batchWriteErrorMessage), Input: v})
			}
		}
	}

	if len(batchErr.FailedWrites) > 0 {
		return batchErr
	}

	return nil
}

func (db *Database[T]) createBatches(items []T) ([][]batchItemWriteRequest[T], error) {
	numToBatch := maxRecordsToInsertOnBatchOperation
	if db.numToBatchInsert > 0 || db.numToBatchInsert < maxRecordsToInsertOnBatchOperation {
		numToBatch = db.numToBatchInsert
	}

	bins := len(items) / numToBatch
	if len(items)%numToBatch != 0 {
		bins += 1
	}

	toInsert := make([][]batchItemWriteRequest[T], bins)
	j := 0
	for i, it := range items {
		marItem, err := db.attributeMarshaller(it)
		if err != nil {
			return [][]batchItemWriteRequest[T]{}, err
		}
		wr := types.WriteRequest{PutRequest: &types.PutRequest{Item: marItem}}
		if i%numToBatch == 0 && i > 0 {
			j += 1
			toInsert[j] = []batchItemWriteRequest[T]{{WriteRequest: wr, Item: it}}
		} else {
			sub := toInsert[j]
			sub = append(sub, batchItemWriteRequest[T]{WriteRequest: wr, Item: it})
			toInsert[j] = sub
		}
	}

	return toInsert, nil
}
