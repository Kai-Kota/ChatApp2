package store

import (
	"context"
	"fmt"
	"os"

	"github.com/Kai-Kota/ChatApp2/backend/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type IConnectionStore interface {
	SaveConnection(connectionID, roomID, userID string, connectedAt int64) error
	DeleteConnection(connectionID string) error
	GetConnectionsByRoom(roomID string) ([]types.WebSocketConnection, error)
	GetConnection(connectionID string) (*types.WebSocketConnection, error)
}

type ConnectionStore struct {
	client    *dynamodb.Client
	tableName string
}

func NewConnectionStore(ctx context.Context, tableName string) IConnectionStore {
	var cfg aws.Config
	var err error

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "ap-northeast-1"
	}
	endpoint := os.Getenv("DYNAMODB_ENDPOINT")

	if endpoint != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				if service == dynamodb.ServiceID {
					return aws.Endpoint{
						URL:               endpoint,
						HostnameImmutable: true,
					}, nil
				}
				return aws.Endpoint{}, &aws.EndpointNotFoundError{}
			})
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(region), config.WithEndpointResolverWithOptions(customResolver))
	} else {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(region))
	}
	if err != nil {
		panic(err)
	}

	client := dynamodb.NewFromConfig(cfg)

	return &ConnectionStore{
		client:    client,
		tableName: tableName,
	}
}

func (s *ConnectionStore) SaveConnection(connectionID, roomID, userID string, connectedAt int64) error {
	conn := types.WebSocketConnection{
		ConnectionID: connectionID,
		RoomID:       roomID,
		UserID:       userID,
		ConnectedAt:  connectedAt,
	}

	av, err := attributevalue.MarshalMap(conn)
	if err != nil {
		return err
	}

	_, err = s.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(s.tableName),
		Item:      av,
	})
	return err
}

func (s *ConnectionStore) DeleteConnection(connectionID string) error {
	_, err := s.client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(s.tableName),
		Key: map[string]ddbtypes.AttributeValue{
			"connection_id": &ddbtypes.AttributeValueMemberS{Value: connectionID},
		},
	})
	return err
}

func (s *ConnectionStore) GetConnectionsByRoom(roomID string) ([]types.WebSocketConnection, error) {
	out, err := s.client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(s.tableName),
		IndexName:              aws.String("RoomConnectionsIndex"),
		KeyConditionExpression: aws.String("room_id = :rid"),
		ExpressionAttributeValues: map[string]ddbtypes.AttributeValue{
			":rid": &ddbtypes.AttributeValueMemberS{Value: roomID},
		},
	})
	if err != nil {
		return nil, err
	}

	var connections []types.WebSocketConnection
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &connections); err != nil {
		return nil, err
	}

	return connections, nil
}

func (s *ConnectionStore) GetConnection(connectionID string) (*types.WebSocketConnection, error) {
	out, err := s.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key: map[string]ddbtypes.AttributeValue{
			"connection_id": &ddbtypes.AttributeValueMemberS{Value: connectionID},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, fmt.Errorf("connection not found")
	}

	var conn types.WebSocketConnection
	if err := attributevalue.UnmarshalMap(out.Item, &conn); err != nil {
		return nil, err
	}

	return &conn, nil
}
