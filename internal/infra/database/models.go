package database

import (
	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
)

type RoomModel struct {
	Id        string
	Admin_id  string
	Name      string
	Category  string
	CreatedAt string
	UpdatedAt string
}

func (m *RoomModel) ToEntity() (*entity.Room, error) {
	id, err := valueobject.NewIdWith(m.Id)
	if err != nil {
		return nil, err
	}

	adminId, err := valueobject.NewUserIdWith(m.Admin_id)
	if err != nil {
		return nil, err
	}

	name, err := valueobject.NewRoomNameWith(m.Name)
	if err != nil {
		return nil, err
	}

	category, err := valueobject.NewRoomCategoryWith(m.Category)
	if err != nil {
		return nil, err
	}

	created_at, err := valueobject.NewTimestampWith(m.CreatedAt)
	if err != nil {
		return nil, err
	}

	updated_at, err := valueobject.NewTimestampWith(m.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return entity.NewRoomWith(id, adminId, name, category, created_at, updated_at)
}

type MessageModel struct {
	Id         string
	RoomId     string
	SenderId   string
	SenderName string
	Text       string
	CreatedAt  string
}

func (m *MessageModel) ToEntity() (*entity.Message, error) {
	id, err := valueobject.NewIdWith(m.Id)
	if err != nil {
		return nil, err
	}

	roomId, err := valueobject.NewIdWith(m.RoomId)
	if err != nil {
		return nil, err
	}

	senderId, err := valueobject.NewUserIdWith(m.SenderId)
	if err != nil {
		return nil, err
	}

	senderName, err := valueobject.NewUserNameWith(m.SenderName)
	if err != nil {
		return nil, err
	}

	text, err := valueobject.NewMessageTextWith(m.Text)
	if err != nil {
		return nil, err
	}

	createdAt, err := valueobject.NewTimestampWith(m.CreatedAt)
	if err != nil {
		return nil, err
	}

	return entity.NewMessageWith(id, roomId, senderId, senderName, text, createdAt)
}
