package usecase

import (
	"context"
	"database/sql"
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/errors"
	gateway_pkg "github.com/sesaquecruz/go-chat-api/internal/domain/gateway"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/test/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUpdateRoomUseCase_ShouldUpdateRoomWhenRoomInputIsValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := valueobject.NewID()
	adminId, _ := valueobject.NewAuth0IDWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Need for Speed")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	createdAt := valueobject.NewDateTime()
	updatedAt, _ := valueobject.NewDateTimeWith(createdAt.StringValue())

	room, _ := entity.NewRoomWith(
		id,
		adminId,
		name,
		category,
		createdAt,
		updatedAt,
	)

	ctx := context.Background()
	input := UpdateRoomUseCaseInput{
		Id:       id.Value(),
		AdminId:  adminId.Value(),
		Name:     "Rust",
		Category: "Tech",
	}

	gateway := mock.NewMockRoomGatewayInterface(ctrl)
	gateway.
		EXPECT().
		FindById(gomock.Any(), gomock.Any()).
		Do(func(c context.Context, i *valueobject.ID) {
			assert.Equal(t, ctx, c)
			assert.Equal(t, id.Value(), i.Value())
		}).
		Return(room, nil).
		Times(1)

	gateway.
		EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Do(func(c context.Context, r *entity.Room) {
			assert.Equal(t, ctx, c)
			assert.Equal(t, id.Value(), r.Id().Value())
			assert.Equal(t, adminId.Value(), r.AdminId().Value())
			assert.Equal(t, input.Name, r.Name().Value())
			assert.Equal(t, input.Category, r.Category().Value())
			assert.True(t, createdAt.TimeValue().Equal(*r.CreatedAt().TimeValue()))
			assert.True(t, updatedAt.TimeValue().Before(*r.UpdatedAt().TimeValue()))
		}).
		Return(nil).
		Times(1)

	useCase := NewUpdateRoomUseCase(gateway)
	err := useCase.Execute(ctx, &input)
	assert.Nil(t, err)
}

func TestUpdateRoomUseCase_ShouldReturnAGatewayErrorWhenRoomIdDoesNotExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	input := UpdateRoomUseCaseInput{
		Id:       valueobject.NewID().Value(),
		AdminId:  "auth0|64c8457bb160e37c8c34533b",
		Name:     "Rust",
		Category: "Tech",
	}

	gateway := mock.NewMockRoomGatewayInterface(ctrl)
	gateway.
		EXPECT().
		FindById(gomock.Any(), gomock.Any()).
		Return(nil, sql.ErrNoRows).
		Times(1)

	useCase := NewUpdateRoomUseCase(gateway)
	err := useCase.Execute(ctx, &input)
	assert.NotNil(t, err)
	assert.IsType(t, &errors.GatewayError{}, err)
	assert.EqualError(t, err, gateway_pkg.ErrNotFoundRoom)
}

func TestUpdateRoomUseCase_ShouldReturnAValidationErrorWhenAdminIdIsNotRoomAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := valueobject.NewID()
	adminId, _ := valueobject.NewAuth0IDWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Need for Speed")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	createdAt := valueobject.NewDateTime()
	updatedAt, _ := valueobject.NewDateTimeWith(createdAt.StringValue())

	room, _ := entity.NewRoomWith(
		id,
		adminId,
		name,
		category,
		createdAt,
		updatedAt,
	)

	ctx := context.Background()
	input := UpdateRoomUseCaseInput{
		Id:       id.Value(),
		AdminId:  "auth0|64c8457bb160e37c8c34533c",
		Name:     "Rust",
		Category: "Tech",
	}

	gateway := mock.NewMockRoomGatewayInterface(ctrl)
	gateway.
		EXPECT().
		FindById(gomock.Any(), gomock.Any()).
		Do(func(c context.Context, i *valueobject.ID) {
			assert.Equal(t, ctx, c)
			assert.Equal(t, id.Value(), i.Value())
		}).
		Return(room, nil).
		Times(1)

	useCase := NewUpdateRoomUseCase(gateway)
	err := useCase.Execute(ctx, &input)
	assert.NotNil(t, err)
	assert.IsType(t, &errors.ValidationError{}, err)
	assert.EqualError(t, err, valueobject.ErrInvalidId)
}

func TestUpdateRoomUseCase_ShouldReturnAValidationErrorWhenInputIsInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		input *UpdateRoomUseCaseInput
		err   error
	}{
		{
			input: &UpdateRoomUseCaseInput{
				Id:       "fdfagferr",
				AdminId:  "",
				Name:     "",
				Category: "",
			},
			err: errors.NewValidationError(valueobject.ErrInvalidId),
		},
		{
			input: &UpdateRoomUseCaseInput{
				Id:       valueobject.NewID().Value(),
				AdminId:  "fdadaervc233",
				Name:     "",
				Category: "",
			},
			err: errors.NewValidationError(valueobject.ErrInvalidId),
		},
		{
			input: &UpdateRoomUseCaseInput{
				Id:       valueobject.NewID().Value(),
				AdminId:  "auth0|64c8457bb160e37c8c34533c",
				Name:     "",
				Category: "",
			},
			err: errors.NewValidationError(valueobject.ErrRequiredRoomName),
		},
		{
			input: &UpdateRoomUseCaseInput{
				Id:       valueobject.NewID().Value(),
				AdminId:  "auth0|64c8457bb160e37c8c34533c",
				Name:     "Rust",
				Category: "Anything",
			},
			err: errors.NewValidationError(valueobject.ErrInvalidRoomCategory),
		},
	}

	ctx := context.Background()
	gateway := mock.NewMockRoomGatewayInterface(ctrl)
	useCase := NewUpdateRoomUseCase(gateway)

	for _, test := range testCases {
		err := useCase.Execute(ctx, test.input)
		assert.IsType(t, &errors.ValidationError{}, err)
		assert.EqualError(t, err, test.err.Error())
	}
}
