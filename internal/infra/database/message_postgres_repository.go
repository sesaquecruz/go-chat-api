package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/repository"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/internal/infra/database/model"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
)

type MessagePostgresRepository struct {
	db     *sql.DB
	logger *log.Logger
}

func NewMessagePostgresRepository(db *sql.DB) *MessagePostgresRepository {
	return &MessagePostgresRepository{
		db:     db,
		logger: log.NewLogger("MessagePostgresRepository"),
	}
}

func (r *MessagePostgresRepository) Save(ctx context.Context, message *entity.Message) error {
	stmt, err := r.db.PrepareContext(ctx, `
		INSERT INTO messages (id, room_id, sender_id, sender_name, text, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
	`)
	if err != nil {
		r.logger.Error(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		message.Id().Value(),
		message.RoomId().Value(),
		message.SenderId().Value(),
		message.SenderName().Value(),
		message.Text().Value(),
		message.CreatedAt().Value(),
	)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *MessagePostgresRepository) FindById(ctx context.Context, id *valueobject.Id) (*entity.Message, error) {
	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, room_id, sender_id, sender_name, text, created_at 
		FROM messages 
		WHERE id = $1
	`)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	defer stmt.Close()

	var m model.MessageModel

	err = stmt.QueryRowContext(ctx, id.Value()).Scan(
		&m.Id,
		&m.RoomId,
		&m.SenderId,
		&m.SenderName,
		&m.Text,
		&m.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFoundMessage
		}

		r.logger.Error(err)
		return nil, err
	}

	message, err := m.ToEntity()
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}

	return message, nil
}
