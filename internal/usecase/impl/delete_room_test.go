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

func TestDeleteRoomUseCase_ShouldDeleteARoomWhenDataIsValid(t *testing.T) {
	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("A Game")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	savedRoom := entity.NewRoom(adminId, name, category)

	ctx := context.Background()
	input := &usecase.DeleteRoomUseCaseInput{
		Id:      savedRoom.Id().Value(),
		AdminId: savedRoom.AdminId().Value(),
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

	roomRepository.
		EXPECT().
		Update(mock.Anything, mock.Anything).
		Run(func(c context.Context, r *entity.Room) {
			assert.Equal(t, ctx, c)
			assert.Equal(t, input.Id, r.Id().Value())
		}).
		Return(nil).
		Once()

	useCase := NewDeleteRoomUseCase(roomRepository)

	err := useCase.Execute(ctx, input)
	assert.Nil(t, err)
}

func TestDeleteRoomUseCase_ShouldReturnAnErrorWhenDataIsInvalid(t *testing.T) {
	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("A Game")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	savedRoom := entity.NewRoom(adminId, name, category)

	ctx := context.Background()

	testCases := []struct {
		test  string
		input *usecase.DeleteRoomUseCaseInput
		err   error
	}{
		{
			"empty id",
			&usecase.DeleteRoomUseCaseInput{
				Id:      "",
				AdminId: "auth0|64c8457bb160e37c8c34533b",
			},
			valueobject.ErrRequiredId,
		},
		{
			"empty admin id",
			&usecase.DeleteRoomUseCaseInput{
				Id:      "b3588483-4795-434a-877c-dcd158d6caa7",
				AdminId: "",
			},
			valueobject.ErrRequiredUserId,
		},
		{
			"invalid id",
			&usecase.DeleteRoomUseCaseInput{
				Id:      "dfaioewurqredfa",
				AdminId: "auth0|64c8457bb160e37c8c34533b",
			},
			valueobject.ErrInvalidId,
		},
		{
			"invalid admin id",
			&usecase.DeleteRoomUseCaseInput{
				Id:      "b3588483-4795-434a-877c-dcd158d6caa7",
				AdminId: "fdafiuero3c8c34533b",
			},
			valueobject.ErrInvalidUserId,
		},
		{
			"invalid room admin",
			&usecase.DeleteRoomUseCaseInput{
				Id:      savedRoom.Id().Value(),
				AdminId: "auth0|64c8457bb160e37c8c34533c",
			},
			entity.ErrInvalidRoomAdmin,
		},
	}

	roomRepository := mocks.NewRoomRepositoryMock(t)

	roomRepository.
		EXPECT().
		FindById(mock.Anything, mock.Anything).
		Run(func(c context.Context, i *valueobject.Id) {
			assert.Equal(t, ctx, c)
			assert.Equal(t, savedRoom.Id().Value(), i.Value())
		}).
		Return(savedRoom, nil).
		Once()

	useCase := NewDeleteRoomUseCase(roomRepository)

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := useCase.Execute(ctx, tc.input)
			assert.NotNil(t, err)
			assert.ErrorIs(t, err, tc.err)
		})
	}
}

func TestDeleteRoomUseCase_ShouldReturnAnErrorOnRepositoryError(t *testing.T) {
	ctx := context.Background()
	input := &usecase.DeleteRoomUseCaseInput{
		Id:      "b3588483-4795-434a-877c-dcd158d6caa7",
		AdminId: "auth0|64c8457bb160e37c8c34533b",
	}

	roomRepository := mocks.NewRoomRepositoryMock(t)

	roomRepository.
		EXPECT().
		FindById(mock.Anything, mock.Anything).
		Return(nil, repository.ErrNotFoundRoom).
		Once()

	useCase := NewDeleteRoomUseCase(roomRepository)

	err := useCase.Execute(ctx, input)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, repository.ErrNotFoundRoom)
}
