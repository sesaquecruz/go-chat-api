package valueobject

import (
	"strings"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
)

type MessageText struct {
	value string
}

func NewMessageTextWith(text string) (*MessageText, error) {
	value := strings.TrimSpace(text)

	if value == "" {
		return nil, validation.ErrRequiredMessageText
	}

	if len(value) > 100 {
		return nil, validation.ErrSizeMessageText
	}

	return &MessageText{value: value}, nil
}

func (t *MessageText) Value() string {
	return t.value
}
