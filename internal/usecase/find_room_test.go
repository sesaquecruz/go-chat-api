package usecase

import (
	"context"
	"database/sql"
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/errors"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/test/mock"
	"github.com/stretchr/testify/assert"

	"go.uber.org/mock/gomock"
)

func TestShouldReturnARoomWhenRoomIdExists(t *testing.T) {
	adminId, _ := valueobject.NewAuth0IDWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Need for Speed")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	room, _ := entity.NewRoom(adminId, name, category)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	input := FindRoomUseCaseInput{RoomId: room.Id().Value()}

	gateway := mock.NewMockRoomGatewayInterface(ctrl)
	gateway.
		EXPECT().
		FindById(gomock.Any(), gomock.Any()).
		Do(func(c context.Context, id *valueobject.ID) {
			assert.Equal(t, ctx, c)
			assert.Equal(t, room.Id().Value(), id.Value())
		}).
		Return(room, nil).
		Times(1)

	useCase := NewFindRoomUseCase(gateway)
	output, err := useCase.Execute(ctx, &input)
	assert.NotNil(t, output)
	assert.Nil(t, err)
	assert.Equal(t, room.Id().Value(), output.Id)
	assert.Equal(t, room.AdminId().Value(), output.AdminId)
	assert.Equal(t, room.Name().Value(), output.Name)
	assert.Equal(t, room.Category().Value(), output.Category)
	assert.Equal(t, room.CreatedAt().StringValue(), output.CreatedAt)
	assert.Equal(t, room.UpdatedAt().StringValue(), output.UpdatedAt)
}

func TestShouldReturnAnErrorWhenRoomIdDoesNotExist(t *testing.T) {
	id := valueobject.NewID().Value()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	input := FindRoomUseCaseInput{RoomId: id}

	gateway := mock.NewMockRoomGatewayInterface(ctrl)
	gateway.
		EXPECT().
		FindById(gomock.Any(), gomock.Any()).
		Do(func(c context.Context, i *valueobject.ID) {
			assert.Equal(t, ctx, c)
			assert.Equal(t, id, i.Value())
		}).
		Return(nil, sql.ErrNoRows).
		Times(1)

	useCase := NewFindRoomUseCase(gateway)
	output, err := useCase.Execute(ctx, &input)
	assert.Nil(t, output)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestShouldReturnAnErrorWhenRoomIdIsInvalid(t *testing.T) {
	id := "fdaiuo13j23ufoesfsfd"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	input := FindRoomUseCaseInput{RoomId: id}

	gateway := mock.NewMockRoomGatewayInterface(ctrl)
	gateway.
		EXPECT().
		FindById(gomock.Any(), gomock.Any()).
		Times(0)

	useCase := NewFindRoomUseCase(gateway)
	output, err := useCase.Execute(ctx, &input)
	assert.Nil(t, output)
	assert.NotNil(t, err)
	assert.IsType(t, &errors.ValidationError{}, err)
	assert.EqualError(t, err, valueobject.ErrInvalidId)
}
