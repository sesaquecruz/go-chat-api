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

func TestUserName_ShouldReturnAValidationErrorWhenValueIsInvalid(t *testing.T) {
	testCases := []struct {
		test  string
		value string
		err   error
	}{
		{
			"empty value",
			"",
			ErrRequiredUserName,
		},
		{
			"invalid value size",
			"dfaiuerewnvdiuoriewruuiwqeuqwe89123jladjsdasadiou23",
			ErrInvalidUserName,
		},
	}

	assert.Equal(t, len(testCases[1].value), 51)

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			userName, err := NewUserNameWith(tc.value)
			assert.Nil(t, userName)
			assert.NotNil(t, err)
			assert.ErrorIs(t, err, tc.err)
			assert.IsType(t, validation.ValidationError(""), err)
		})
	}
}
