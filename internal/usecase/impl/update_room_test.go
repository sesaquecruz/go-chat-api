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

func TestUpdateRoomUseCase_ShouldUpdateARoomWhenDataIsValid(t *testing.T) {
	id := valueobject.NewId()
	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("A Game")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	createdAt := valueobject.NewTimestamp()
	updatedAt, _ := valueobject.NewTimestampWith(createdAt.Value())
	savedRoom := entity.NewRoomWith(id, adminId, name, category, createdAt, updatedAt)

	ctx := context.Background()
	input := &usecase.UpdateRoomUseCaseInput{
		Id:       id.Value(),
		AdminId:  adminId.Value(),
		Name:     "A Programming Language",
		Category: "Tech",
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
			assert.Equal(t, input.AdminId, r.AdminId().Value())
			assert.Equal(t, input.Name, r.Name().Value())
			assert.Equal(t, input.Category, r.Category().Value())
			assert.True(t, createdAt.Time().Equal(r.CreatedAt().Time()))
			assert.True(t, updatedAt.Time().Before(r.UpdatedAt().Time()))
		}).
		Return(nil).
		Once()

	useCase := NewUpdateRoomUseCase(roomRepository)

	err := useCase.Execute(ctx, input)
	assert.Nil(t, err)
}

func TestUpdateRoomUseCase_ShouldReturnAnErrorWhenDataIsInvalid(t *testing.T) {
	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("A Game")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	savedRoom := entity.NewRoom(adminId, name, category)

	ctx := context.Background()

	testCases := []struct {
		test  string
		input *usecase.UpdateRoomUseCaseInput
		err   error
	}{
		{
			"empty id",
			&usecase.UpdateRoomUseCaseInput{
				Id:       "",
				AdminId:  "auth0|64c8457bb160e37c8c34533c",
				Name:     "A Programming Language",
				Category: "Tech",
			},
			valueobject.ErrRequiredId,
		},
		{
			"empty admin id",
			&usecase.UpdateRoomUseCaseInput{
				Id:       "b3588483-4795-434a-877c-dcd158d6caa7",
				AdminId:  "",
				Name:     "A Programming Language",
				Category: "Tech",
			},
			valueobject.ErrRequiredUserId,
		},
		{
			"empty name",
			&usecase.UpdateRoomUseCaseInput{
				Id:       "b3588483-4795-434a-877c-dcd158d6caa7",
				AdminId:  "auth0|64c8457bb160e37c8c34533c",
				Name:     "",
				Category: "Tech",
			},
			valueobject.ErrRequiredRoomName,
		},
		{
			"empty category",
			&usecase.UpdateRoomUseCaseInput{
				Id:       "b3588483-4795-434a-877c-dcd158d6caa7",
				AdminId:  "auth0|64c8457bb160e37c8c34533c",
				Name:     "A Programming Language",
				Category: "",
			},
			valueobject.ErrRequiredRoomCategory,
		},
		{
			"invalid id",
			&usecase.UpdateRoomUseCaseInput{
				Id:       "fdafak12j17921",
				AdminId:  "auth0|64c8457bb160e37c8c34533c",
				Name:     "A Programming Language",
				Category: "Tech",
			},
			valueobject.ErrInvalidId,
		},
		{
			"invalid admin id",
			&usecase.UpdateRoomUseCaseInput{
				Id:       "b3588483-4795-434a-877c-dcd158d6caa7",
				AdminId:  "fadjkv89123192hfdf",
				Name:     "A Programming Language",
				Category: "Tech",
			},
			valueobject.ErrInvalidUserId,
		},
		{
			"invalid category",
			&usecase.UpdateRoomUseCaseInput{
				Id:       "b3588483-4795-434a-877c-dcd158d6caa7",
				AdminId:  "auth0|64c8457bb160e37c8c34533c",
				Name:     "A Programming Language",
				Category: "fadferiouk1j23",
			},
			valueobject.ErrInvalidRoomCategory,
		},
		{
			"invalid room admin",
			&usecase.UpdateRoomUseCaseInput{
				Id:       savedRoom.Id().Value(),
				AdminId:  "auth0|64c8457bb160e37c8c34533c",
				Name:     "A Programming Language",
				Category: "Tech",
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

	useCase := NewUpdateRoomUseCase(roomRepository)

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := useCase.Execute(ctx, tc.input)
			assert.NotNil(t, err)
			assert.ErrorIs(t, err, tc.err)
		})
	}
}

func TestUpdateRoomUseCase_ShouldReturnAnErrorOnRepositoryError(t *testing.T) {
	ctx := context.Background()
	input := &usecase.UpdateRoomUseCaseInput{
		Id:       "b3588483-4795-434a-877c-dcd158d6caa7",
		AdminId:  "auth0|64c8457bb160e37c8c34533c",
		Name:     "A Programming Language",
		Category: "Tech",
	}

	roomRepository := mocks.NewRoomRepositoryMock(t)

	roomRepository.
		EXPECT().
		FindById(mock.Anything, mock.Anything).
		Return(nil, repository.ErrNotFoundRoom).
		Once()

	useCase := NewUpdateRoomUseCase(roomRepository)

	err := useCase.Execute(ctx, input)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, repository.ErrNotFoundRoom)
}
