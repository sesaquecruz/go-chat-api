package model

import (
	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
)

type MessageModel struct {
	Id         string
	RoomId     string
	SenderId   string
	SenderName string
	Text       string
	CreatedAt  string
}

func NewMessageModel(message *entity.Message) *MessageModel {
	model := MessageModel{}

	model.Id = message.Id().Value()
	model.RoomId = message.RoomId().Value()
	model.SenderId = message.SenderId().Value()
	model.SenderName = message.SenderName().Value()
	model.Text = message.Text().Value()
	model.CreatedAt = message.CreatedAt().Value()

	return &model
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

	message := entity.NewMessageWith(id, roomId, senderId, senderName, text, createdAt)

	return message, nil
}
