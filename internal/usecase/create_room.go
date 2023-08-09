package usecase

import (
	"context"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/gateway"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/pkg"
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
	logger      *pkg.Logger
}

func NewCreateRoomUseCase(roomGateway gateway.RoomGatewayInterface) *CreateRoomUseCase {
	return &CreateRoomUseCase{
		roomGateway: roomGateway,
		logger:      pkg.NewLogger("CreateRoomUseCase"),
	}
}

func (uc *CreateRoomUseCase) Execute(ctx context.Context, input *CreateRoomUseCaseInput) (*CreateRoomUseCaseOutput, error) {
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
		uc.logger.Error(err)
		return nil, err
	}

	err = uc.roomGateway.Save(ctx, room)
	if err != nil {
		uc.logger.Error(err)
		return nil, err
	}

	return &CreateRoomUseCaseOutput{
		RoomId: room.Id().Value(),
	}, nil
}
