package impl

import (
	"context"
	"errors"

	"github.com/sesaquecruz/go-chat-api/internal/domain/repository"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
)

type FindRoomUseCase struct {
	roomRepository repository.RoomRepository
	logger         *log.Logger
}

func NewFindRoomUseCase(roomRepository repository.RoomRepository) *FindRoomUseCase {
	return &FindRoomUseCase{
		roomRepository: roomRepository,
		logger:         log.NewLogger("FindRoomUseCase"),
	}
}

func (u *FindRoomUseCase) Execute(
	ctx context.Context,
	input *usecase.FindRoomUseCaseInput,
) (*usecase.FindRoomUseCaseOutput, error) {

	id, err := valueobject.NewIdWith(input.Id)
	if err != nil {
		return nil, err
	}

	room, err := u.roomRepository.FindById(ctx, id)
	if err != nil {
		if !errors.Is(err, repository.ErrNotFoundRoom) {
			u.logger.Error(err)
		}

		return nil, err
	}

	if room.IsDeleted() {
		return nil, repository.ErrNotFoundRoom
	}

	output := &usecase.FindRoomUseCaseOutput{
		Id:        room.Id().Value(),
		AdminId:   room.AdminId().Value(),
		Name:      room.Name().Value(),
		Category:  room.Category().Value(),
		CreatedAt: room.CreatedAt().Value(),
		UpdatedAt: room.UpdatedAt().Value(),
	}

	return output, nil
}
