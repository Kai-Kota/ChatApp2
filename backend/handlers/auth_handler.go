package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Kai-Kota/ChatApp2/backend/domain"
	"github.com/Kai-Kota/ChatApp2/backend/types"
	"github.com/aws/aws-lambda-go/events"
)

type IAuthHandler interface {
	Signup(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)
	Login(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)
	HandleOptions(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)
}

type AuthHandler struct {
	authDomain domain.IAuthDomain
}

func NewAuthHandler(authDomain domain.IAuthDomain) IAuthHandler {
	return &AuthHandler{
		authDomain: authDomain,
	}
}

func (h *AuthHandler) Signup(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	if event.Body == "" {
		return errResponse(http.StatusBadRequest, "empty body"), nil
	}

	var input types.Auth
	if err := json.Unmarshal([]byte(event.Body), &input); err != nil {
		return errResponse(http.StatusBadRequest, "invalid JSON"), nil
	}
	if input.UserName == "" || input.Password == "" {
		return errResponse(http.StatusBadRequest, "'user_name' and 'password' required"), nil
	}

	if err := h.authDomain.Signup(input.UserName, input.Password); err != nil {
		return errResponse(http.StatusInternalServerError, err.Error()), nil
	}

	return response(http.StatusCreated, map[string]string{"message": "signup successful"}), nil
}

func (h *AuthHandler) Login(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	if event.Body == "" {
		return errResponse(http.StatusBadRequest, "empty body"), nil
	}

	var input types.Auth
	if err := json.Unmarshal([]byte(event.Body), &input); err != nil {
		return errResponse(http.StatusBadRequest, "invalid JSON"), nil
	}
	if input.UserName == "" || input.Password == "" {
		return errResponse(http.StatusBadRequest, "'user_name' and 'password' required"), nil
	}

	token, err := h.authDomain.Login(input.UserName, input.Password)
	if err != nil {
		return errResponse(http.StatusUnauthorized, err.Error()), nil
	}

	return response(http.StatusOK, map[string]string{"token": token}), nil
}

func (h *AuthHandler) HandleOptions(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "http://localhost:3000",
			"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type, Authorization",
		},
		Body:            "",
		IsBase64Encoded: false,
	}, nil
}

func response(code int, object interface{}) events.APIGatewayV2HTTPResponse {
	marshalled, err := json.Marshal(object)
	if err != nil {
		return errResponse(http.StatusInternalServerError, err.Error())
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: code,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "http://localhost:3000",
			"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type",
		},
		Body:            string(marshalled),
		IsBase64Encoded: false,
	}
}

func errResponse(status int, body string) events.APIGatewayV2HTTPResponse {
	message := map[string]string{
		"message": body,
	}

	messageBytes, _ := json.Marshal(&message)

	return events.APIGatewayV2HTTPResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "http://localhost:3000",
			"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type",
		},
		Body: string(messageBytes),
	}
}
