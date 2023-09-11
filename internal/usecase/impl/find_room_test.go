package impl

import (
	"context"
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/repository"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFindRoomUseCase_ShouldReturnARoomWhenDataIsValid(t *testing.T) {
	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("A Game")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	savedRoom := entity.NewRoom(adminId, name, category)

	ctx := context.Background()
	input := &usecase.FindRoomUseCaseInput{
		Id: savedRoom.Id().Value(),
	}

	roomRepository := mocks.NewRoomRepositoryMock(t)

	roomRepository.
		EXPECT().
		FindById(mock.Anything, mock.Anything).
		Run(func(c context.Context, i *valueobject.Id) {
			assert.Equal(t, ctx, c)
			assert.Equal(t, input.Id, i.Value())
		}).
		Return(savedRoom, nil).
		Once()

	useCase := NewFindRoomUseCase(roomRepository)
	output, err := useCase.Execute(ctx, input)
	assert.NotNil(t, output)
	assert.Nil(t, err)
	assert.Equal(t, savedRoom.Id().Value(), output.Id)
	assert.Equal(t, savedRoom.AdminId().Value(), output.AdminId)
	assert.Equal(t, savedRoom.Name().Value(), output.Name)
	assert.Equal(t, savedRoom.Category().Value(), output.Category)
	assert.Equal(t, savedRoom.CreatedAt().Value(), output.CreatedAt)
	assert.Equal(t, savedRoom.UpdatedAt().Value(), output.UpdatedAt)
}

func TestFindRoomUseCase_ShouldReturnAnErrorWhenDataIsInvalid(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		test  string
		input *usecase.FindRoomUseCaseInput
		err   error
	}{
		{
			"empty id",
			&usecase.FindRoomUseCaseInput{
				Id: "",
			},
			valueobject.ErrRequiredId,
		},
		{
			"invalid id",
			&usecase.FindRoomUseCaseInput{
				Id: "dfaioewurqredfa",
			},
			valueobject.ErrInvalidId,
		},
	}

	roomRepository := mocks.NewRoomRepositoryMock(t)
	useCase := NewFindRoomUseCase(roomRepository)

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			output, err := useCase.Execute(ctx, tc.input)
			assert.Nil(t, output)
			assert.NotNil(t, err)
			assert.ErrorIs(t, err, tc.err)
		})
	}
}

func TestFindRoomUseCase_ShouldReturnAnErrorOnRepositoryError(t *testing.T) {
	ctx := context.Background()
	input := &usecase.FindRoomUseCaseInput{
		Id: "b3588483-4795-434a-877c-dcd158d6caa7",
	}

	roomRepository := mocks.NewRoomRepositoryMock(t)

	roomRepository.
		EXPECT().
		FindById(mock.Anything, mock.Anything).
		Return(nil, repository.ErrNotFoundRoom).
		Once()

	useCase := NewFindRoomUseCase(roomRepository)

	output, err := useCase.Execute(ctx, input)
	assert.Nil(t, output)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, repository.ErrNotFoundRoom)
}
