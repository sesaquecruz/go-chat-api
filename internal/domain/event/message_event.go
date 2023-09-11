package event

import (
	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
)

type MessageEvent struct {
	Id         string `json:"id"`
	RoomId     string `json:"room_id"`
	SenderId   string `json:"sender_id"`
	SenderName string `json:"sender_name"`
	Text       string `json:"text"`
	CreatedAt  string `json:"created_at"`
}

func NewMessageEvent(message *entity.Message) *MessageEvent {
	messageEvent := &MessageEvent{
		Id:         message.Id().Value(),
		RoomId:     message.RoomId().Value(),
		SenderId:   message.SenderId().Value(),
		SenderName: message.SenderName().Value(),
		Text:       message.Text().Value(),
		CreatedAt:  message.CreatedAt().Value(),
	}

	return messageEvent
}
