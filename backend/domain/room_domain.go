package domain

import (
	"errors"

	"github.com/Kai-Kota/ChatApp2/backend/store"
	"github.com/Kai-Kota/ChatApp2/backend/types"
)

type IRoomDomain interface {
	CreateRoom(partnerUserName, createdBy string) (*types.Room, error)
	GetRoom(roomID string) (*types.GetRoomResponse, error)
	GetUserRooms(userID string) ([]types.Room, error)
}

type RoomDomain struct {
	roomStore store.IRoomStore
}

func NewRoomDomain(roomStore store.IRoomStore) IRoomDomain {
	return &RoomDomain{
		roomStore: roomStore,
	}
}

func (d *RoomDomain) CreateRoom(partnerUserName, createdBy string) (*types.Room, error) {
	if partnerUserName == "" {
		return nil, errors.New("user_name is required")
	}
	if createdBy == "" {
		return nil, errors.New("created_by is required")
	}
	// 既存のペア部屋がないかチェック（両ユーザーが参加している共通room_idを探索）
	userRooms, err := d.roomStore.GetUserRooms(createdBy)
	if err != nil {
		return nil, err
	}
	partnerRooms, err := d.roomStore.GetUserRooms(partnerUserName)
	if err != nil {
		return nil, err
	}

	roomIndex := make(map[string]types.Room, len(userRooms))
	for _, r := range userRooms {
		roomIndex[r.RoomID] = r
	}
	for _, pr := range partnerRooms {
		if r, ok := roomIndex[pr.RoomID]; ok {
			// 共通の部屋が既に存在するため、それを返す
			existing := r
			return &existing, nil
		}
	}

	// 共通部屋が無い場合、新規作成（相手の名前を部屋名として使用）
	room, err := d.roomStore.CreateRoom(partnerUserName, createdBy)
	if err != nil {
		return nil, err
	}
	if err := d.roomStore.AddMember(room.RoomID, partnerUserName, "member"); err != nil {
		return nil, err
	}
	room.MemberCount = 2
	return room, nil
}

func (d *RoomDomain) GetRoom(roomID string) (*types.GetRoomResponse, error) {
	if roomID == "" {
		return nil, errors.New("room_id is required")
	}

	room, err := d.roomStore.GetRoom(roomID)
	if err != nil {
		return nil, err
	}

	members, err := d.roomStore.GetRoomMembers(roomID)
	if err != nil {
		return nil, err
	}

	return &types.GetRoomResponse{
		Room:    *room,
		Members: members,
	}, nil
}

func (d *RoomDomain) GetUserRooms(userID string) ([]types.Room, error) {
	if userID == "" {
		return nil, errors.New("user_id is required")
	}
	rooms, err := d.roomStore.GetUserRooms(userID)
	if err != nil {
		return nil, err
	}

	// 各部屋のメンバーを確認し、DM（2人部屋）の場合は相手のユーザー名を部屋名として返す
	for i := range rooms {
		members, err := d.roomStore.GetRoomMembers(rooms[i].RoomID)
		if err != nil {
			continue
		}
		if len(members) == 2 {
			// 自分以外のユーザー名を部屋名に設定
			other := members[0]
			if other == userID && len(members) > 1 {
				other = members[1]
			}
			rooms[i].RoomName = other
		}
	}
	return rooms, nil
}
