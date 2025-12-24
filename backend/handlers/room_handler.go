package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Kai-Kota/ChatApp2/backend/domain"
	"github.com/Kai-Kota/ChatApp2/backend/types"
	"github.com/aws/aws-lambda-go/events"
)

type IRoomHandler interface {
	CreateRoom(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)
	GetRoom(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)
	GetUserRooms(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)
	HandleOptions(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)
}

type RoomHandler struct {
	roomDomain domain.IRoomDomain
}

func NewRoomHandler(roomDomain domain.IRoomDomain) IRoomHandler {
	return &RoomHandler{
		roomDomain: roomDomain,
	}
}

func (h *RoomHandler) CreateRoom(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	if event.Body == "" {
		return errResponse(http.StatusBadRequest, "empty body"), nil
	}

	var req types.CreateRoomRequest
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return errResponse(http.StatusBadRequest, "invalid JSON"), nil
	}

	// x-user-idヘッダーから取得（API Gateway V2では小文字に正規化される）
	createdBy := ""
	for key, value := range event.Headers {
		if key == "x-user-id" || key == "X-User-Id" {
			createdBy = value
			break
		}
	}
	if createdBy == "" {
		return errResponse(http.StatusUnauthorized, "x-user-id header required"), nil
	}

	if req.UserName == "" {
		return errResponse(http.StatusBadRequest, "user_name is required"), nil
	}

	room, err := h.roomDomain.CreateRoom(req.UserName, createdBy)
	if err != nil {
		return errResponse(http.StatusInternalServerError, err.Error()), nil
	}

	return response(http.StatusCreated, room), nil
}

func (h *RoomHandler) GetRoom(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	roomID, ok := event.PathParameters["room_id"]
	if !ok || roomID == "" {
		return errResponse(http.StatusBadRequest, "room_id required"), nil
	}

	roomData, err := h.roomDomain.GetRoom(roomID)
	if err != nil {
		return errResponse(http.StatusNotFound, err.Error()), nil
	}

	return response(http.StatusOK, roomData), nil
}

func (h *RoomHandler) GetUserRooms(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// x-user-idヘッダーから取得（API Gateway V2では小文字に正規化される）
	userID := ""
	for key, value := range event.Headers {
		if key == "x-user-id" || key == "X-User-Id" {
			userID = value
			break
		}
	}
	if userID == "" {
		return errResponse(http.StatusUnauthorized, "x-user-id header required"), nil
	}

	rooms, err := h.roomDomain.GetUserRooms(userID)
	if err != nil {
		return errResponse(http.StatusInternalServerError, err.Error()), nil
	}

	return response(http.StatusOK, rooms), nil
}

func (h *RoomHandler) HandleOptions(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "http://localhost:3000",
			"Access-Control-Allow-Methods":     "GET, POST, PUT, DELETE, OPTIONS",
			"Access-Control-Allow-Headers":     "Content-Type, Authorization, x-user-id",
			"Access-Control-Allow-Credentials": "true",
		},
		Body:            "",
		IsBase64Encoded: false,
	}, nil
}
