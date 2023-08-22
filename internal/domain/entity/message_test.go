package entity

import (
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestMessage_ShouldCreateAMessageWhenDataAreValid(t *testing.T) {
	id := valueobject.NewId()
	roomId := valueobject.NewId()
	senderId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	senderName, _ := valueobject.NewUserNameWith("username")
	text, _ := valueobject.NewMessageTextWith("a simple message")
	createdAt := valueobject.NewTimestamp()

	message, err := NewMessage(roomId, senderId, senderName, text)
	assert.NotNil(t, message)
	assert.Nil(t, err)
	assert.NotNil(t, message.Id())
	assert.Equal(t, roomId.Value(), message.RoomId().Value())
	assert.Equal(t, senderId.Value(), message.SenderId().Value())
	assert.Equal(t, senderName.Value(), message.SenderName().Value())
	assert.Equal(t, text.Value(), message.Text().Value())
	assert.NotNil(t, message.CreatedAt())

	message, err = NewMessageWith(id, roomId, senderId, senderName, text, createdAt)
	assert.NotNil(t, message)
	assert.Nil(t, err)
	assert.Equal(t, id.Value(), message.Id().Value())
	assert.Equal(t, roomId.Value(), message.RoomId().Value())
	assert.Equal(t, senderId.Value(), message.SenderId().Value())
	assert.Equal(t, senderName.Value(), message.SenderName().Value())
	assert.Equal(t, text.Value(), message.Text().Value())
	assert.Equal(t, createdAt.Value(), message.CreatedAt().Value())
}

func TestMessage_ShouldReturnAnErrorWhenCreateAMessageWithInvalidData(t *testing.T) {
	id := valueobject.NewId()
	roomId := valueobject.NewId()
	senderId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	senderName, _ := valueobject.NewUserNameWith("username")
	text, _ := valueobject.NewMessageTextWith("a simple message")
	createdAt := valueobject.NewTimestamp()

	testCases := []struct {
		testName   string
		errNew     error
		errWith    error
		id         *valueobject.Id
		roomId     *valueobject.Id
		senderId   *valueobject.UserId
		senderName *valueobject.UserName
		text       *valueobject.MessageText
		createdAt  *valueobject.Timestamp
	}{
		{
			"nil id",
			nil,
			validation.ErrRequiredMessageId,
			nil,
			roomId,
			senderId,
			senderName,
			text,
			createdAt,
		},
		{
			"nil room id",
			validation.ErrRequiredMessageRoomId,
			validation.ErrRequiredMessageRoomId,
			id,
			nil,
			senderId,
			senderName,
			text,
			createdAt,
		},
		{
			"nil sender id",
			validation.ErrRequiredMessageSenderId,
			validation.ErrRequiredMessageSenderId,
			id,
			roomId,
			nil,
			senderName,
			text,
			createdAt,
		},
		{
			"nil sender name",
			validation.ErrRequiredMessageSenderName,
			validation.ErrRequiredMessageSenderName,
			id,
			roomId,
			senderId,
			nil,
			text,
			createdAt,
		},
		{
			"nil text",
			validation.ErrRequiredMessageText,
			validation.ErrRequiredMessageText,
			id,
			roomId,
			senderId,
			senderName,
			nil,
			createdAt,
		},
		{
			"nil created at",
			nil,
			validation.ErrRequiredMessageCreatedAt,
			id,
			roomId,
			senderId,
			senderName,
			text,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			_, err := NewMessage(
				test.roomId,
				test.senderId,
				test.senderName,
				test.text,
			)
			assert.Equal(t, test.errNew, err)

			_, err = NewMessageWith(
				test.id,
				test.roomId,
				test.senderId,
				test.senderName,
				test.text,
				test.createdAt,
			)
			assert.Equal(t, test.errWith, err)
		})
	}
}
