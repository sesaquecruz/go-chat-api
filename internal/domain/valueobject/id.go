package valueobject

import (
	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"

	"github.com/google/uuid"
)

const (
	ErrRequiredId = validation.ValidationError("id is required")
	ErrInvalidId  = validation.ValidationError("id is invalid")
)

type Id struct {
	value string
}

func NewId() *Id {
	id, _ := NewIdWith(uuid.New().String())
	return id
}

func NewIdWith(value string) (*Id, error) {
	if value == "" {
		return nil, ErrRequiredId
	}

	id, err := uuid.Parse(value)
	if err != nil || id == uuid.Nil {
		return nil, ErrInvalidId
	}

	return &Id{value: id.String()}, nil
}

func (id *Id) Value() string {
	return id.value
}
