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

func TestUserId_ShouldReturnAValidationErrorWhenValueIsInvalid(t *testing.T) {
	testCases := []struct {
		test  string
		value string
		err   error
	}{
		{
			"empty value",
			"",
			ErrRequiredUserId,
		},
		{
			"invalid value",
			"kj12389013kjfdsf9819023jkhfjds",
			ErrInvalidUserId,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			id, err := NewUserIdWith(tc.value)
			assert.Nil(t, id)
			assert.NotNil(t, err)
			assert.ErrorIs(t, err, tc.err)
			assert.IsType(t, validation.ValidationError(""), err)
		})
	}
}
