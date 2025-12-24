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
	roomTableName := os.Getenv("ROOM_TABLE")
	membersTableName := os.Getenv("MEMBERS_TABLE")

	room_store := store.NewRoomStore(context.TODO(), roomTableName, membersTableName)
	room_domain := domain.NewRoomDomain(room_store)
	room_handler := handlers.NewRoomHandler(room_domain)

	lambda.Start(func(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		if event.RequestContext.HTTP.Method == "OPTIONS" {
			return room_handler.HandleOptions(ctx, event)
		}
		return room_handler.GetRoom(ctx, event)
	})
}
