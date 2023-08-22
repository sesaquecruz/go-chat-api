package valueobject

import (
	"time"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
)

const timestampLayout = "2006-01-02T15:04:05.999999Z"

type Timestamp struct {
	value time.Time
}

func NewTimestamp() *Timestamp {
	timestamp, _ := NewTimestampWith(time.Now().UTC().Format(timestampLayout))
	return timestamp
}

func NewTimestampWith(value string) (*Timestamp, error) {
	if value == "" {
		return nil, validation.ErrRequiredTimestamp
	}

	t, err := time.Parse(timestampLayout, value)
	if err != nil {
		return nil, validation.ErrInvalidTimestamp
	}

	t = t.UTC()

	return &Timestamp{value: t}, nil
}

func (t *Timestamp) Value() time.Time {
	return t.value
}

func (t *Timestamp) String() string {
	return t.value.Format(timestampLayout)
}
