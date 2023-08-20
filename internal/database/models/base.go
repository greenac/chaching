package models

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/greenac/chaching/internal/env"
)

type DynamoConfig struct {
	MainTable string
	Env       env.GoEnv
	Region    string
	Url       string
	Profile   string
	Index1    string
	Index2    string
}

type IDatabaseClient interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
	CreateTable(ctx context.Context, params *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error)
	DeleteTable(ctx context.Context, params *dynamodb.DeleteTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteTableOutput, error)
	BatchWriteItem(ctx context.Context, params *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error)
}

type ModelType string

const (
	ModelTypeCompany     ModelType = "company"
	ModelTypeDataPoint   ModelType = "dataPoint"
	ModelTypeTransaction ModelType = "transaction"
)

const (
	DbPartitionKey = "pk"
	DbSearchKey    = "sk"
	DbGsi1Key      = "gsi1"
	DbGsi2Key      = "gsi2"
)

type IDbModel interface {
	Keys() ModelKeys
}

type BaseDbModel struct {
	Pk string `json:"-" dynamodbav:"pk"`
	Sk string `json:"-" dynamodbav:"sk"`
}

type BaseDbModelWith1GlobalKeys struct {
	BaseDbModel
	GPK1 string `dynamodbav:"gpk1,omitempty" json:"-"`
	GSK1 string `dynamodbav:"gsk1,omitempty" json:"-"`
}

type BaseDbModelWith2GlobalKeys struct {
	BaseDbModelWith1GlobalKeys
	GPK2 string `dynamodbav:"gpk2,omitempty" json:"-"`
	GSK2 string `dynamodbav:"gsk2,omitempty" json:"-"`
}

type DataBaseModel struct {
	BaseDbModelWith2GlobalKeys
}

type ModelKeys struct {
	Pk   string
	Sk   string
	Gpk1 string
	Gsk1 string
	Gpk2 string
	Gsk2 string
}

func GetModelKeys(mt ModelType) ModelKeys {
	var mk ModelKeys

	switch mt {
	case ModelTypeCompany:
		mk = ModelKeys{
			Pk: "type#company#",
			Sk: "companyName#",
		}
	case ModelTypeDataPoint:
		mk = ModelKeys{
			Pk: "type#dataPoint#compay#",
			Sk: "timeStamp#",
		}
	case ModelTypeTransaction:
		mk = ModelKeys{
			Pk:   "transaction#",
			Sk:   "companyName#",
			Gpk1: "transaction#",
			Gsk1: "createdAt#",
		}
	}

	return mk
}
