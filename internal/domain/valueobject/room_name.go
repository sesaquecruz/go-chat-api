package valueobject

import "github.com/sesaquecruz/go-chat-api/internal/domain/validation"

const (
	ErrRequiredRoomName = validation.ValidationError("room name is required")
	ErrInvalidRoomName  = validation.ValidationError("room name must not have more than 50 characters")
)

type RoomName struct {
	value string
}

func NewRoomNameWith(value string) (*RoomName, error) {
	if value == "" {
		return nil, ErrRequiredRoomName
	}
	if len(value) > 50 {
		return nil, ErrInvalidRoomName
	}

	return &RoomName{value: value}, nil
}

func (n *RoomName) Value() string {
	return n.value
}
