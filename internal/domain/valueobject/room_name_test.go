package valueobject

import (
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"

	"github.com/stretchr/testify/assert"
)

func TestRoomName_ShouldCreateANewRoomNameWhenValueIsValid(t *testing.T) {
	value := "A Room Name"
	name, err := NewRoomNameWith(value)
	assert.NotNil(t, name)
	assert.Nil(t, err)
	assert.Equal(t, value, name.Value())
}

func TestRoomName_ShouldReturnARequiredRoomNameErrorWhenValueIsEmpty(t *testing.T) {
	value := ""
	name, err := NewRoomNameWith(value)
	assert.Nil(t, name)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, validation.ErrRequiredRoomName)
}

func TestRoomName_ShouldReturnAMaxSizeRoomNameErrorWhenValueHasMoreThan20Characters(t *testing.T) {
	value := "dfaiuerewnvdiuoriewruuiwqeuqwe89123jladjsdasadiou23"
	assert.Equal(t, len(value), 51)
	name, err := NewRoomNameWith(value)
	assert.Nil(t, name)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, validation.ErrMaxSizeRoomName)
}
