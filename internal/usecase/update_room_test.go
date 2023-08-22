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

func TestUpdateRoomUseCase_ShouldUpdateARoomWhenRoomInputIsValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := valueobject.NewId()
	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Need for Speed")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	createdAt := valueobject.NewTimestamp()
	updatedAt, _ := valueobject.NewTimestampWith(createdAt.String())

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

	repository := mock.NewMockRoomRepositoryInterface(ctrl)
	repository.
		EXPECT().
		FindById(gomock.Any(), gomock.Any()).
		Do(func(c context.Context, i *valueobject.Id) {
			assert.Equal(t, ctx, c)
			assert.Equal(t, id.Value(), i.Value())
		}).
		Return(room, nil).
		Times(1)

	repository.
		EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Do(func(c context.Context, r *entity.Room) {
			assert.Equal(t, ctx, c)
			assert.Equal(t, id.Value(), r.Id().Value())
			assert.Equal(t, adminId.Value(), r.AdminId().Value())
			assert.Equal(t, input.Name, r.Name().Value())
			assert.Equal(t, input.Category, r.Category().Value())
			assert.True(t, createdAt.Value().Equal(r.CreatedAt().Value()))
			assert.True(t, updatedAt.Value().Before(r.UpdatedAt().Value()))
		}).
		Return(nil).
		Times(1)

	useCase := NewUpdateRoomUseCase(repository)
	err := useCase.Execute(ctx, &input)
	assert.Nil(t, err)
}

func TestUpdateRoomUseCase_ShouldReturnANotFoundErrorWhenIdDoesNotExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	input := UpdateRoomUseCaseInput{
		Id:       valueobject.NewId().Value(),
		AdminId:  "auth0|64c8457bb160e37c8c34533b",
		Name:     "Rust",
		Category: "Tech",
	}

	repository := mock.NewMockRoomRepositoryInterface(ctrl)
	repository.
		EXPECT().
		FindById(gomock.Any(), gomock.Any()).
		Return(nil, sql.ErrNoRows).
		Times(1)

	useCase := NewUpdateRoomUseCase(repository)
	err := useCase.Execute(ctx, &input)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, validation.ErrNotFoundRoom)
}

func TestUpdateRoomUseCase_ShouldReturnAnAdminErrorWhenAdminIdIsInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := valueobject.NewId()
	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Need for Speed")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	createdAt := valueobject.NewTimestamp()
	updatedAt, _ := valueobject.NewTimestampWith(createdAt.String())

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

	repository := mock.NewMockRoomRepositoryInterface(ctrl)
	repository.
		EXPECT().
		FindById(gomock.Any(), gomock.Any()).
		Do(func(c context.Context, i *valueobject.Id) {
			assert.Equal(t, ctx, c)
			assert.Equal(t, id.Value(), i.Value())
		}).
		Return(room, nil).
		Times(1)

	useCase := NewUpdateRoomUseCase(repository)
	err := useCase.Execute(ctx, &input)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, validation.ErrInvalidRoomAdmin)
}

func TestUpdateRoomUseCase_ShouldReturnAnErrorWhenInputIsInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		err   error
		input *UpdateRoomUseCaseInput
	}{
		{
			err: validation.ErrInvalidId,
			input: &UpdateRoomUseCaseInput{
				Id:       "fdfagferr",
				AdminId:  "",
				Name:     "",
				Category: "",
			},
		},
		{
			err: validation.ErrInvalidUserId,
			input: &UpdateRoomUseCaseInput{
				Id:       valueobject.NewId().Value(),
				AdminId:  "fdadaervc233",
				Name:     "",
				Category: "",
			},
		},
		{
			err: validation.ErrRequiredRoomName,
			input: &UpdateRoomUseCaseInput{
				Id:       valueobject.NewId().Value(),
				AdminId:  "auth0|64c8457bb160e37c8c34533c",
				Name:     "",
				Category: "",
			},
		},
		{
			err: validation.ErrInvalidRoomCategory,
			input: &UpdateRoomUseCaseInput{
				Id:       valueobject.NewId().Value(),
				AdminId:  "auth0|64c8457bb160e37c8c34533c",
				Name:     "Rust",
				Category: "Anything",
			},
		},
	}

	ctx := context.Background()
	repository := mock.NewMockRoomRepositoryInterface(ctrl)
	useCase := NewUpdateRoomUseCase(repository)

	for _, test := range testCases {
		err := useCase.Execute(ctx, test.input)
		assert.ErrorIs(t, err, test.err)
	}
}
