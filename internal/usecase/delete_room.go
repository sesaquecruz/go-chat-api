package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sesaquecruz/go-chat-api/internal/domain/repository"
	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
)

type DeleteRoomUseCaseInput struct {
	Id      string
	AdminId string
}

type DeleteRoomUseCaseInterface interface {
	Execute(ctx context.Context, input *DeleteRoomUseCaseInput) error
}

type DeleteRoomUseCase struct {
	roomRepository repository.RoomRepositoryInterface
	logger         *log.Logger
}

func NewDeleteRoomUseCase(roomRepository repository.RoomRepositoryInterface) *DeleteRoomUseCase {
	return &DeleteRoomUseCase{
		roomRepository: roomRepository,
		logger:         log.NewLogger("DeleteRoomUseCase"),
	}
}

func (u *DeleteRoomUseCase) Execute(ctx context.Context, input *DeleteRoomUseCaseInput) error {
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
		if errors.Is(err, sql.ErrNoRows) {
			return validation.ErrNotFoundRoom
		}

		u.logger.Error(err)
		return validation.NewInternalError(err)
	}

	if room.AdminId().Value() != adminId.Value() {
		return validation.ErrInvalidRoomAdmin
	}

	err = u.roomRepository.Delete(ctx, id)
	if err != nil {
		u.logger.Error(err)
		return validation.NewInternalError(err)
	}

	return nil
}
