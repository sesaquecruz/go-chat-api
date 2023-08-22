package usecase

import (
	"context"
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/test/mock"
	"github.com/stretchr/testify/assert"

	"go.uber.org/mock/gomock"
)

func TestCreateMessageUseCase_ShouldCreateAMessageInputIsValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Need for Speed")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	room, _ := entity.NewRoom(adminId, name, category)

	ctx := context.Background()
	input := CreateMessageUseCaseInput{
		RoomId:     room.Id().Value(),
		SenderId:   "auth0|64c8457bb160e37c8c34533g",
		SenderName: "username",
		Text:       "A simple message",
	}

	var message *entity.Message

	roomRepository := mock.NewMockRoomRepositoryInterface(ctrl)
	messageRepository := mock.NewMockMessageRepositoryInterface(ctrl)
	messageGateway := mock.NewMockMessageGatewayInterface(ctrl)

	roomRepository.
		EXPECT().
		FindById(gomock.Any(), gomock.Any()).
		Do(func(c context.Context, id *valueobject.Id) {
			assert.Equal(t, ctx, c)
			assert.Equal(t, input.RoomId, id.Value())
		}).
		Return(room, nil).
		Times(1)

	messageRepository.
		EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Do(func(c context.Context, m *entity.Message) {
			assert.Equal(t, ctx, c)
			assert.Equal(t, input.RoomId, m.RoomId().Value())
			assert.Equal(t, input.SenderId, m.SenderId().Value())
			assert.Equal(t, input.SenderName, m.SenderName().Value())
			assert.Equal(t, input.Text, m.Text().Value())
			message = m
		}).
		Return(nil).
		Times(1)

	messageGateway.
		EXPECT().
		Send(gomock.Any(), gomock.Any()).
		Do(func(c context.Context, m *entity.Message) {
			assert.Equal(t, ctx, c)
			assert.Equal(t, message.Id().Value(), m.Id().Value())
			assert.Equal(t, message.RoomId().Value(), m.RoomId().Value())
			assert.Equal(t, message.SenderId().Value(), m.SenderId().Value())
			assert.Equal(t, message.SenderName().Value(), m.SenderName().Value())
			assert.Equal(t, message.Text().Value(), m.Text().Value())
			assert.Equal(t, message.CreatedAt().String(), m.CreatedAt().String())
		}).
		Return(nil).
		Times(1)

	useCase := NewCreateMessageUseCase(roomRepository, messageRepository, messageGateway)

	output, err := useCase.Execute(ctx, &input)
	assert.NotNil(t, output)
	assert.Nil(t, err)
	assert.Equal(t, message.Id().Value(), output.MessageId)
}
