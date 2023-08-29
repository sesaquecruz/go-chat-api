package entity

import (
	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
)

const ErrInvalidRoomAdmin = validation.UnauthorizedError("room admin is invalid")

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
) *Room {
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
) *Room {
	return &Room{
		id:        id,
		adminId:   adminId,
		name:      name,
		category:  category,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
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

func (r *Room) UpdateName(name *valueobject.RoomName) {
	r.name = name
	r.updatedAt = valueobject.NewTimestamp()
}

func (r *Room) UpdateCategory(category *valueobject.RoomCategory) {
	r.category = category
	r.updatedAt = valueobject.NewTimestamp()
}

func (r *Room) ValidateAdmin(adminId *valueobject.UserId) error {
	if r.adminId.Value() != adminId.Value() {
		return ErrInvalidRoomAdmin
	}

	return nil
}
