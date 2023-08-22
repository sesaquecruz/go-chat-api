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

func TestFindRoomUseCase_ShouldReturnARoomWhenIdExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Need for Speed")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	room, _ := entity.NewRoom(adminId, name, category)

	ctx := context.Background()
	input := FindRoomUseCaseInput{RoomId: room.Id().Value()}

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

	useCase := NewFindRoomUseCase(repository)
	output, err := useCase.Execute(ctx, &input)
	assert.NotNil(t, output)
	assert.Nil(t, err)
	assert.Equal(t, room.Id().Value(), output.Id)
	assert.Equal(t, room.AdminId().Value(), output.AdminId)
	assert.Equal(t, room.Name().Value(), output.Name)
	assert.Equal(t, room.Category().Value(), output.Category)
	assert.Equal(t, room.CreatedAt().String(), output.CreatedAt)
	assert.Equal(t, room.UpdatedAt().String(), output.UpdatedAt)
}

func TestFindRoomUseCase_ShouldReturnANotFoundErrorWhenIdDoesNotExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := valueobject.NewId().Value()

	ctx := context.Background()
	input := FindRoomUseCaseInput{RoomId: id}

	repository := mock.NewMockRoomRepositoryInterface(ctrl)
	repository.
		EXPECT().
		FindById(gomock.Any(), gomock.Any()).
		Do(func(c context.Context, i *valueobject.Id) {
			assert.Equal(t, ctx, c)
			assert.Equal(t, id, i.Value())
		}).
		Return(nil, sql.ErrNoRows).
		Times(1)

	useCase := NewFindRoomUseCase(repository)
	output, err := useCase.Execute(ctx, &input)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, validation.ErrNotFoundRoom)
}

func TestFindRoomUseCase_ShouldReturnAnIdErrorWhenIdIsInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := "fdaiuo13j23ufoesfsfd"

	ctx := context.Background()
	input := FindRoomUseCaseInput{RoomId: id}

	repository := mock.NewMockRoomRepositoryInterface(ctrl)

	useCase := NewFindRoomUseCase(repository)
	output, err := useCase.Execute(ctx, &input)
	assert.Nil(t, output)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, validation.ErrInvalidId)
}
