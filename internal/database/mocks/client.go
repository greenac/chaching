package mocks

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/greenac/chaching/internal/database/models"
)

type ClientMock struct {
	PutItemOutput        dynamodb.PutItemOutput
	PutItemError         error
	QueryOutput          dynamodb.QueryOutput
	QueryError           error
	GetItemOutput        dynamodb.GetItemOutput
	GetItemError         error
	DeleteItemOutput     dynamodb.DeleteItemOutput
	DeleteItemError      error
	UpdateItemOutput     dynamodb.UpdateItemOutput
	UpdateItemError      error
	CreateTableOutput    dynamodb.CreateTableOutput
	CreateTableError     error
	DeleteTableOutput    dynamodb.DeleteTableOutput
	DeleteTableError     error
	BatchWriteItemOutput dynamodb.BatchWriteItemOutput
	BatchWriteItemError  error
}

var _ models.IDatabaseClient = (*ClientMock)(nil)

func (c ClientMock) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return &c.PutItemOutput, c.PutItemError
}

func (c ClientMock) Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	return &c.QueryOutput, c.QueryError
}

func (c ClientMock) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return &c.GetItemOutput, c.GetItemError
}

func (c ClientMock) DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	return &c.DeleteItemOutput, c.DeleteItemError
}

func (c ClientMock) UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	return &c.UpdateItemOutput, c.UpdateItemError
}

func (c ClientMock) CreateTable(ctx context.Context, params *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
	return &c.CreateTableOutput, c.CreateTableError
}

func (c ClientMock) DeleteTable(ctx context.Context, params *dynamodb.DeleteTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteTableOutput, error) {
	return &c.DeleteTableOutput, c.DeleteTableError
}

func (c ClientMock) BatchWriteItem(ctx context.Context, params *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error) {
	return &c.BatchWriteItemOutput, c.BatchWriteItemError
}
