package usecase

import (
	"context"
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain"
	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/test/mock"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestShouldCreateANewRoomWhenInputsAreValid(t *testing.T) {
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

	usecase := NewCreateRoomUseCase(gateway)
	output, err := usecase.Execute(ctx, &input)
	assert.NotNil(t, output)
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput.RoomId, output.RoomId)
}

func TestShouldReturnAValidationErrorWhenAdminIdIsInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	input := CreateRoomUseCaseInput{
		AdminId:  "",
		Name:     "A Room",
		Category: "Game",
	}

	usecase := NewCreateRoomUseCase(mock.NewMockRoomGatewayInterface(ctrl))
	output, err := usecase.Execute(ctx, &input)
	assert.Nil(t, output)
	assert.IsType(t, &domain.ValidationError{}, err)
	assert.EqualError(t, err, valueobject.ErrRequiredId)
}

func TestShouldReturnAValidationErrorWhenNameIsInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	input := CreateRoomUseCaseInput{
		AdminId:  "auth0|64c8457bb160e37c8c34533b",
		Name:     "",
		Category: "Game",
	}

	usecase := NewCreateRoomUseCase(mock.NewMockRoomGatewayInterface(ctrl))
	output, err := usecase.Execute(ctx, &input)
	assert.Nil(t, output)
	assert.IsType(t, &domain.ValidationError{}, err)
	assert.EqualError(t, err, valueobject.ErrRequiredRoomName)
}

func TestShouldReturnAValidationErrorWhenCategoryIsInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	input := CreateRoomUseCaseInput{
		AdminId:  "auth0|64c8457bb160e37c8c34533b",
		Name:     "A Room",
		Category: "",
	}

	usecase := NewCreateRoomUseCase(mock.NewMockRoomGatewayInterface(ctrl))
	output, err := usecase.Execute(ctx, &input)
	assert.Nil(t, output)
	assert.IsType(t, &domain.ValidationError{}, err)
	assert.EqualError(t, err, valueobject.ErrRequiredRoomCategory)
}

func TestShouldReturnAGatewayErrorWhenGatewayReturnAnError(t *testing.T) {
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
		Return(domain.NewGatewayError("gateway error")).
		Times(1)

	usecase := NewCreateRoomUseCase(gateway)
	output, err := usecase.Execute(ctx, &input)
	assert.Nil(t, output)
	assert.NotNil(t, err)
	assert.IsType(t, &domain.GatewayError{}, err)
	assert.EqualError(t, err, "gateway error")
}