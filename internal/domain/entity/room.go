package entity

import (
	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
)

type Room struct {
	id        *valueobject.ID
	adminId   *valueobject.Auth0ID
	name      *valueobject.RoomName
	category  *valueobject.RoomCategory
	createdAt *valueobject.Timestamp
	updatedAt *valueobject.Timestamp
}

func NewRoom(
	adminId *valueobject.Auth0ID,
	name *valueobject.RoomName,
	category *valueobject.RoomCategory,
) (*Room, error) {
	now := valueobject.NewTimestamp()
	room := &Room{
		id:        valueobject.NewID(),
		adminId:   adminId,
		name:      name,
		category:  category,
		createdAt: now,
		updatedAt: now,
	}

	if err := room.Validate(); err != nil {
		return nil, err
	}

	return room, nil
}

func NewRoomWith(
	id *valueobject.ID,
	adminId *valueobject.Auth0ID,
	name *valueobject.RoomName,
	category *valueobject.RoomCategory,
	createdAt *valueobject.Timestamp,
	updatedAt *valueobject.Timestamp,
) (*Room, error) {
	room := &Room{
		id:        id,
		adminId:   adminId,
		name:      name,
		category:  category,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}

	if err := room.Validate(); err != nil {
		return nil, err
	}

	return room, nil
}

func (r *Room) Validate() error {
	if r.id == nil {
		return validation.ErrRequiredRoomId
	}
	if r.adminId == nil {
		return validation.ErrRequiredRoomAdminId
	}
	if r.name == nil {
		return validation.ErrRequiredRoomName
	}
	if r.category == nil {
		return validation.ErrRequiredRoomCategory
	}
	if r.createdAt == nil {
		return validation.ErrRequiredRoomCreatedAt
	}
	if r.updatedAt == nil {
		return validation.ErrRequiredRoomUpdatedAt
	}

	return nil
}

func (r *Room) Id() *valueobject.ID {
	return r.id
}

func (r *Room) AdminId() *valueobject.Auth0ID {
	return r.adminId
}

func (r *Room) Name() *valueobject.RoomName {
	return r.name
}

func (r *Room) Category() *valueobject.RoomCategory {
	return r.category
}

func (r *Room) CreatedAt() *valueobject.Timestamp {
	return r.createdAt
}

func (r *Room) UpdatedAt() *valueobject.Timestamp {
	return r.updatedAt
}

func (r *Room) UpdateName(name *valueobject.RoomName) error {
	if name == nil {
		return validation.ErrRequiredRoomName
	}

	r.name = name
	r.updatedAt = valueobject.NewTimestamp()
	return nil
}

func (r *Room) UpdateCategory(category *valueobject.RoomCategory) error {
	if category == nil {
		return validation.ErrRequiredRoomCategory
	}

	r.category = category
	r.updatedAt = valueobject.NewTimestamp()
	return nil
}
