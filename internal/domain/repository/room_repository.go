package repository

import (
	"context"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/pagination"
	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
)

const ErrNotFoundRoom = validation.NotFoundError("room not found")

type RoomRepository interface {
	Save(ctx context.Context, room *entity.Room) error
	FindById(ctx context.Context, id *valueobject.Id) (*entity.Room, error)
	Search(ctx context.Context, query *pagination.Query) (*pagination.Page[*entity.Room], error)
	Update(ctx context.Context, room *entity.Room) error
	Delete(ctx context.Context, id *valueobject.Id) error
}
