package valueobject

import (
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestShouldCreateANewID(t *testing.T) {
	id := NewID()
	assert.NotNil(t, id)
	_, err := uuid.Parse(id.Value())
	assert.Nil(t, err)
}

func TestShouldCreateANewIDWhenValueIsValid(t *testing.T) {
	value := uuid.New().String()
	id, err := NewIDWith(value)
	assert.NotNil(t, id)
	assert.Nil(t, err)
	assert.Equal(t, value, id.Value())
}

func TestShouldReturnARequiredIDErrorWhenValueIsInvalid(t *testing.T) {
	value := ""
	id, err := NewIDWith(value)
	assert.Nil(t, id)
	assert.NotNil(t, err)
	assert.IsType(t, &domain.ValidationError{}, err)
	assert.EqualError(t, err, ErrRequiredId)
}

func TestShouldReturnAnInvalidIDErrorWhenValueIsInvalid(t *testing.T) {
	value := "kj12389013kjfdsf9819023jkhfjds"
	id, err := NewIDWith(value)
	assert.Nil(t, id)
	assert.NotNil(t, err)
	assert.IsType(t, &domain.ValidationError{}, err)
	assert.EqualError(t, err, ErrInvalidId)
}