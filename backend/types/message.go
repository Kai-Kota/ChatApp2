package types

type Message struct {
	RoomID    string `dynamodbav:"room_id" json:"room_id"`
	Timestamp int64  `dynamodbav:"timestamp" json:"timestamp"`
	UserID    string `dynamodbav:"user_id" json:"user_id"`
	Content   string `dynamodbav:"content" json:"content"`
	MessageID string `dynamodbav:"message_id" json:"message_id"`
}

type WebSocketConnection struct {
	ConnectionID string `dynamodbav:"connection_id" json:"connection_id"`
	RoomID       string `dynamodbav:"room_id" json:"room_id"`
	UserID       string `dynamodbav:"user_id" json:"user_id"`
	ConnectedAt  int64  `dynamodbav:"connected_at" json:"connected_at"`
}

type SendMessageRequest struct {
	Action  string `json:"action"`
	RoomID  string `json:"room_id"`
	Content string `json:"content"`
}

type GetMessagesRequest struct {
	RoomID string `json:"room_id"`
	Limit  int    `json:"limit,omitempty"`
}
