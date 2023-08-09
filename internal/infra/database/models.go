package database

import (
	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
)

type RoomModel struct {
	Id         string
	Admin_id   string
	Name       string
	Category   string
	Created_at string
	Updated_at string
}

func (rm *RoomModel) toEntity() (*entity.Room, error) {
	id, err := valueobject.NewIDWith(rm.Id)
	if err != nil {
		return nil, err
	}

	adminId, err := valueobject.NewAuth0IDWith(rm.Admin_id)
	if err != nil {
		return nil, err
	}

	name, err := valueobject.NewRoomNameWith(rm.Name)
	if err != nil {
		return nil, err
	}

	category, err := valueobject.NewRoomCategoryWith(rm.Category)
	if err != nil {
		return nil, err
	}

	created_at, err := valueobject.NewDateTimeWith(rm.Created_at)
	if err != nil {
		return nil, err
	}

	updated_at, err := valueobject.NewDateTimeWith(rm.Updated_at)
	if err != nil {
		return nil, err
	}

	return entity.NewRoomWith(id, adminId, name, category, created_at, updated_at)
}
