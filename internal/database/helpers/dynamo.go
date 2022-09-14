package helpers

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	con "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/greenac/chaching/internal/database/models"
	"github.com/greenac/chaching/internal/env"
	genErr "github.com/greenac/chaching/internal/error"
	"strings"
)

const DatabaseKeySeparator = "="

func DynamoClient(ctx context.Context, config models.DynamoConfig) (models.IDatabaseClient, genErr.IGenError) {
	var cfg aws.Config
	var err error

	switch config.Env {
	case env.GoEnvLocal:
		resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           config.Url,
				SigningRegion: "localhost",
			}, nil
		})
		cfg, err = con.LoadDefaultConfig(
			ctx,
			con.WithRegion(config.Region),
			con.WithEndpointResolverWithOptions(resolver),
		)
	case env.GoEnvDev:
		cfg, err = con.LoadDefaultConfig(
			ctx,
			con.WithSharedConfigProfile(config.Profile),
			con.WithRegion(config.Region),
		)
	default:
		cfg, err = con.LoadDefaultConfig(
			ctx,
			con.WithRegion(config.Region),
		)
	}

	if err != nil {
		return nil, &genErr.GenError{Messages: []string{"SetupDatabase::Failed to load config with error: " + err.Error()}}
	}

	return dynamodb.NewFromConfig(cfg), nil
}

func CreateCompositeKey(key string, values ...string) string {
	val := strings.Builder{}
	for i, v := range values {
		if i == len(values)-1 {
			val.WriteString(v)
		} else {
			val.WriteString(v + DatabaseKeySeparator)
		}
	}

	return key + val.String()
}
