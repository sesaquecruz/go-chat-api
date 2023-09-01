package model

import (
	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
)

type RoomModel struct {
	Id        string
	Admin_id  string
	Name      string
	Category  string
	CreatedAt string
	UpdatedAt string
}

func (m *RoomModel) ToEntity() (*entity.Room, error) {
	id, err := valueobject.NewIdWith(m.Id)
	if err != nil {
		return nil, err
	}

	adminId, err := valueobject.NewUserIdWith(m.Admin_id)
	if err != nil {
		return nil, err
	}

	name, err := valueobject.NewRoomNameWith(m.Name)
	if err != nil {
		return nil, err
	}

	category, err := valueobject.NewRoomCategoryWith(m.Category)
	if err != nil {
		return nil, err
	}

	created_at, err := valueobject.NewTimestampWith(m.CreatedAt)
	if err != nil {
		return nil, err
	}

	updated_at, err := valueobject.NewTimestampWith(m.UpdatedAt)
	if err != nil {
		return nil, err
	}

	room := entity.NewRoomWith(id, adminId, name, category, created_at, updated_at)

	return room, nil
}
