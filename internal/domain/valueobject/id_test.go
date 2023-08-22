package valueobject

import (
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestId_ShouldCreateAnId(t *testing.T) {
	id := NewId()
	assert.NotNil(t, id)

	_, err := uuid.Parse(id.Value())
	assert.Nil(t, err)
}

func TestId_ShouldCreateAnIdWhenValueIsValid(t *testing.T) {
	value := uuid.New().String()
	id, err := NewIdWith(value)
	assert.NotNil(t, id)
	assert.Nil(t, err)
	assert.Equal(t, value, id.Value())
}

func TestId_ShouldReturnARequiredIdErrorWhenValueIsEmpty(t *testing.T) {
	value := ""
	id, err := NewIdWith(value)
	assert.Nil(t, id)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, validation.ErrRequiredId)
}

func TestId_ShouldReturnAnInvalidIdErrorWhenValueIsInvalid(t *testing.T) {
	value := "kj12389013kjfdsf9819023jkhfjds"
	id, err := NewIdWith(value)
	assert.Nil(t, id)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, validation.ErrInvalidId)
}
