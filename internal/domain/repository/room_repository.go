package repository

import (
	"context"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/repository/pagination"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
)

type RoomRepositoryInterface interface {
	Save(ctx context.Context, room *entity.Room) error
	FindById(ctx context.Context, id *valueobject.Id) (*entity.Room, error)
	Search(ctx context.Context, query *pagination.Query) (*pagination.Page[*entity.Room], error)
	Update(ctx context.Context, room *entity.Room) error
	Delete(ctx context.Context, id *valueobject.Id) error
}
