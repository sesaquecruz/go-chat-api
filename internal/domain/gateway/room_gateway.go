package gateway

import (
	"context"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
)

const ErrInvalidRoom = "room is invalid"
const ErrAlreadyExistsRoom = "room already exists"

type RoomGatewayInterface interface {
	Save(ctx context.Context, room *entity.Room) error
	FindById(ctx context.Context, id *valueobject.ID) (*entity.Room, error)
}
