package store

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Kai-Kota/ChatApp2/backend/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type IMessageStore interface {
	SaveMessage(roomID, userID, content string) (*types.Message, error)
	GetMessages(roomID string, limit int) ([]types.Message, error)
}

type MessageStore struct {
	client    *dynamodb.Client
	tableName string
}

func NewMessageStore(ctx context.Context, tableName string) IMessageStore {
	var cfg aws.Config
	var err error

	region := "ap-northeast-1"
	endpoint := os.Getenv("DYNAMODB_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8000"
	}

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
	if err != nil {
		panic(err)
	}

	client := dynamodb.NewFromConfig(cfg)

	return &MessageStore{
		client:    client,
		tableName: tableName,
	}
}

func (s *MessageStore) SaveMessage(roomID, userID, content string) (*types.Message, error) {
	message := types.Message{
		RoomID:    roomID,
		Timestamp: time.Now().Unix(),
		UserID:    userID,
		Content:   content,
		MessageID: uuid.New().String(),
	}

	av, err := attributevalue.MarshalMap(message)
	if err != nil {
		return nil, err
	}

	_, err = s.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(s.tableName),
		Item:      av,
	})
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (s *MessageStore) GetMessages(roomID string, limit int) ([]types.Message, error) {
	if limit <= 0 {
		limit = 50
	}

	out, err := s.client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(s.tableName),
		KeyConditionExpression: aws.String("room_id = :rid"),
		ExpressionAttributeValues: map[string]ddbtypes.AttributeValue{
			":rid": &ddbtypes.AttributeValueMemberS{Value: roomID},
		},
		Limit:            aws.Int32(int32(limit)),
		ScanIndexForward: aws.Bool(false), // 降順（新しい順）
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}

	var messages []types.Message
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &messages); err != nil {
		return nil, err
	}

	// 結果を古い順に並び替え（チャットは上から古い順が一般的）
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}
