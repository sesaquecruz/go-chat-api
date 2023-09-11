package model

import (
	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
)

type RoomModel struct {
	Id        string
	AdminId   string
	Name      string
	Category  string
	CreatedAt string
	UpdatedAt string
	DeletedAt *string
}

func NewRoomModel(room *entity.Room) *RoomModel {
	model := RoomModel{}

	model.Id = room.Id().Value()
	model.AdminId = room.AdminId().Value()
	model.Name = room.Name().Value()
	model.Category = room.Category().Value()
	model.CreatedAt = room.CreatedAt().Value()
	model.UpdatedAt = room.UpdatedAt().Value()

	if room.DeletedAt() != nil {
		deleteAt := room.DeletedAt().Value()
		model.DeletedAt = &deleteAt
	}

	return &model
}

func (m *RoomModel) ToEntity() (*entity.Room, error) {
	id, err := valueobject.NewIdWith(m.Id)
	if err != nil {
		return nil, err
	}

	adminId, err := valueobject.NewUserIdWith(m.AdminId)
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

	createdAt, err := valueobject.NewTimestampWith(m.CreatedAt)
	if err != nil {
		return nil, err
	}

	updatedAt, err := valueobject.NewTimestampWith(m.UpdatedAt)
	if err != nil {
		return nil, err
	}

	var deletedAt *valueobject.Timestamp = nil

	if m.DeletedAt != nil {
		deletedAt, err = valueobject.NewTimestampWith(*m.DeletedAt)
		if err != nil {
			return nil, err
		}
	}

	room := entity.NewRoomWith(id, adminId, name, category, createdAt, updatedAt, deletedAt)

	return room, nil
}
