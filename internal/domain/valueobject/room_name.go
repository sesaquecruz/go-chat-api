package valueobject

import "github.com/sesaquecruz/go-chat-api/internal/domain/validation"

type RoomName struct {
	value string
}

func NewRoomNameWith(value string) (*RoomName, error) {
	if value == "" {
		return nil, validation.ErrRequiredRoomName
	}
	if len(value) > 50 {
		return nil, validation.ErrMaxSizeRoomName
	}

	return &RoomName{value: value}, nil
}

func (n *RoomName) Value() string {
	return n.value
}
