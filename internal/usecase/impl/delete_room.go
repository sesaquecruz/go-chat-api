package impl

import (
	"context"
	"errors"

	"github.com/sesaquecruz/go-chat-api/internal/domain/repository"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
)

type DeleteRoomUseCase struct {
	roomRepository repository.RoomRepository
	logger         *log.Logger
}

func NewDeleteRoomUseCase(roomRepository repository.RoomRepository) *DeleteRoomUseCase {
	return &DeleteRoomUseCase{
		roomRepository: roomRepository,
		logger:         log.NewLogger("DeleteRoomUseCase"),
	}
}

func (u *DeleteRoomUseCase) Execute(ctx context.Context, input *usecase.DeleteRoomUseCaseInput) error {
	id, err := valueobject.NewIdWith(input.Id)
	if err != nil {
		return err
	}

	adminId, err := valueobject.NewUserIdWith(input.AdminId)
	if err != nil {
		return err
	}

	room, err := u.roomRepository.FindById(ctx, id)
	if err != nil {
		if !errors.Is(err, repository.ErrNotFoundRoom) {
			u.logger.Error(err)
		}

		return err
	}

	if room.IsDeleted() {
		return repository.ErrNotFoundRoom
	}

	err = room.ValidateAdmin(adminId)
	if err != nil {
		return err
	}

	err = room.Delete()
	if err != nil {
		return err
	}

	err = u.roomRepository.Update(ctx, room)
	if err != nil {
		u.logger.Error(err)
		return err
	}

	return nil
}
