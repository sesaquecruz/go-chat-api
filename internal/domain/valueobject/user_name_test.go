package valueobject

import (
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"

	"github.com/stretchr/testify/assert"
)

func TestUserName_ShouldCreateAnUserNameWhenValueIsValid(t *testing.T) {
	value := "An username"
	name, err := NewUserNameWith(value)
	assert.NotNil(t, name)
	assert.Nil(t, err)
	assert.Equal(t, value, name.Value())
}

func TestUserName_ShouldReturnARequiredUserNameErrorWhenValueIsEmpty(t *testing.T) {
	value := ""
	name, err := NewUserNameWith(value)
	assert.Nil(t, name)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, validation.ErrRequiredUserName)
}

func TestUserName_ShouldReturnASizeUserNameErrorWhenValueHasMoreThan50Characters(t *testing.T) {
	value := "dfaiuerewnvdiuoriewruuiwqeuqwe89123jladjsdasadiou23"
	assert.Equal(t, len(value), 51)

	name, err := NewUserNameWith(value)
	assert.Nil(t, name)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, validation.ErrSizeUserName)
}
