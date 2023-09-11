package valueobject

import (
	"testing"
	"time"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"

	"github.com/stretchr/testify/assert"
)

func TestTimestamp_ShouldCreateATimestampWhenValueIsValid(t *testing.T) {
	timestamp1 := NewTimestamp()
	assert.NotNil(t, timestamp1)
	assert.Equal(t, timestamp1.Value(), timestamp1.Time().Format(timestampLayout))

	value := time.Now().UTC().Format(timestampLayout)
	timestamp2, err := NewTimestampWith(value)
	assert.NotNil(t, timestamp2)
	assert.Nil(t, err)
	assert.Equal(t, timestamp2.Value(), timestamp2.Time().Format(timestampLayout))

	timestamp2, err = NewTimestampWith(timestamp1.Value())
	assert.NotNil(t, timestamp2)
	assert.Nil(t, err)
	assert.Equal(t, timestamp1.Value(), timestamp2.Value())
	assert.True(t, timestamp1.Time().Equal(timestamp2.Time()))
}

func TestTimestamp_ShouldReturnAValidationErrorWhenValueIsInvalid(t *testing.T) {
	testCases := []struct {
		test  string
		value string
		err   error
	}{
		{
			"empty value",
			"",
			ErrRequiredTimestamp,
		},
		{
			"invalid value",
			"2006-01-02T15:04:00.999999",
			ErrInvalidTimestamp,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			timestamp, err := NewTimestampWith(tc.value)
			assert.Nil(t, timestamp)
			assert.NotNil(t, err)
			assert.ErrorIs(t, err, tc.err)
			assert.IsType(t, validation.ValidationError(""), err)
		})
	}
}
