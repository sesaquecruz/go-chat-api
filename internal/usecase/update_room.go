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

type UpdateRoomUseCaseInput struct {
	Id       string
	AdminId  string
	Name     string
	Category string
}

type UpdateRoomUseCaseInterface interface {
	Execute(ctx context.Context, input *UpdateRoomUseCaseInput) error
}

type UpdateRoomUseCase struct {
	roomRepository repository.RoomRepositoryInterface
	logger         *log.Logger
}

func NewUpdateRoomUseCase(roomRepository repository.RoomRepositoryInterface) *UpdateRoomUseCase {
	return &UpdateRoomUseCase{
		roomRepository: roomRepository,
		logger:         log.NewLogger("UpdateRoomUseCase"),
	}
}

func (u *UpdateRoomUseCase) Execute(ctx context.Context, input *UpdateRoomUseCaseInput) error {
	id, err := valueobject.NewIdWith(input.Id)
	if err != nil {
		return err
	}

	adminId, err := valueobject.NewUserIdWith(input.AdminId)
	if err != nil {
		return err
	}

	name, err := valueobject.NewRoomNameWith(input.Name)
	if err != nil {
		return err
	}

	category, err := valueobject.NewRoomCategoryWith(input.Category)
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

	if adminId.Value() != room.AdminId().Value() {
		return validation.ErrInvalidRoomAdmin
	}

	err = room.UpdateName(name)
	if err != nil {
		u.logger.Error(err)
		return validation.NewInternalError(err)
	}

	err = room.UpdateCategory(category)
	if err != nil {
		u.logger.Error(err)
		return validation.NewInternalError(err)
	}

	err = u.roomRepository.Update(ctx, room)
	if err != nil {
		u.logger.Error(err)
		return validation.NewInternalError(err)
	}

	return nil
}
