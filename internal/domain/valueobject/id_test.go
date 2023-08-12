package valueobject

import (
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestID_ShouldCreateANewID(t *testing.T) {
	id := NewID()
	assert.NotNil(t, id)
	_, err := uuid.Parse(id.Value())
	assert.Nil(t, err)
}

func TestID_ShouldCreateANewIDWhenValueIsValid(t *testing.T) {
	value := uuid.New().String()
	id, err := NewIDWith(value)
	assert.NotNil(t, id)
	assert.Nil(t, err)
	assert.Equal(t, value, id.Value())
}

func TestID_ShouldReturnARequiredIDErrorWhenValueIsInvalid(t *testing.T) {
	value := ""
	id, err := NewIDWith(value)
	assert.Nil(t, id)
	assert.NotNil(t, err)
	assert.IsType(t, &errors.ValidationError{}, err)
	assert.EqualError(t, err, ErrRequiredId)
}

func TestID_ShouldReturnAnInvalidIDErrorWhenValueIsInvalid(t *testing.T) {
	value := "kj12389013kjfdsf9819023jkhfjds"
	id, err := NewIDWith(value)
	assert.Nil(t, id)
	assert.NotNil(t, err)
	assert.IsType(t, &errors.ValidationError{}, err)
	assert.EqualError(t, err, ErrInvalidId)
}
