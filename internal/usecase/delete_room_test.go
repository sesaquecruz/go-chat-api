package usecase

import (
	"context"
	"database/sql"
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/test/mock"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDeleteRoom_ShouldDeleteARoomWhenInputIsValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Need for Speed")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	room, _ := entity.NewRoom(adminId, name, category)

	ctx := context.Background()
	input := DeleteRoomUseCaseInput{
		Id:      room.Id().Value(),
		AdminId: room.AdminId().Value(),
	}

	repository := mock.NewMockRoomRepositoryInterface(ctrl)

	repository.
		EXPECT().
		FindById(gomock.Any(), gomock.Any()).
		Do(func(c context.Context, id *valueobject.Id) {
			assert.Equal(t, ctx, c)
			assert.Equal(t, room.Id().Value(), id.Value())
		}).
		Return(room, nil).
		Times(1)

	repository.
		EXPECT().
		Delete(gomock.Any(), gomock.Any()).
		Do(func(c context.Context, i *valueobject.Id) {
			assert.Equal(t, ctx, c)
			assert.Equal(t, room.Id().Value(), i.Value())
		}).
		Return(nil).
		Times(1)

	useCase := NewDeleteRoomUseCase(repository)
	err := useCase.Execute(ctx, &input)
	assert.Nil(t, err)
}

func TestDeleteRoom_ShouldReturnAnIdErrorWhenInputIsInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		input *DeleteRoomUseCaseInput
		err   error
	}{
		{
			input: &DeleteRoomUseCaseInput{
				Id:      "dfaioewurqredfa",
				AdminId: "auth0|64c8457bb160e37c8c34533b",
			},
			err: validation.ErrInvalidId,
		},
		{
			input: &DeleteRoomUseCaseInput{
				Id:      valueobject.NewId().Value(),
				AdminId: "fdafiuero3c8c34533b",
			},
			err: validation.ErrInvalidId,
		},
	}

	ctx := context.Background()
	repository := mock.NewMockRoomRepositoryInterface(ctrl)
	useCase := NewDeleteRoomUseCase(repository)

	for _, test := range testCases {
		err := useCase.Execute(ctx, test.input)
		assert.NotNil(t, err)
		assert.ErrorIs(t, test.err, validation.ErrInvalidId)
	}
}

func TestDeleteRoom_ShouldReturnAnNotFoundErrorWhenRoomDoesNotExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	input := DeleteRoomUseCaseInput{
		Id:      valueobject.NewId().Value(),
		AdminId: "auth0|64c8457bb160e37c8c34533b",
	}

	repository := mock.NewMockRoomRepositoryInterface(ctrl)

	repository.
		EXPECT().
		FindById(gomock.Any(), gomock.Any()).
		Return(nil, sql.ErrNoRows).
		Times(1)

	useCase := NewDeleteRoomUseCase(repository)
	err := useCase.Execute(ctx, &input)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, validation.ErrNotFoundRoom)
}

func TestDeleteRoom_ShouldReturnAnAuthorizationErrorWhenAdminIdIsInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Need for Speed")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	room, _ := entity.NewRoom(adminId, name, category)

	ctx := context.Background()
	input := DeleteRoomUseCaseInput{
		Id:      room.Id().Value(),
		AdminId: "auth0|64c8457bb160e37c8c34533c",
	}

	repository := mock.NewMockRoomRepositoryInterface(ctrl)

	repository.
		EXPECT().
		FindById(gomock.Any(), gomock.Any()).
		Do(func(c context.Context, id *valueobject.Id) {
			assert.Equal(t, ctx, c)
			assert.Equal(t, room.Id().Value(), id.Value())
		}).
		Return(room, nil).
		Times(1)

	useCase := NewDeleteRoomUseCase(repository)
	err := useCase.Execute(ctx, &input)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, validation.ErrInvalidRoomAdmin)
}
