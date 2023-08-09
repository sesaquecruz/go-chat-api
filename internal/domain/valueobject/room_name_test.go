package valueobject

import (
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestShouldCreateANewRoomNameWhenValueIsValid(t *testing.T) {
	value := "A Room Name"
	name, err := NewRoomNameWith(value)
	assert.NotNil(t, name)
	assert.Nil(t, err)
	assert.Equal(t, value, name.Value())
}

func TestShouldReturnARequiredRoomNameErrorWhenValueIsEmpty(t *testing.T) {
	value := ""
	name, err := NewRoomNameWith(value)
	assert.Nil(t, name)
	assert.NotNil(t, err)
	assert.IsType(t, &domain.ValidationError{}, err)
	assert.EqualError(t, err, ErrRequiredRoomName)
}

func TestShouldReturnAMaxSizeRoomNameErrorWhenValueHasMoreThan20Characters(t *testing.T) {
	value := "dfaiuerewnvdiuoriewruuiwqeuqwe89123jladjsdasadiou23"
	assert.Equal(t, len(value), 51)
	name, err := NewRoomNameWith(value)
	assert.Nil(t, name)
	assert.NotNil(t, err)
	assert.IsType(t, &domain.ValidationError{}, err)
	assert.EqualError(t, err, ErrMaxSizeRoomName)
}
