package valueobject

import "github.com/sesaquecruz/go-chat-api/internal/domain/errors"

const ErrRequiredRoomCategory = "room category is required"
const ErrInvalidRoomCategory = "room category is invalid"

type RoomCategory struct {
	value string
}

func NewRoomCategoryWith(value string) (*RoomCategory, error) {
	if value == "" {
		return nil, errors.NewValidationError(ErrRequiredRoomCategory)
	}

	switch value {
	case "General":
	case "Tech":
	case "Game":
	case "Book":
	case "Movie":
	case "Music":
	case "Language":
	case "Science":
	default:
		return nil, errors.NewValidationError(ErrInvalidRoomCategory)
	}

	return &RoomCategory{value: value}, nil
}

func (c *RoomCategory) Value() string {
	return c.value
}
