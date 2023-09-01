package valueobject

import (
	"strings"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
)

const (
	ErrRequiredMessageText = validation.ValidationError("message text is required")
	ErrInvalidMessageText  = validation.ValidationError("message test length must be less than or equal to 100")
)

type MessageText struct {
	value string
}

func NewMessageTextWith(text string) (*MessageText, error) {
	value := strings.TrimSpace(text)

	if value == "" {
		return nil, ErrRequiredMessageText
	}

	if len(value) > 100 {
		return nil, ErrInvalidMessageText
	}

	return &MessageText{value: value}, nil
}

func (t *MessageText) Value() string {
	return t.value
}
