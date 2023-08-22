package valueobject

import (
	"github.com/google/uuid"
	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
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
		return nil, validation.ErrRequiredId
	}

	id, err := uuid.Parse(value)
	if err != nil || id == uuid.Nil {
		return nil, validation.ErrInvalidId
	}

	return &Id{value: id.String()}, nil
}

func (id *Id) Value() string {
	return id.value
}
