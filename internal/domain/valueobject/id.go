package valueobject

import (
	"github.com/sesaquecruz/go-chat-api/internal/domain/errors"

	"github.com/google/uuid"
)

const ErrRequiredId = "id is required"
const ErrInvalidId = "id is invalid"

type ID struct {
	value string
}

func NewID() *ID {
	return &ID{value: uuid.NewString()}
}

func NewIDWith(value string) (*ID, error) {
	if value == "" {
		return nil, errors.NewValidationError(ErrRequiredId)
	}

	id, err := uuid.Parse(value)
	if id == uuid.Nil || err != nil {
		return nil, errors.NewValidationError(ErrInvalidId)
	}

	return &ID{value: id.String()}, nil
}

func (id *ID) Value() string {
	return id.value
}
