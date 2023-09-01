package impl

import (
	"context"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/repository"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
)

type CreateRoomUseCase struct {
	roomRepository repository.RoomRepository
	logger         *log.Logger
}

func NewCreateRoomUseCase(roomRepository repository.RoomRepository) *CreateRoomUseCase {
	return &CreateRoomUseCase{
		roomRepository: roomRepository,
		logger:         log.NewLogger("CreateRoomUseCase"),
	}
}

func (u *CreateRoomUseCase) Execute(
	ctx context.Context,
	input *usecase.CreateRoomUseCaseInput,
) (*usecase.CreateRoomUseCaseOutput, error) {

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

	room := entity.NewRoom(adminId, name, category)

	err = u.roomRepository.Save(ctx, room)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	output := &usecase.CreateRoomUseCaseOutput{
		RoomId: room.Id().Value(),
	}

	return output, nil
}
