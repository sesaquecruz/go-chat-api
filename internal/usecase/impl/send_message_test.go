package impl

import (
	"context"
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/event"
	"github.com/sesaquecruz/go-chat-api/internal/domain/repository"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSendMessageUseCase_ShouldCreateAMessageWhenDataIsValid(t *testing.T) {
	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("A Game")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	roomSaved := entity.NewRoom(adminId, name, category)

	messageCreated := &entity.Message{}

	ctx := context.Background()
	input := &usecase.SendMessageUseCaseInput{
		RoomId:     roomSaved.Id().Value(),
		SenderId:   roomSaved.AdminId().Value(),
		SenderName: "An username",
		Text:       "A text",
	}

	roomRepository := mocks.NewRoomRepositoryMock(t)
	messageRepository := mocks.NewMessageRepositoryMock(t)
	messageEventGateway := mocks.NewMessageEventGatewayMock(t)

	roomRepository.EXPECT().
		FindById(mock.Anything, mock.Anything).
		Run(func(c context.Context, i *valueobject.Id) {
			assert.Equal(t, ctx, c)
			assert.Equal(t, input.RoomId, i.Value())
		}).
		Return(roomSaved, nil).
		Once()

	messageRepository.
		EXPECT().
		Save(mock.Anything, mock.Anything).
		Run(func(c context.Context, m *entity.Message) {
			assert.Equal(t, ctx, c)
			assert.Equal(t, input.RoomId, m.RoomId().Value())
			assert.Equal(t, input.SenderId, m.SenderId().Value())
			assert.Equal(t, input.SenderName, m.SenderName().Value())
			assert.Equal(t, input.Text, m.Text().Value())
			messageCreated = m
		}).
		Return(nil).
		Once()

	messageEventGateway.
		EXPECT().
		Send(mock.Anything, mock.Anything).
		Run(func(c context.Context, e *event.MessageEvent) {
			assert.Equal(t, ctx, c)
			assert.Equal(t, messageCreated.Id().Value(), e.Id)
			assert.Equal(t, messageCreated.RoomId().Value(), e.RoomId)
			assert.Equal(t, messageCreated.SenderId().Value(), e.SenderId)
			assert.Equal(t, messageCreated.SenderName().Value(), e.SenderName)
			assert.Equal(t, messageCreated.Text().Value(), e.Text)
			assert.Equal(t, messageCreated.CreatedAt().Value(), e.CreatedAt)
		}).
		Return(nil).
		Once()

	useCase := NewSendMessageUseCase(roomRepository, messageRepository, messageEventGateway)

	output, err := useCase.Execute(ctx, input)
	assert.NotNil(t, output)
	assert.Nil(t, err)
	assert.Equal(t, messageCreated.Id().Value(), output.MessageId)
}

func TestSendMessageUseCase_ShouldReturnAnErrorWhenDataIsInvalid(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		test  string
		input *usecase.SendMessageUseCaseInput
		err   error
	}{
		{
			"empty room id",
			&usecase.SendMessageUseCaseInput{
				RoomId:     "",
				SenderId:   "auth0|64c8457bb160e37c8c34533b",
				SenderName: "An username",
				Text:       "A text",
			},
			valueobject.ErrRequiredId,
		},
		{
			"empty sender id",
			&usecase.SendMessageUseCaseInput{
				RoomId:     "b3588483-4795-434a-877c-dcd158d6caa7",
				SenderId:   "",
				SenderName: "An username",
				Text:       "A text",
			},
			valueobject.ErrRequiredUserId,
		},
		{
			"empty sender name",
			&usecase.SendMessageUseCaseInput{
				RoomId:     "b3588483-4795-434a-877c-dcd158d6caa7",
				SenderId:   "auth0|64c8457bb160e37c8c34533b",
				SenderName: "",
				Text:       "A text",
			},
			valueobject.ErrRequiredUserName,
		},
		{
			"empty text",
			&usecase.SendMessageUseCaseInput{
				RoomId:     "b3588483-4795-434a-877c-dcd158d6caa7",
				SenderId:   "auth0|64c8457bb160e37c8c34533b",
				SenderName: "An username",
				Text:       "",
			},
			valueobject.ErrRequiredMessageText,
		},
	}

	roomRepository := mocks.NewRoomRepositoryMock(t)
	messageRepository := mocks.NewMessageRepositoryMock(t)
	messageEventGateway := mocks.NewMessageEventGatewayMock(t)

	useCase := NewSendMessageUseCase(roomRepository, messageRepository, messageEventGateway)

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			output, err := useCase.Execute(ctx, tc.input)
			assert.Nil(t, output)
			assert.NotNil(t, err)
			assert.ErrorIs(t, err, tc.err)
		})
	}
}

func TestSendMessageUseCase_ShouldReturnAnErrorOnRepositoryError(t *testing.T) {
	ctx := context.Background()
	input := &usecase.SendMessageUseCaseInput{
		RoomId:     "b3588483-4795-434a-877c-dcd158d6caa7",
		SenderId:   "auth0|64c8457bb160e37c8c34533b",
		SenderName: "An username",
		Text:       "A text",
	}

	roomRepository := mocks.NewRoomRepositoryMock(t)
	messageRepository := mocks.NewMessageRepositoryMock(t)
	messageEventGateway := mocks.NewMessageEventGatewayMock(t)

	roomRepository.EXPECT().
		FindById(mock.Anything, mock.Anything).
		Return(nil, repository.ErrNotFoundMessage).
		Once()

	useCase := NewSendMessageUseCase(roomRepository, messageRepository, messageEventGateway)

	output, err := useCase.Execute(ctx, input)
	assert.Nil(t, output)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, repository.ErrNotFoundMessage)
}
