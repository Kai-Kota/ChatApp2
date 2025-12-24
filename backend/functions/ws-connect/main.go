package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Kai-Kota/ChatApp2/backend/store"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionID := request.RequestContext.ConnectionID

	// クエリパラメータからroom_idとuser_idを取得
	roomID := request.QueryStringParameters["room_id"]
	userID := request.QueryStringParameters["user_id"]

	if roomID == "" || userID == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"message":"room_id and user_id are required"}`,
		}, nil
	}

	connectionsTableName := os.Getenv("CONNECTIONS_TABLE")
	connStore := store.NewConnectionStore(ctx, connectionsTableName)

	// 接続情報を保存
	err := connStore.SaveConnection(connectionID, roomID, userID, time.Now().Unix())
	if err != nil {
		fmt.Printf("Failed to save connection: %v\n", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"message":"Failed to save connection"}`,
		}, nil
	}

	fmt.Printf("WebSocket connected: connectionID=%s, roomID=%s, userID=%s\n", connectionID, roomID, userID)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       `{"message":"Connected"}`,
	}, nil
}

func main() {
	lambda.Start(handler)
}
