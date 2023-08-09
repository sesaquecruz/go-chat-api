package database

import (
	"context"
	"database/sql"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/pkg"
)

type RoomPostgresGateway struct {
	db     *sql.DB
	logger *pkg.Logger
}

func NewRoomPostgresGateway(db *sql.DB) *RoomPostgresGateway {
	return &RoomPostgresGateway{
		db:     db,
		logger: pkg.NewLogger("RoomPostgresGateway"),
	}
}

func (rg *RoomPostgresGateway) Save(ctx context.Context, room *entity.Room) error {
	stmt, err := rg.db.PrepareContext(ctx, "INSERT INTO rooms (id, admin_id, name, category, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		rg.logger.Error(err)
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
		rg.logger.Error(err)
		return err
	}

	return nil
}

func (rg *RoomPostgresGateway) FindById(ctx context.Context, id *valueobject.ID) (*entity.Room, error) {
	stmt, err := rg.db.PrepareContext(ctx, "SELECT id, admin_id, name, category, created_at, updated_at FROM rooms WHERE id = $1")
	if err != nil {
		rg.logger.Error(err)
		return nil, err
	}
	defer stmt.Close()

	var rm RoomModel
	err = stmt.QueryRowContext(ctx, id.Value()).Scan(
		&rm.Id,
		&rm.Admin_id,
		&rm.Name,
		&rm.Category,
		&rm.Created_at,
		&rm.Updated_at,
	)
	if err != nil {
		rg.logger.Error(err)
		return nil, err
	}

	return rm.toEntity()
}
