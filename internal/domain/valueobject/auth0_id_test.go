package valueobject

import (
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestShouldCreateANewAuth0IDWhenValueIsValid(t *testing.T) {
	value := "auth0|64c8457bb160e37c8c34533b"
	id, err := NewAuth0IDWith(value)
	assert.NotNil(t, id)
	assert.Nil(t, err)
	assert.Equal(t, value, id.Value())
}

func TestShouldReturnARequiredAuth0IDErrorWhenValueIsInvalid(t *testing.T) {
	value := ""
	id, err := NewAuth0IDWith(value)
	assert.Nil(t, id)
	assert.NotNil(t, err)
	assert.IsType(t, &domain.ValidationError{}, err)
	assert.EqualError(t, err, ErrRequiredId)
}

func TestShouldReturnAnInvalidAuth0IDErrorWhenValueIsInvalid(t *testing.T) {
	value := "kj12389013kjfdsf9819023jkhfjds"
	id, err := NewAuth0IDWith(value)
	assert.Nil(t, id)
	assert.NotNil(t, err)
	assert.IsType(t, &domain.ValidationError{}, err)
	assert.EqualError(t, err, ErrInvalidId)
}
