package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Kai-Kota/ChatApp2/backend/store"
	"github.com/Kai-Kota/ChatApp2/backend/types"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
)

func handler(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionID := request.RequestContext.ConnectionID

	// リクエストボディをパース
	var req types.SendMessageRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"message":"Invalid request body"}`,
		}, nil
	}

	connectionsTableName := os.Getenv("CONNECTIONS_TABLE")
	messagesTableName := os.Getenv("MESSAGES_TABLE")

	connStore := store.NewConnectionStore(ctx, connectionsTableName)
	msgStore := store.NewMessageStore(ctx, messagesTableName)

	// 送信者の接続情報を取得
	conn, err := connStore.GetConnection(connectionID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"message":"Connection not found"}`,
		}, nil
	}

	// メッセージを保存
	message, err := msgStore.SaveMessage(req.RoomID, conn.UserID, req.Content)
	if err != nil {
		fmt.Printf("Failed to save message: %v\n", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"message":"Failed to save message"}`,
		}, nil
	}

	// 同じ部屋の全接続を取得
	connections, err := connStore.GetConnectionsByRoom(req.RoomID)
	if err != nil {
		fmt.Printf("Failed to get connections: %v\n", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"message":"Failed to get connections"}`,
		}, nil
	}

	// API Gateway Management APIクライアントを作成
	endpoint := fmt.Sprintf("https://%s/%s", request.RequestContext.DomainName, request.RequestContext.Stage)
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-northeast-1"))
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"message":"Failed to load config"}`,
		}, nil
	}

	client := apigatewaymanagementapi.NewFromConfig(cfg, func(o *apigatewaymanagementapi.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	// メッセージを全接続にブロードキャスト
	messageData, _ := json.Marshal(message)
	for _, c := range connections {
		_, err := client.PostToConnection(ctx, &apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: aws.String(c.ConnectionID),
			Data:         messageData,
		})
		if err != nil {
			fmt.Printf("Failed to send to connection %s: %v\n", c.ConnectionID, err)
			// 接続が切れている場合は削除
			connStore.DeleteConnection(c.ConnectionID)
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       `{"message":"Message sent"}`,
	}, nil
}

func main() {
	lambda.Start(handler)
}
