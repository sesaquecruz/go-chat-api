package valueobject

import (
	"github.com/google/uuid"
	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
)

type ID struct {
	value uuid.UUID
}

func NewID() *ID {
	return &ID{value: uuid.New()}
}

func NewIDWith(value string) (*ID, error) {
	if value == "" {
		return nil, validation.ErrRequiredId
	}

	id, err := uuid.Parse(value)
	if err != nil || id == uuid.Nil {
		return nil, validation.ErrInvalidId
	}

	return &ID{value: id}, nil
}

func (id *ID) Value() string {
	return id.value.String()
}
