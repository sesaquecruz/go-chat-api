package impl

import (
	"context"
	"errors"
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateRoomUseCase_ShouldCreateARoomWhenDataIsValid(t *testing.T) {
	roomCreated := &entity.Room{}

	ctx := context.Background()
	input := &usecase.CreateRoomUseCaseInput{
		AdminId:  "auth0|64c8457bb160e37c8c34533b",
		Name:     "A Game",
		Category: "Game",
	}

	roomRepository := mocks.NewRoomRepositoryMock(t)

	roomRepository.
		EXPECT().
		Save(mock.Anything, mock.Anything).
		Run(func(c context.Context, r *entity.Room) {
			assert.Equal(t, ctx, c)
			assert.Equal(t, input.AdminId, r.AdminId().Value())
			assert.Equal(t, input.Name, r.Name().Value())
			assert.Equal(t, input.Category, r.Category().Value())
			roomCreated = r
		}).
		Return(nil).
		Once()

	useCase := NewCreateRoomUseCase(roomRepository)

	output, err := useCase.Execute(ctx, input)
	assert.NotNil(t, output)
	assert.Nil(t, err)
	assert.Equal(t, roomCreated.Id().Value(), output.RoomId)
}

func TestCreateRoomUseCase_ShouldReturnAnErrorWhenDataIsInvalid(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		test  string
		input *usecase.CreateRoomUseCaseInput
		err   error
	}{
		{
			"empty admin id",
			&usecase.CreateRoomUseCaseInput{
				AdminId:  "",
				Name:     "A Game",
				Category: "Game",
			},
			valueobject.ErrRequiredUserId,
		},
		{
			"empty name",
			&usecase.CreateRoomUseCaseInput{
				AdminId:  "auth0|64c8457bb160e37c8c34533b",
				Name:     "",
				Category: "Game",
			},
			valueobject.ErrRequiredRoomName,
		},
		{
			"empty category",
			&usecase.CreateRoomUseCaseInput{
				AdminId:  "auth0|64c8457bb160e37c8c34533b",
				Name:     "A Game",
				Category: "",
			},
			valueobject.ErrRequiredRoomCategory,
		},
	}

	roomRepository := mocks.NewRoomRepositoryMock(t)
	useCase := NewCreateRoomUseCase(roomRepository)

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			output, err := useCase.Execute(ctx, tc.input)
			assert.Nil(t, output)
			assert.NotNil(t, err)
			assert.ErrorIs(t, err, tc.err)
		})
	}
}

func TestCreateRoomUseCase_ShouldReturnAnErrorOnRepositoryError(t *testing.T) {
	ctx := context.Background()
	input := &usecase.CreateRoomUseCaseInput{
		AdminId:  "auth0|64c8457bb160e37c8c34533b",
		Name:     "A Game",
		Category: "Game",
	}

	roomRepository := mocks.NewRoomRepositoryMock(t)

	roomRepository.
		EXPECT().
		Save(mock.Anything, mock.Anything).
		Return(errors.New("a repository error")).
		Once()

	useCase := NewCreateRoomUseCase(roomRepository)

	output, err := useCase.Execute(ctx, input)
	assert.Nil(t, output)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "a repository error")
}
