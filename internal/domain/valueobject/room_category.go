package valueobject

import "github.com/sesaquecruz/go-chat-api/internal/domain/validation"

const (
	ErrRequiredRoomCategory = validation.ValidationError("room category is required")
	ErrInvalidRoomCategory  = validation.ValidationError("room category is invalid")
)

type RoomCategory struct {
	value string
}

func NewRoomCategoryWith(value string) (*RoomCategory, error) {
	if value == "" {
		return nil, ErrRequiredRoomCategory
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
		return nil, ErrInvalidRoomCategory
	}

	return &RoomCategory{value: value}, nil
}

func (c *RoomCategory) Value() string {
	return c.value
}
