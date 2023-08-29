package entity

import (
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
)

type Message struct {
	id         *valueobject.Id
	roomId     *valueobject.Id
	senderId   *valueobject.UserId
	senderName *valueobject.UserName
	text       *valueobject.MessageText
	createdAt  *valueobject.Timestamp
}

func NewMessage(
	roomId *valueobject.Id,
	senderId *valueobject.UserId,
	senderName *valueobject.UserName,
	text *valueobject.MessageText,
) *Message {
	return NewMessageWith(
		valueobject.NewId(),
		roomId,
		senderId,
		senderName,
		text,
		valueobject.NewTimestamp(),
	)
}

func NewMessageWith(
	id *valueobject.Id,
	roomId *valueobject.Id,
	senderId *valueobject.UserId,
	senderName *valueobject.UserName,
	text *valueobject.MessageText,
	createdAt *valueobject.Timestamp,
) *Message {
	return &Message{
		id:         id,
		roomId:     roomId,
		senderId:   senderId,
		senderName: senderName,
		text:       text,
		createdAt:  createdAt,
	}
}

func (m *Message) Id() *valueobject.Id {
	return m.id
}

func (m *Message) RoomId() *valueobject.Id {
	return m.roomId
}

func (m *Message) SenderId() *valueobject.UserId {
	return m.senderId
}

func (m *Message) SenderName() *valueobject.UserName {
	return m.senderName
}

func (m *Message) Text() *valueobject.MessageText {
	return m.text
}

func (m *Message) CreatedAt() *valueobject.Timestamp {
	return m.createdAt
}
