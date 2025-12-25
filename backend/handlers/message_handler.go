package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Kai-Kota/ChatApp2/backend/domain"
	"github.com/aws/aws-lambda-go/events"
)

type IMessageHandler interface {
	GetMessages(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)
	HandleOptions(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)
}

type MessageHandler struct {
	domain domain.IMessageDomain
}

func NewMessageHandler(d domain.IMessageDomain) IMessageHandler {
	return &MessageHandler{domain: d}
}

func (h *MessageHandler) GetMessages(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	roomID, ok := event.PathParameters["room_id"]
	if !ok || roomID == "" {
		return errResponse(http.StatusBadRequest, "room_id required"), nil
	}
	limit := 20
	if lstr := event.QueryStringParameters["limit"]; lstr != "" {
		if l, err := strconv.Atoi(lstr); err == nil {
			limit = l
		}
	}
	messages, err := h.domain.GetMessages(roomID, limit)
	if err != nil {
		return errResponse(http.StatusInternalServerError, err.Error()), nil
	}
	return response(http.StatusOK, messages), nil
}

func (h *MessageHandler) HandleOptions(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, OPTIONS",
			"Access-Control-Allow-Headers":     "Content-Type, Authorization, x-user-id",
			"Access-Control-Allow-Credentials": "false",
		},
		Body:            "",
		IsBase64Encoded: false,
	}, nil
}
