package valueobject

import (
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"

	"github.com/stretchr/testify/assert"
)

func TestUserId_ShouldCreateAnUserIdWhenValueIsValid(t *testing.T) {
	value := "auth0|64c8457bb160e37c8c34533b"
	id, err := NewUserIdWith(value)
	assert.NotNil(t, id)
	assert.Nil(t, err)
	assert.Equal(t, value, id.Value())
}

func TestUserId_ShouldReturnARequiredUserIdErrorWhenValueIsEmpty(t *testing.T) {
	value := ""
	id, err := NewUserIdWith(value)
	assert.Nil(t, id)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, validation.ErrRequiredUserId)
}

func TestUserId_ShouldReturnAnInvalidUserIdErrorWhenValueIsInvalid(t *testing.T) {
	value := "kj12389013kjfdsf9819023jkhfjds"
	id, err := NewUserIdWith(value)
	assert.Nil(t, id)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, validation.ErrInvalidUserId)
}
