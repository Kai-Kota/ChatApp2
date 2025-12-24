package main

import (
	"context"
	"os"

	"github.com/Kai-Kota/ChatApp2/backend/domain"
	"github.com/Kai-Kota/ChatApp2/backend/handlers"
	"github.com/Kai-Kota/ChatApp2/backend/store"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	messagesTableName := os.Getenv("MESSAGES_TABLE")

	msgStore := store.NewMessageStore(context.TODO(), messagesTableName)
	msgDomain := domain.NewMessageDomain(msgStore)
	handler := handlers.NewMessageHandler(msgDomain)

	lambda.Start(func(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		if event.RequestContext.HTTP.Method == "OPTIONS" {
			return handler.HandleOptions(ctx, event)
		}
		return handler.GetMessages(ctx, event)
	})
}
