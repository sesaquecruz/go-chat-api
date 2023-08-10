package database

import (
	"context"
	"database/sql"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
)

type RoomPostgresGateway struct {
	db     *sql.DB
	logger *log.Logger
}

func NewRoomPostgresGateway(db *sql.DB) *RoomPostgresGateway {
	return &RoomPostgresGateway{
		db:     db,
		logger: log.NewLogger("RoomPostgresGateway"),
	}
}

func (g *RoomPostgresGateway) Save(ctx context.Context, room *entity.Room) error {
	stmt, err := g.db.PrepareContext(ctx, "INSERT INTO rooms (id, admin_id, name, category, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		g.logger.Error(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		room.Id().Value(),
		room.AdminId().Value(),
		room.Name().Value(),
		room.Category().Value(),
		room.CreatedAt().StringValue(),
		room.UpdatedAt().StringValue(),
	)
	if err != nil {
		g.logger.Error(err)
		return err
	}

	return nil
}

func (g *RoomPostgresGateway) FindById(ctx context.Context, id *valueobject.ID) (*entity.Room, error) {
	stmt, err := g.db.PrepareContext(ctx, "SELECT id, admin_id, name, category, created_at, updated_at FROM rooms WHERE id = $1")
	if err != nil {
		g.logger.Error(err)
		return nil, err
	}
	defer stmt.Close()

	var r RoomModel
	err = stmt.QueryRowContext(ctx, id.Value()).Scan(
		&r.Id,
		&r.Admin_id,
		&r.Name,
		&r.Category,
		&r.Created_at,
		&r.Updated_at,
	)
	if err != nil {
		g.logger.Error(err)
		return nil, err
	}

	return r.ToEntity()
}
