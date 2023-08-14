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

func TestUpdateRoomUseCase_ShouldUpdateRoomWhenRoomInputIsValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := valueobject.NewID()
	adminId, _ := valueobject.NewAuth0IDWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Need for Speed")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	createdAt := valueobject.NewTimestamp()
	updatedAt, _ := valueobject.NewTimestampWith(createdAt.Value())

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
			assert.True(t, createdAt.Time().Equal(r.CreatedAt().Time()))
			assert.True(t, updatedAt.Time().Before(r.UpdatedAt().Time()))
		}).
		Return(nil).
		Times(1)

	useCase := NewUpdateRoomUseCase(gateway)
	err := useCase.Execute(ctx, &input)
	assert.Nil(t, err)
}

func TestUpdateRoomUseCase_ShouldReturnANotFoundErrorWhenRoomIdDoesNotExist(t *testing.T) {
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
	assert.ErrorIs(t, err, validation.ErrNotFoundRoom)
}

func TestUpdateRoomUseCase_ShouldReturnAnAdminErrorWhenAdminIdIsNotRoomAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := valueobject.NewID()
	adminId, _ := valueobject.NewAuth0IDWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Need for Speed")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	createdAt := valueobject.NewTimestamp()
	updatedAt, _ := valueobject.NewTimestampWith(createdAt.Value())

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
	assert.ErrorIs(t, err, validation.ErrInvalidRoomAdmin)
}

func TestUpdateRoomUseCase_ShouldReturnAnErrorWhenInputIsInvalid(t *testing.T) {
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
			err: validation.ErrInvalidId,
		},
		{
			input: &UpdateRoomUseCaseInput{
				Id:       valueobject.NewID().Value(),
				AdminId:  "fdadaervc233",
				Name:     "",
				Category: "",
			},
			err: validation.ErrInvalidId,
		},
		{
			input: &UpdateRoomUseCaseInput{
				Id:       valueobject.NewID().Value(),
				AdminId:  "auth0|64c8457bb160e37c8c34533c",
				Name:     "",
				Category: "",
			},
			err: validation.ErrRequiredRoomName,
		},
		{
			input: &UpdateRoomUseCaseInput{
				Id:       valueobject.NewID().Value(),
				AdminId:  "auth0|64c8457bb160e37c8c34533c",
				Name:     "Rust",
				Category: "Anything",
			},
			err: validation.ErrInvalidRoomCategory,
		},
	}

	ctx := context.Background()
	gateway := mock.NewMockRoomGatewayInterface(ctrl)
	useCase := NewUpdateRoomUseCase(gateway)

	for _, test := range testCases {
		err := useCase.Execute(ctx, test.input)
		assert.ErrorIs(t, err, test.err)
	}
}
