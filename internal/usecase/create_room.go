package usecase

import (
	"context"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/errors"
	"github.com/sesaquecruz/go-chat-api/internal/domain/gateway"
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
	roomGateway gateway.RoomGatewayInterface
	logger      *log.Logger
}

func NewCreateRoomUseCase(roomGateway gateway.RoomGatewayInterface) *CreateRoomUseCase {
	return &CreateRoomUseCase{
		roomGateway: roomGateway,
		logger:      log.NewLogger("CreateRoomUseCase"),
	}
}

func (u *CreateRoomUseCase) Execute(ctx context.Context, input *CreateRoomUseCaseInput) (*CreateRoomUseCaseOutput, error) {
	adminId, err := valueobject.NewAuth0IDWith(input.AdminId)
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

	err = u.roomGateway.Save(ctx, room)
	if err != nil {
		u.logger.Error(err)
		return nil, errors.NewGatewayError(err.Error())
	}

	return &CreateRoomUseCaseOutput{
		RoomId: room.Id().Value(),
	}, nil
}
