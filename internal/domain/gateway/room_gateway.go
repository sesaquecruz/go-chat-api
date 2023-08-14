package gateway

import (
	"context"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
)

type RoomGatewayInterface interface {
	Save(ctx context.Context, room *entity.Room) error
	FindById(ctx context.Context, id *valueobject.ID) (*entity.Room, error)
	Update(ctx context.Context, room *entity.Room) error
	Delete(ctx context.Context, id *valueobject.ID) error
}
