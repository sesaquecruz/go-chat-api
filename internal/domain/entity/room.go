package entity

import (
	"github.com/sesaquecruz/go-chat-api/internal/domain"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
)

const ErrRequiredRoomId = "room id is required"
const ErrRequiredRoomAdminId = "room admin id is required"
const ErrRequiredRoomName = "room name is required"
const ErrRequiredRoomCategory = "room category is required"
const ErrRequiredRoomCreatedAt = "room created at is required"
const ErrRequiredRoomUpdatedAt = "room updated at is required"

type Room struct {
	id        *valueobject.ID
	adminId   *valueobject.Auth0ID
	name      *valueobject.RoomName
	category  *valueobject.RoomCategory
	createdAt *valueobject.DateTime
	updatedAt *valueobject.DateTime
}

func NewRoom(
	adminId *valueobject.Auth0ID,
	name *valueobject.RoomName,
	category *valueobject.RoomCategory,
) (*Room, error) {
	now := valueobject.NewDateTime()
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
	createdAt *valueobject.DateTime,
	updatedAt *valueobject.DateTime,
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
		return domain.NewValidationError(ErrRequiredRoomId)
	}
	if r.adminId == nil {
		return domain.NewValidationError(ErrRequiredRoomAdminId)
	}
	if r.name == nil {
		return domain.NewValidationError(ErrRequiredRoomName)
	}
	if r.category == nil {
		return domain.NewValidationError(ErrRequiredRoomCategory)
	}
	if r.createdAt == nil {
		return domain.NewValidationError(ErrRequiredRoomCreatedAt)
	}
	if r.updatedAt == nil {
		return domain.NewValidationError(ErrRequiredRoomUpdatedAt)
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

func (r *Room) CreatedAt() *valueobject.DateTime {
	return r.createdAt
}

func (r *Room) UpdatedAt() *valueobject.DateTime {
	return r.updatedAt
}

func (r *Room) UpdateName(name *valueobject.RoomName) error {
	if name == nil {
		return domain.NewValidationError(ErrRequiredRoomName)
	}

	r.name = name
	r.updatedAt = valueobject.NewDateTime()
	return nil
}

func (r *Room) UpdateCategory(category *valueobject.RoomCategory) error {
	if category == nil {
		return domain.NewValidationError(ErrRequiredRoomCategory)
	}

	r.category = category
	r.updatedAt = valueobject.NewDateTime()
	return nil
}
