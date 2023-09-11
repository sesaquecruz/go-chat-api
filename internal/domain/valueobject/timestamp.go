package valueobject

import (
	"time"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
)

const timestampLayout = "2006-01-02T15:04:05.999999Z"

const (
	ErrRequiredTimestamp = validation.ValidationError("timestamp is required")
	ErrInvalidTimestamp  = validation.ValidationError("timestamp is invalid")
)

type Timestamp struct {
	value time.Time
}

func NewTimestamp() *Timestamp {
	timestamp, _ := NewTimestampWith(time.Now().UTC().Format(timestampLayout))
	return timestamp
}

func NewTimestampWith(value string) (*Timestamp, error) {
	if value == "" {
		return nil, ErrRequiredTimestamp
	}

	t, err := time.Parse(timestampLayout, value)
	if err != nil {
		return nil, ErrInvalidTimestamp
	}

	t = t.UTC()

	return &Timestamp{value: t}, nil
}

func (t *Timestamp) Value() string {
	return t.value.Format(timestampLayout)
}

func (t *Timestamp) Time() time.Time {
	return t.value
}
