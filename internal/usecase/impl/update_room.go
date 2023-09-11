package impl

import (
	"context"
	"errors"

	"github.com/sesaquecruz/go-chat-api/internal/domain/repository"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
)

type UpdateRoomUseCase struct {
	roomRepository repository.RoomRepository
	logger         *log.Logger
}

func NewUpdateRoomUseCase(roomRepository repository.RoomRepository) *UpdateRoomUseCase {
	return &UpdateRoomUseCase{
		roomRepository: roomRepository,
		logger:         log.NewLogger("UpdateRoomUseCase"),
	}
}

func (u *UpdateRoomUseCase) Execute(ctx context.Context, input *usecase.UpdateRoomUseCaseInput) error {
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

	room.UpdateName(name)
	room.UpdateCategory(category)

	err = u.roomRepository.Update(ctx, room)
	if err != nil {
		u.logger.Error(err)
		return err
	}

	return nil
}
