package valueobject

import (
	"testing"
	"time"

	"github.com/sesaquecruz/go-chat-api/internal/domain/errors"

	"github.com/stretchr/testify/assert"
)

func TestDateTime_ShouldCreateANewDateTime(t *testing.T) {
	datetime := NewDateTime()
	assert.NotNil(t, datetime)
	timeValue := datetime.TimeValue()
	assert.NotNil(t, timeValue)
	parsedValue, err := time.Parse(DateTimeLayout, datetime.StringValue())
	assert.NotNil(t, parsedValue)
	assert.Nil(t, err)
	assert.True(t, timeValue.Equal(parsedValue))
}

func TestDateTime_ShouldCreateANewDateTimeWhenValueIsValid(t *testing.T) {
	value := time.Now().UTC().Format(DateTimeLayout)
	datetime, err := NewDateTimeWith(value)
	assert.NotNil(t, datetime)
	assert.Nil(t, err)
	timeValue := datetime.TimeValue()
	assert.NotNil(t, timeValue)
	parsedValue, err := time.Parse(DateTimeLayout, datetime.StringValue())
	assert.NotNil(t, parsedValue)
	assert.Nil(t, err)
	assert.True(t, timeValue.Equal(parsedValue))
}

func TestDateTime_ShouldCreateEqualsDateTime(t *testing.T) {
	datetime1 := NewDateTime()
	assert.NotNil(t, datetime1)
	datetime2, err := NewDateTimeWith(datetime1.StringValue())
	assert.NotNil(t, datetime2)
	assert.Nil(t, err)
	assert.Equal(t, datetime1.StringValue(), datetime2.StringValue())
	assert.True(t, datetime1.TimeValue().Equal(*datetime2.TimeValue()))
}

func TestDateTime_ShouldReturnARequiredDateTimeErrorWhenValueIsEmpty(t *testing.T) {
	value := ""
	datetime, err := NewDateTimeWith(value)
	assert.Nil(t, datetime)
	assert.NotNil(t, err)
	assert.IsType(t, &errors.ValidationError{}, err)
	assert.EqualError(t, err, ErrRequiredDateTime)
}

func TestDateTime_ShouldReturnAInvalidDateTimeErrorWhenValueIsInvalid(t *testing.T) {
	value := "2006-01-02T15:04:00.999999"
	datetime, err := NewDateTimeWith(value)
	assert.Nil(t, datetime)
	assert.NotNil(t, err)
	assert.IsType(t, &errors.ValidationError{}, err)
	assert.EqualError(t, err, ErrInvalidDateTime)
}
