package provider

import (
	"context"
	"log/slog"
	"trama/internal/infra/awsddb"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type InfraProvider struct {
	userRepository *awsddb.UserRepository
	dynamo         *dynamodb.Client
}

func (p *Provider) provisionInfra() {
	p.provisionDynamoDB()
	p.provisionUserRepository()
}

func (p *Provider) provisionDynamoDB() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		slog.Error("Failed to load AWS config", "error", err)
		panic(err)
	}
	p.dynamo = dynamodb.NewFromConfig(cfg)
}

func (p *Provider) provisionUserRepository() {
	p.userRepository = awsddb.NewUserRepository(p.dynamo, p.env.DynamoDBUsersTableName)
}
