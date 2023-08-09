package valueobject

import "github.com/sesaquecruz/go-chat-api/internal/domain"

const ErrRequiredRoomCategory = "room category is required"
const ErrInvalidRoomCategory = "room category is invalid"

type RoomCategory struct {
	value string
}

func NewRoomCategoryWith(value string) (*RoomCategory, error) {
	if value == "" {
		return nil, domain.NewValidationError(ErrRequiredRoomCategory)
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
		return nil, domain.NewValidationError(ErrInvalidRoomCategory)
	}

	return &RoomCategory{value: value}, nil
}

func (rc *RoomCategory) Value() string {
	return rc.value
}
