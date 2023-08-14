package valueobject

import (
	"testing"
	"time"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"

	"github.com/stretchr/testify/assert"
)

func TestTimestamp_ShouldCreateANewTimestamp(t *testing.T) {
	timestamp := NewTimestamp()
	assert.NotNil(t, timestamp)
	timeValue := timestamp.Time()
	stringValue := timestamp.Value()
	assert.Equal(t, stringValue, timeValue.Format(timestampLayout))
}

func TestTimestamp_ShouldCreateANewTimestampWhenValueIsValid(t *testing.T) {
	value := time.Now().UTC().Format(timestampLayout)
	timestamp, err := NewTimestampWith(value)
	assert.NotNil(t, timestamp)
	assert.Nil(t, err)
	timeValue := timestamp.Time()
	stringValue := timestamp.Value()
	assert.Equal(t, stringValue, timeValue.Format(timestampLayout))
}

func TestTimestamp_ShouldCreateEqualsTimestamps(t *testing.T) {
	timestamp1 := NewTimestamp()
	timestamp2, _ := NewTimestampWith(timestamp1.Value())
	assert.Equal(t, timestamp1.Value(), timestamp2.Value())
	assert.True(t, timestamp1.Time().Equal(timestamp2.Time()))
}

func TestTimestamp_ShouldReturnARequiredTimestampErrorWhenValueIsEmpty(t *testing.T) {
	value := ""
	timestamp, err := NewTimestampWith(value)
	assert.Nil(t, timestamp)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, validation.ErrRequiredTimestamp)
}

func TestTimestamp_ShouldReturnAInvalidTimestampErrorWhenValueIsInvalid(t *testing.T) {
	value := "2006-01-02T15:04:00.999999"
	timestamp, err := NewTimestampWith(value)
	assert.Nil(t, timestamp)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, validation.ErrInvalidTimestamp)
}
