package valueobject

import (
	"testing"
	"time"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"

	"github.com/stretchr/testify/assert"
)

func TestTimestamp_ShouldCreateATimestamp(t *testing.T) {
	timestamp := NewTimestamp()
	assert.NotNil(t, timestamp)

	timeValue := timestamp.Value()
	stringValue := timestamp.String()
	assert.Equal(t, stringValue, timeValue.Format(timestampLayout))
}

func TestTimestamp_ShouldCreateATimestampWhenValueIsValid(t *testing.T) {
	value := time.Now().UTC().Format(timestampLayout)
	timestamp, err := NewTimestampWith(value)
	assert.NotNil(t, timestamp)
	assert.Nil(t, err)

	timeValue := timestamp.Value()
	stringValue := timestamp.String()
	assert.Equal(t, stringValue, timeValue.Format(timestampLayout))
}

func TestTimestamp_ShouldCreateEqualsTimestamps(t *testing.T) {
	timestamp1 := NewTimestamp()
	timestamp2, err := NewTimestampWith(timestamp1.String())
	assert.NotNil(t, timestamp1)
	assert.Nil(t, err)
	assert.True(t, timestamp1.Value().Equal(timestamp2.Value()))
	assert.Equal(t, timestamp1.String(), timestamp2.String())
}

func TestTimestamp_ShouldReturnARequiredTimestampErrorWhenValueIsEmpty(t *testing.T) {
	value := ""
	timestamp, err := NewTimestampWith(value)
	assert.Nil(t, timestamp)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, validation.ErrRequiredTimestamp)
}

func TestTimestamp_ShouldReturnAnInvalidTimestampErrorWhenValueIsInvalid(t *testing.T) {
	value := "2006-01-02T15:04:00.999999"
	timestamp, err := NewTimestampWith(value)
	assert.Nil(t, timestamp)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, validation.ErrInvalidTimestamp)
}
