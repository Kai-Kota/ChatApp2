package types

type Room struct {
	RoomID      string `dynamodbav:"room_id" json:"room_id"`
	RoomName    string `dynamodbav:"room_name" json:"room_name"`
	CreatedBy   string `dynamodbav:"created_by" json:"created_by"`
	CreatedAt   int64  `dynamodbav:"created_at" json:"created_at"`
	MemberCount int    `dynamodbav:"member_count" json:"member_count"`
}

type RoomMember struct {
	RoomID   string `dynamodbav:"room_id" json:"room_id"`
	UserID   string `dynamodbav:"user_id" json:"user_id"`
	JoinedAt int64  `dynamodbav:"joined_at" json:"joined_at"`
	Role     string `dynamodbav:"role" json:"role"` // "admin" or "member"
}

type CreateRoomRequest struct {
	UserName string `json:"user_name"` // 相手のユーザー名
}

type GetRoomResponse struct {
	Room    Room     `json:"room"`
	Members []string `json:"members"`
}
