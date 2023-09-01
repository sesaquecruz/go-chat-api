package repository

import (
	"context"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
)

const ErrNotFoundMessage = validation.NotFoundError("message not found")

type MessageRepository interface {
	Save(ctx context.Context, message *entity.Message) error
	FindById(ctx context.Context, id *valueobject.Id) (*entity.Message, error)
}
