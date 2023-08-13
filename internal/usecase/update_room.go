package usecase

import (
	"context"
	"database/sql"
	"errors"

	domain_errors "github.com/sesaquecruz/go-chat-api/internal/domain/errors"
	"github.com/sesaquecruz/go-chat-api/internal/domain/gateway"
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
	roomGateway gateway.RoomGatewayInterface
	logger      *log.Logger
}

func NewUpdateRoomUseCase(roomGateway gateway.RoomGatewayInterface) *UpdateRoomUseCase {
	return &UpdateRoomUseCase{
		roomGateway: roomGateway,
		logger:      log.NewLogger("UpdateRoomUseCase"),
	}
}

func (u *UpdateRoomUseCase) Execute(ctx context.Context, input *UpdateRoomUseCaseInput) error {
	id, err := valueobject.NewIDWith(input.Id)
	if err != nil {
		return err
	}

	adminId, err := valueobject.NewAuth0IDWith(input.AdminId)
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

	room, err := u.roomGateway.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain_errors.NewGatewayError(gateway.ErrNotFoundRoom)
		}

		u.logger.Error(err)
		return domain_errors.NewGatewayError(err.Error())
	}

	if adminId.Value() != room.AdminId().Value() {
		return domain_errors.NewAuthorizationError("invalid room admin")
	}

	if err := room.UpdateName(name); err != nil {
		u.logger.Error(err)
		return err
	}

	if err := room.UpdateCategory(category); err != nil {
		u.logger.Error(err)
		return err
	}

	if err := u.roomGateway.Update(ctx, room); err != nil {
		return domain_errors.NewGatewayError(err.Error())
	}

	return nil
}
