package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/test/mock"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateRoomUseCase_ShouldCreateANewRoomWhenInputIsValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	input := CreateRoomUseCaseInput{
		AdminId:  "auth0|64c8457bb160e37c8c34533b",
		Name:     "A Room",
		Category: "Game",
	}
	expectedOutput := &CreateRoomUseCaseOutput{}

	gateway := mock.NewMockRoomGatewayInterface(ctrl)
	gateway.
		EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Do(func(c context.Context, r *entity.Room) {
			assert.Equal(t, ctx, c)
			assert.Nil(t, r.Validate())
			assert.Equal(t, input.AdminId, r.AdminId().Value())
			assert.Equal(t, input.Name, r.Name().Value())
			assert.Equal(t, input.Category, r.Category().Value())
			expectedOutput.RoomId = r.Id().Value()
		}).
		Return(nil).
		Times(1)

	useCase := NewCreateRoomUseCase(gateway)
	output, err := useCase.Execute(ctx, &input)
	assert.NotNil(t, output)
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput.RoomId, output.RoomId)
}

func TestCreateRoomUseCase_ShouldReturnAnIdErrorWhenAdminIdIsInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	input := CreateRoomUseCaseInput{
		AdminId:  "",
		Name:     "A Room",
		Category: "Game",
	}

	gateway := mock.NewMockRoomGatewayInterface(ctrl)

	useCase := NewCreateRoomUseCase(gateway)
	output, err := useCase.Execute(ctx, &input)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, validation.ErrRequiredId)
}

func TestCreateRoomUseCase_ShouldReturnANameErrorWhenNameIsInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	input := CreateRoomUseCaseInput{
		AdminId:  "auth0|64c8457bb160e37c8c34533b",
		Name:     "",
		Category: "Game",
	}

	gateway := mock.NewMockRoomGatewayInterface(ctrl)

	useCase := NewCreateRoomUseCase(gateway)
	output, err := useCase.Execute(ctx, &input)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, validation.ErrRequiredRoomName)
}

func TestCreateRoomUseCase_ShouldReturnACategoryErrorWhenCategoryIsInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	input := CreateRoomUseCaseInput{
		AdminId:  "auth0|64c8457bb160e37c8c34533b",
		Name:     "A Room",
		Category: "",
	}

	gateway := mock.NewMockRoomGatewayInterface(ctrl)

	useCase := NewCreateRoomUseCase(gateway)
	output, err := useCase.Execute(ctx, &input)
	assert.Nil(t, output)
	assert.ErrorIs(t, err, validation.ErrRequiredRoomCategory)
}

func TestCreateRoomUseCase_ShouldReturnAnInternalErrorWhenOccursAGatewayError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	input := CreateRoomUseCaseInput{
		AdminId:  "auth0|64c8457bb160e37c8c34533b",
		Name:     "A Room",
		Category: "Game",
	}

	gateway := mock.NewMockRoomGatewayInterface(ctrl)
	gateway.
		EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(validation.NewInternalError(errors.New("a gateway error"))).
		Times(1)

	useCase := NewCreateRoomUseCase(gateway)
	output, err := useCase.Execute(ctx, &input)
	assert.Nil(t, output)
	assert.NotNil(t, err)
	assert.IsType(t, &validation.InternalError{}, err)
	assert.EqualError(t, err, "a gateway error")
}
