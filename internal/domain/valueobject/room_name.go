package valueobject

import "github.com/sesaquecruz/go-chat-api/internal/domain/errors"

const ErrRequiredRoomName = "room name is required"
const ErrMaxSizeRoomName = "room name must not have more than 20 characters"

type RoomName struct {
	value string
}

func NewRoomNameWith(value string) (*RoomName, error) {
	if value == "" {
		return nil, errors.NewValidationError(ErrRequiredRoomName)
	}
	if len(value) > 50 {
		return nil, errors.NewValidationError(ErrMaxSizeRoomName)
	}

	return &RoomName{value: value}, nil
}

func (n *RoomName) Value() string {
	return n.value
}
