package valueobject

import "github.com/sesaquecruz/go-chat-api/internal/domain"

const ErrRequiredRoomName = "room name is required"
const ErrMaxSizeRoomName = "room name must not have more than 20 characters"

type RoomName struct {
	value string
}

func NewRoomNameWith(value string) (*RoomName, error) {
	if value == "" {
		return nil, domain.NewValidationError(ErrRequiredRoomName)
	}
	if len(value) > 50 {
		return nil, domain.NewValidationError(ErrMaxSizeRoomName)
	}

	return &RoomName{value: value}, nil
}

func (rn *RoomName) Value() string {
	return rn.value
}
