package store

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"golang.org/x/crypto/bcrypt"
)

type IAuthStore interface {
	Signup(userName, password string) error
	Login(userName, password string) (string, error)
}

type AuthStore struct {
	client    *dynamodb.Client
	tableName string
}

func NewAuthStore(ctx context.Context, tableName string) IAuthStore {
	var cfg aws.Config
	var err error

	region := os.Getenv("AWS_REGION")
	endpoint := os.Getenv("DYNAMODB_ENDPOINT")

	fmt.Printf("[DynamoDB Store] Region: %s\n", region)
	fmt.Printf("[DynamoDB Store] Endpoint: %s\n", endpoint)

	// Check for local DynamoDB endpoint
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
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	return &AuthStore{
		client:    client,
		tableName: tableName,
	}
}

func (s *AuthStore) Signup(userName, password string) error {
	type userItem struct {
		ID       string `dynamodbav:"id"`
		UserName string `dynamodbav:"user_name"`
		Password string `dynamodbav:"password"`
	}

	item := userItem{ID: userName, UserName: userName, Password: password}
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = s.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName:                aws.String(s.tableName),
		Item:                     av,
		ConditionExpression:      aws.String("attribute_not_exists(#id)"),
		ExpressionAttributeNames: map[string]string{"#id": "id"},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthStore) Login(userName, password string) (string, error) {
	// Fetch user by id
	out, err := s.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key: map[string]ddbtypes.AttributeValue{
			"id": &ddbtypes.AttributeValueMemberS{Value: userName},
		},
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}
	if len(out.Item) == 0 {
		return "", errors.New("user not found")
	}

	type userItem struct {
		ID       string `dynamodbav:"id"`
		UserName string `dynamodbav:"user_name"`
		Password string `dynamodbav:"password"`
	}
	var item userItem
	if err := attributevalue.UnmarshalMap(out.Item, &item); err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(item.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Return a simple token placeholder
	return fmt.Sprintf("token:%s", userName), nil
}
