package entity

import (
	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
)

type Room struct {
	id        *valueobject.Id
	adminId   *valueobject.UserId
	name      *valueobject.RoomName
	category  *valueobject.RoomCategory
	createdAt *valueobject.Timestamp
	updatedAt *valueobject.Timestamp
}

func NewRoom(
	adminId *valueobject.UserId,
	name *valueobject.RoomName,
	category *valueobject.RoomCategory,
) (*Room, error) {
	now := valueobject.NewTimestamp()

	return NewRoomWith(
		valueobject.NewId(),
		adminId,
		name,
		category,
		now,
		now,
	)
}

func NewRoomWith(
	id *valueobject.Id,
	adminId *valueobject.UserId,
	name *valueobject.RoomName,
	category *valueobject.RoomCategory,
	createdAt *valueobject.Timestamp,
	updatedAt *valueobject.Timestamp,
) (*Room, error) {
	if id == nil {
		return nil, validation.ErrRequiredRoomId
	}
	if adminId == nil {
		return nil, validation.ErrRequiredRoomAdminId
	}
	if name == nil {
		return nil, validation.ErrRequiredRoomName
	}
	if category == nil {
		return nil, validation.ErrRequiredRoomCategory
	}
	if createdAt == nil {
		return nil, validation.ErrRequiredRoomCreatedAt
	}
	if updatedAt == nil {
		return nil, validation.ErrRequiredRoomUpdatedAt
	}

	return &Room{
		id:        id,
		adminId:   adminId,
		name:      name,
		category:  category,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
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

func (r *Room) Id() *valueobject.Id {
	return r.id
}

func (r *Room) AdminId() *valueobject.UserId {
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
