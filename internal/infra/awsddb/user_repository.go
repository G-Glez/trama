package awsddb

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"trama/internal/presentation/api/auth"
)

type UserItem struct {
	Username     string `dynamodbav:"username"`
	PasswordHash string `dynamodbav:"password_hash"`
	CreatedAt    string `dynamodbav:"created_at"`
}

type UserRepository struct {
	db    *dynamodb.Client
	table string
}

func NewUserRepository(db *dynamodb.Client, table string) *UserRepository {
	return &UserRepository{db: db, table: table}
}

func (r *UserRepository) Create(ctx context.Context, user auth.User) error {
	item, err := attributevalue.MarshalMap(UserItem(user))
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.table),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(username)"),
	})
	if err != nil {
		var condErr *types.ConditionalCheckFailedException
		if errors.As(err, &condErr) {
			return auth.ErrUserExists
		}
		return fmt.Errorf("dynamodb put: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (auth.User, error) {
	result, err := r.db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"username": &types.AttributeValueMemberS{Value: username},
		},
	})
	if err != nil {
		return auth.User{}, fmt.Errorf("dynamodb get: %w", err)
	}
	if result.Item == nil {
		return auth.User{}, auth.ErrUserNotFound
	}

	var item UserItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return auth.User{}, fmt.Errorf("unmarshal: %w", err)
	}

	return auth.User(item), nil
}
