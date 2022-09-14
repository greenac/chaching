package helpers

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/greenac/chaching/internal/database/models"
	genErr "github.com/greenac/chaching/internal/error"
)

func CreateTable(ctx context.Context, client models.IDatabaseClient) genErr.IGenError {
	var input dynamodb.CreateTableInput

	_, err := client.CreateTable(ctx, &input)
	if err != nil {
		return &genErr.GenError{Messages: []string{"CreateTable:Failed to create table with error: " + err.Error()}}
	}

	return nil
}

func DeleteTable(ctx context.Context, client models.IDatabaseClient, tableName string) genErr.IGenError {
	input := dynamodb.DeleteTableInput{TableName: aws.String(tableName)}

	_, err := client.DeleteTable(ctx, &input)
	if err != nil {
		return &genErr.GenError{Messages: []string{"DeleteTable:Failed to delete table with error: " + err.Error()}}
	}

	return nil
}
