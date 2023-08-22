package usecase

import (
	"context"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/repository"
	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
)

type CreateRoomUseCaseInput struct {
	AdminId  string
	Name     string
	Category string
}

type CreateRoomUseCaseOutput struct {
	RoomId string
}

type CreateRoomUseCaseInterface interface {
	Execute(ctx context.Context, input *CreateRoomUseCaseInput) (*CreateRoomUseCaseOutput, error)
}

type CreateRoomUseCase struct {
	roomRepository repository.RoomRepositoryInterface
	logger         *log.Logger
}

func NewCreateRoomUseCase(roomRepository repository.RoomRepositoryInterface) *CreateRoomUseCase {
	return &CreateRoomUseCase{
		roomRepository: roomRepository,
		logger:         log.NewLogger("CreateRoomUseCase"),
	}
}

func (u *CreateRoomUseCase) Execute(ctx context.Context, input *CreateRoomUseCaseInput) (*CreateRoomUseCaseOutput, error) {
	adminId, err := valueobject.NewUserIdWith(input.AdminId)
	if err != nil {
		return nil, err
	}

	name, err := valueobject.NewRoomNameWith(input.Name)
	if err != nil {
		return nil, err
	}

	category, err := valueobject.NewRoomCategoryWith(input.Category)
	if err != nil {
		return nil, err
	}

	room, err := entity.NewRoom(adminId, name, category)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	err = u.roomRepository.Save(ctx, room)
	if err != nil {
		u.logger.Error(err)
		return nil, validation.NewInternalError(err)
	}

	output := &CreateRoomUseCaseOutput{
		RoomId: room.Id().Value(),
	}

	return output, nil
}
