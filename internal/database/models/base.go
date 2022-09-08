package models

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"time"
)

type IClient interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
	CreateTable(ctx context.Context, params *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error)
	DeleteTable(ctx context.Context, params *dynamodb.DeleteTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteTableOutput, error)
}

type ModelType string

const (
	ModelTypeCompany     ModelType = "company"
	ModelTypeDataPoint   ModelType = "dataPoint"
	ModelTypeTransaction ModelType = "transaction"
)

type BaseDbModel struct {
	Pk        string    `json:"-" dynamodbav:"pk"`
	Sk        string    `json:"-" dynamodbav:"sk"`
	CreatedAt time.Time `json:"createdAt,omitempty" dynamodbav:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" dynamodbav:"updatedAt"`
}

type BaseDbModelWith1GlobalKeys struct {
	BaseDbModel
	GPK1 string `dynamodbav:"gpk1" json:"-"`
	GSK1 string `dynamodbav:"gsk1" json:"-"`
}

type BaseDbModelWith2GlobalKeys struct {
	BaseDbModelWith1GlobalKeys
	GPK2 string `dynamodbav:"gpk2" json:"-"`
	GSK2 string `dynamodbav:"gsk2" json:"-"`
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
			Pk: "company#",
			Sk: "companyName#",
		}
	case ModelTypeDataPoint:
		mk = ModelKeys{
			Pk:   "dataPoint#",
			Sk:   "companyName#",
			Gpk1: "companyName#",
			Gsk1: "createdAt#",
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
