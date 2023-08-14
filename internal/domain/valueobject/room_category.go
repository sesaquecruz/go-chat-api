package valueobject

import "github.com/sesaquecruz/go-chat-api/internal/domain/validation"

type RoomCategory struct {
	value string
}

func NewRoomCategoryWith(value string) (*RoomCategory, error) {
	if value == "" {
		return nil, validation.ErrRequiredRoomCategory
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
		return nil, validation.ErrInvalidRoomCategory
	}

	return &RoomCategory{value: value}, nil
}

func (c *RoomCategory) Value() string {
	return c.value
}
