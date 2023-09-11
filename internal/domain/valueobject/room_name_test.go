package valueobject

import (
	"testing"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"

	"github.com/stretchr/testify/assert"
)

func TestRoomName_ShouldCreateARoomNameWhenValueIsValid(t *testing.T) {
	value := "A Room Name"
	roomName, err := NewRoomNameWith(value)
	assert.NotNil(t, roomName)
	assert.Nil(t, err)
	assert.Equal(t, value, roomName.Value())
}

func TestRoomName_ShouldReturnAValidationErrorWhenValueIsInvalid(t *testing.T) {
	testCases := []struct {
		test  string
		value string
		err   error
	}{
		{
			"empty value",
			"",
			ErrRequiredRoomName,
		},
		{
			"invalid value size",
			"dfaiuerewnvdiuoriewruuiwqeuqwe89123jladjsdasadiou23",
			ErrInvalidRoomName,
		},
	}

	assert.Equal(t, len(testCases[1].value), 51)

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			roomName, err := NewRoomNameWith(tc.value)
			assert.Nil(t, roomName)
			assert.NotNil(t, err)
			assert.ErrorIs(t, err, tc.err)
			assert.IsType(t, validation.ValidationError(""), err)
		})
	}
}
