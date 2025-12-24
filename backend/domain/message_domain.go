package domain

import (
	"errors"

	"github.com/Kai-Kota/ChatApp2/backend/store"
	"github.com/Kai-Kota/ChatApp2/backend/types"
)

type IMessageDomain interface {
	GetMessages(roomID string, limit int) ([]types.Message, error)
}

type MessageDomain struct {
	store store.IMessageStore
}

func NewMessageDomain(store store.IMessageStore) IMessageDomain {
	return &MessageDomain{store: store}
}

func (d *MessageDomain) GetMessages(roomID string, limit int) ([]types.Message, error) {
	if roomID == "" {
		return nil, errors.New("room_id is required")
	}
	if limit <= 0 {
		limit = 20
	}
	return d.store.GetMessages(roomID, limit)
}
