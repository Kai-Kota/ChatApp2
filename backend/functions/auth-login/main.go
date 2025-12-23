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
	tablename := os.Getenv("AUTH_TABLE")

	auth_store := store.NewAuthStore(context.TODO(), tablename)
	auth_domain := domain.NewAuthDomain(auth_store)
	auth_handler := handlers.NewAuthHandler(auth_domain)

	lambda.Start(func(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		if event.RequestContext.HTTP.Method == "OPTIONS" {
			return auth_handler.HandleOptions(ctx, event)
		}
		return auth_handler.Login(ctx, event)
	})
}
