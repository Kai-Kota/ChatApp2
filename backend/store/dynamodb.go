package store

import (
	"context"
	"fmt"
	"log"

	"github.com/Kai-Kota/ChatApp2/backend/types"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBStore struct {
	client    *dynamodb.Client
	tableName string
}

var _ types.Store = (*DynamoDBStore)(nil)

func NewDynamoDBStore(ctx context.Context, tableName string) *DynamoDBStore {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	return &DynamoDBStore{
		client:    client,
		tableName: tableName,
	}
}

func (d *DynamoDBStore) Put(ctx context.Context, product types.Product) error {
	item, err := attributevalue.MarshalMap(&product)
	if err != nil {
		return fmt.Errorf("unable to marshal product %w", err)
	}

	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &d.tableName,
		Item:      item,
	})

	if err != nil {
		return fmt.Errorf("cannot put item: %w", err)
	}

	return nil
}

func (d *DynamoDBStore) Get(ctx context.Context, id string) (*types.Product, error) {
	response, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &d.tableName,
		Key: map[string]ddbtypes.AttributeValue{
			"id": &ddbtypes.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get item from DynamoDB: %w", err)
	}

	if len(response.Item) == 0 {
		return nil, nil
	}

	product := types.Product{}
	err = attributevalue.UnmarshalMap(response.Item, &product)

	if err != nil {
		return nil, fmt.Errorf("error getting item %w", err)
	}

	return &product, nil
}
