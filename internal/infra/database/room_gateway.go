package database

import (
	"context"
	"database/sql"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/search"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
)

type RoomGateway struct {
	db     *sql.DB
	logger *log.Logger
}

func NewRoomGateway(db *sql.DB) *RoomGateway {
	return &RoomGateway{
		db:     db,
		logger: log.NewLogger("RoomGateway"),
	}
}

func (g *RoomGateway) Save(ctx context.Context, room *entity.Room) error {
	stmt, err := g.db.PrepareContext(ctx, `
		INSERT INTO rooms (id, admin_id, name, category, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
	`)
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
		room.CreatedAt().Value(),
		room.UpdatedAt().Value(),
	)
	if err != nil {
		g.logger.Error(err)
		return err
	}

	return nil
}

func (g *RoomGateway) FindById(ctx context.Context, id *valueobject.ID) (*entity.Room, error) {
	stmt, err := g.db.PrepareContext(ctx, `
		SELECT id, admin_id, name, category, created_at, updated_at 
		FROM rooms 
		WHERE id = $1
	`)
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

	room, err := r.ToEntity()
	if err != nil {
		g.logger.Error(err)
		return nil, err
	}

	return room, nil
}

func (g *RoomGateway) Search(ctx context.Context, query *search.Query) (*search.Page[*entity.Room], error) {
	stmt1, err := g.db.PrepareContext(ctx, `
		SELECT id, admin_id, name, category, created_at, updated_at 
		FROM rooms 
		WHERE ($1 = '') OR (UPPER(name) LIKE '%' || $1 || '%') OR (UPPER(category::text) LIKE '%' || $1 || '%') 
		ORDER BY name `+query.Sort()+`
		LIMIT $2 
		OFFSET $3
	`)
	if err != nil {
		g.logger.Error(err)
		return nil, err
	}
	defer stmt1.Close()

	stmt2, err := g.db.PrepareContext(ctx, `
		SELECT COUNT(*) 
		FROM rooms 
		WHERE ($1 = '') OR (UPPER(name) LIKE '%' || $1 || '%') OR (UPPER(category::text) LIKE '%' || $1 || '%')
	`)
	if err != nil {
		g.logger.Error(err)
		return nil, err
	}
	defer stmt2.Close()

	rows, err := stmt1.QueryContext(ctx, query.Search(), query.Size(), query.Size()*query.Page())
	if err != nil {
		g.logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	var items []*entity.Room
	for rows.Next() {
		var r RoomModel

		err := rows.Scan(
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

		e, err := r.ToEntity()
		if err != nil {
			g.logger.Error(err)
			return nil, err
		}

		items = append(items, e)
	}

	var total int64
	err = stmt2.QueryRowContext(ctx, query.Search()).Scan(&total)
	if err != nil {
		g.logger.Error(err)
		return nil, err
	}

	page := search.NewPage[*entity.Room](query.Page(), query.Size(), total, items)
	return page, nil
}

func (g *RoomGateway) Update(ctx context.Context, room *entity.Room) error {
	stmt, err := g.db.PrepareContext(ctx, `
		UPDATE rooms 
		SET id = $1, admin_id = $2, name = $3, category = $4, created_at = $5, updated_at = $6
	`)
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
		room.CreatedAt().Value(),
		room.UpdatedAt().Value(),
	)
	if err != nil {
		g.logger.Error(err)
		return err
	}

	return nil
}

func (g *RoomGateway) Delete(ctx context.Context, id *valueobject.ID) error {
	stmt, err := g.db.PrepareContext(ctx, `
		DELETE FROM rooms 
		WHERE id = $1
	`)
	if err != nil {
		g.logger.Error(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id.Value())
	if err != nil {
		g.logger.Error(err)
		return err
	}

	return nil
}
