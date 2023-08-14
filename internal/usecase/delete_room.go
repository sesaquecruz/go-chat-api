package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sesaquecruz/go-chat-api/internal/domain/gateway"
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
	roomGateway gateway.RoomGatewayInterface
	logger      *log.Logger
}

func NewDeleteRoomUseCase(roomGateway gateway.RoomGatewayInterface) *DeleteRoomUseCase {
	return &DeleteRoomUseCase{
		roomGateway: roomGateway,
		logger:      log.NewLogger("DeleteRoomUseCase"),
	}
}

func (u *DeleteRoomUseCase) Execute(ctx context.Context, input *DeleteRoomUseCaseInput) error {
	id, err := valueobject.NewIDWith(input.Id)
	if err != nil {
		return err
	}

	adminId, err := valueobject.NewAuth0IDWith(input.AdminId)
	if err != nil {
		return err
	}

	room, err := u.roomGateway.FindById(ctx, id)
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

	if err := u.roomGateway.Delete(ctx, id); err != nil {
		u.logger.Error(err)
		return validation.NewInternalError(err)
	}

	return nil
}
