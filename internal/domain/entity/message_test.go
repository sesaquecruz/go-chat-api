package entity

import (
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"

	"github.com/stretchr/testify/assert"
)

func TestMessage_ShouldCreateAMessageWhenDataIsValid(t *testing.T) {
	id := valueobject.NewId()
	roomId := valueobject.NewId()
	senderId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	senderName, _ := valueobject.NewUserNameWith("a username")
	text, _ := valueobject.NewMessageTextWith("a simple message")
	createdAt := valueobject.NewTimestamp()

	message := NewMessage(roomId, senderId, senderName, text)
	assert.NotNil(t, message.Id())
	assert.Equal(t, roomId.Value(), message.RoomId().Value())
	assert.Equal(t, senderId.Value(), message.SenderId().Value())
	assert.Equal(t, senderName.Value(), message.SenderName().Value())
	assert.Equal(t, text.Value(), message.Text().Value())
	assert.NotNil(t, message.CreatedAt())

	message = NewMessageWith(id, roomId, senderId, senderName, text, createdAt)
	assert.Equal(t, id.Value(), message.Id().Value())
	assert.Equal(t, roomId.Value(), message.RoomId().Value())
	assert.Equal(t, senderId.Value(), message.SenderId().Value())
	assert.Equal(t, senderName.Value(), message.SenderName().Value())
	assert.Equal(t, text.Value(), message.Text().Value())
	assert.Equal(t, createdAt.Value(), message.CreatedAt().Value())
}
