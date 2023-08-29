package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/pagination"
	"github.com/sesaquecruz/go-chat-api/internal/domain/repository"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/internal/infra/database/model"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
)

type RoomPostgresRepository struct {
	db     *sql.DB
	logger *log.Logger
}

func NewRoomPostgresRepository(db *sql.DB) *RoomPostgresRepository {
	return &RoomPostgresRepository{
		db:     db,
		logger: log.NewLogger("RoomPostgresRepository"),
	}
}

func (r *RoomPostgresRepository) Save(ctx context.Context, room *entity.Room) error {
	stmt, err := r.db.PrepareContext(ctx, `
		INSERT INTO rooms (id, admin_id, name, category, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
	`)
	if err != nil {
		r.logger.Error(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		room.Id().Value(),
		room.AdminId().Value(),
		room.Name().Value(),
		room.Category().Value(),
		room.CreatedAt().Value(),
		room.UpdatedAt().Value(),
	)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *RoomPostgresRepository) FindById(ctx context.Context, id *valueobject.Id) (*entity.Room, error) {
	stmt, err := r.db.PrepareContext(ctx, `
		SELECT id, admin_id, name, category, created_at, updated_at 
		FROM rooms 
		WHERE id = $1
	`)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	defer stmt.Close()

	var m model.RoomModel

	err = stmt.QueryRowContext(ctx, id.Value()).Scan(
		&m.Id,
		&m.Admin_id,
		&m.Name,
		&m.Category,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFoundRoom
		}

		r.logger.Error(err)
		return nil, err
	}

	room, err := m.ToEntity()
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}

	return room, nil
}

func (r *RoomPostgresRepository) Search(ctx context.Context, query *pagination.Query) (*pagination.Page[*entity.Room], error) {
	stmt1, err := r.db.PrepareContext(ctx, `
		SELECT id, admin_id, name, category, created_at, updated_at 
		FROM rooms 
		WHERE ($1 = '') OR (UPPER(name) LIKE '%' || $1 || '%') OR (UPPER(category::text) LIKE '%' || $1 || '%') 
		ORDER BY name `+query.Sort()+`
		LIMIT $2 
		OFFSET $3
	`)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	defer stmt1.Close()

	stmt2, err := r.db.PrepareContext(ctx, `
		SELECT COUNT(*) 
		FROM rooms 
		WHERE ($1 = '') OR (UPPER(name) LIKE '%' || $1 || '%') OR (UPPER(category::text) LIKE '%' || $1 || '%')
	`)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	defer stmt2.Close()

	rows, err := stmt1.QueryContext(ctx, query.Search(), query.Size(), query.Size()*query.Page())
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	var items []*entity.Room

	for rows.Next() {
		var m model.RoomModel

		err := rows.Scan(
			&m.Id,
			&m.Admin_id,
			&m.Name,
			&m.Category,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			r.logger.Error(err)
			return nil, err
		}

		room, err := m.ToEntity()
		if err != nil {
			r.logger.Error(err)
			return nil, err
		}

		items = append(items, room)
	}

	var total int64

	err = stmt2.QueryRowContext(ctx, query.Search()).Scan(&total)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}

	page := pagination.NewPage[*entity.Room](query.Page(), query.Size(), total, items)
	return page, nil
}

func (r *RoomPostgresRepository) Update(ctx context.Context, room *entity.Room) error {
	stmt, err := r.db.PrepareContext(ctx, `
		UPDATE rooms 
		SET admin_id = $2, name = $3, category = $4, created_at = $5, updated_at = $6
		WHERE id = $1
	`)
	if err != nil {
		r.logger.Error(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		room.Id().Value(),
		room.AdminId().Value(),
		room.Name().Value(),
		room.Category().Value(),
		room.CreatedAt().Value(),
		room.UpdatedAt().Value(),
	)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *RoomPostgresRepository) Delete(ctx context.Context, id *valueobject.Id) error {
	stmt, err := r.db.PrepareContext(ctx, `
		DELETE FROM rooms 
		WHERE id = $1
	`)
	if err != nil {
		r.logger.Error(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id.Value())
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}
