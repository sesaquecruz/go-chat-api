package valueobject

import (
	"time"

	"github.com/sesaquecruz/go-chat-api/internal/domain/errors"
)

const DateTimeLayout = "2006-01-02T15:04:05.999999Z"

const ErrRequiredDateTime = "datetime is required"
const ErrInvalidDateTime = "datetime is invalid"

type DateTime struct {
	stringValue string
	timeValue   *time.Time
}

func NewDateTime() *DateTime {
	stringValue := time.Now().UTC().Format(DateTimeLayout)
	timeValue, _ := time.Parse(DateTimeLayout, stringValue)

	return &DateTime{
		stringValue: stringValue,
		timeValue:   &timeValue,
	}
}

func NewDateTimeWith(value string) (*DateTime, error) {
	if value == "" {
		return nil, errors.NewValidationError(ErrRequiredDateTime)
	}

	t, err := time.Parse(DateTimeLayout, value)
	if err != nil {
		return nil, errors.NewValidationError(ErrInvalidDateTime)
	}

	timeValue := t.UTC()
	stringValue := timeValue.Format(DateTimeLayout)

	return &DateTime{
		stringValue: stringValue,
		timeValue:   &timeValue,
	}, nil
}

func (dt *DateTime) StringValue() string {
	return dt.stringValue
}

func (dt *DateTime) TimeValue() *time.Time {
	return dt.timeValue
}
