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

type IRoomStore interface {
	CreateRoom(roomName, createdBy string) (*types.Room, error)
	GetRoom(roomID string) (*types.Room, error)
	AddMember(roomID, userID, role string) error
	GetRoomMembers(roomID string) ([]string, error)
	GetUserRooms(userID string) ([]types.Room, error)
}

type RoomStore struct {
	client           *dynamodb.Client
	roomTableName    string
	membersTableName string
}

func NewRoomStore(ctx context.Context, roomTableName, membersTableName string) IRoomStore {
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

	return &RoomStore{
		client:           client,
		roomTableName:    roomTableName,
		membersTableName: membersTableName,
	}
}

func (s *RoomStore) CreateRoom(roomName, createdBy string) (*types.Room, error) {
	roomID := uuid.New().String()
	now := time.Now().Unix()

	room := types.Room{
		RoomID:      roomID,
		RoomName:    roomName,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		MemberCount: 1,
	}

	av, err := attributevalue.MarshalMap(room)
	if err != nil {
		return nil, err
	}

	// Room作成
	_, err = s.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(s.roomTableName),
		Item:      av,
	})
	if err != nil {
		return nil, err
	}

	// 作成者をメンバーとして追加
	if err := s.AddMember(roomID, createdBy, "admin"); err != nil {
		return nil, err
	}

	return &room, nil
}

func (s *RoomStore) GetRoom(roomID string) (*types.Room, error) {
	out, err := s.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(s.roomTableName),
		Key: map[string]ddbtypes.AttributeValue{
			"room_id": &ddbtypes.AttributeValueMemberS{Value: roomID},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, fmt.Errorf("room not found")
	}

	var room types.Room
	if err := attributevalue.UnmarshalMap(out.Item, &room); err != nil {
		return nil, err
	}

	return &room, nil
}

func (s *RoomStore) AddMember(roomID, userID, role string) error {
	member := types.RoomMember{
		RoomID:   roomID,
		UserID:   userID,
		JoinedAt: time.Now().Unix(),
		Role:     role,
	}

	av, err := attributevalue.MarshalMap(member)
	if err != nil {
		return err
	}

	_, err = s.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(s.membersTableName),
		Item:      av,
	})
	return err
}

func (s *RoomStore) GetRoomMembers(roomID string) ([]string, error) {
	out, err := s.client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(s.membersTableName),
		KeyConditionExpression: aws.String("room_id = :rid"),
		ExpressionAttributeValues: map[string]ddbtypes.AttributeValue{
			":rid": &ddbtypes.AttributeValueMemberS{Value: roomID},
		},
	})
	if err != nil {
		return nil, err
	}

	var members []types.RoomMember
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &members); err != nil {
		return nil, err
	}

	userIDs := make([]string, len(members))
	for i, m := range members {
		userIDs[i] = m.UserID
	}

	return userIDs, nil
}

func (s *RoomStore) GetUserRooms(userID string) ([]types.Room, error) {
	// GSIを使ってユーザーが参加している部屋一覧を取得
	out, err := s.client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(s.membersTableName),
		IndexName:              aws.String("UserRoomsIndex"),
		KeyConditionExpression: aws.String("user_id = :uid"),
		ExpressionAttributeValues: map[string]ddbtypes.AttributeValue{
			":uid": &ddbtypes.AttributeValueMemberS{Value: userID},
		},
	})
	if err != nil {
		return nil, err
	}

	var members []types.RoomMember
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &members); err != nil {
		return nil, err
	}

	// 各room_idからRoom情報を取得
	rooms := make([]types.Room, 0, len(members))
	for _, m := range members {
		room, err := s.GetRoom(m.RoomID)
		if err != nil {
			continue
		}
		rooms = append(rooms, *room)
	}

	return rooms, nil
}
