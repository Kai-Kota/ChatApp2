package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Kai-Kota/ChatApp2/backend/store"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionID := request.RequestContext.ConnectionID

	connectionsTableName := os.Getenv("CONNECTIONS_TABLE")
	connStore := store.NewConnectionStore(ctx, connectionsTableName)

	// 接続情報を削除
	err := connStore.DeleteConnection(connectionID)
	if err != nil {
		fmt.Printf("Failed to delete connection: %v\n", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"message":"Failed to delete connection"}`,
		}, nil
	}

	fmt.Printf("WebSocket disconnected: connectionID=%s\n", connectionID)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       `{"message":"Disconnected"}`,
	}, nil
}

func main() {
	lambda.Start(handler)
}
