package valueobject

import (
	"regexp"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
)

var userIdPattern = regexp.MustCompile(`^auth0|[a-fA-F0-9]{24}$`)

type UserId struct {
	value string
}

func NewUserIdWith(value string) (*UserId, error) {
	if value == "" {
		return nil, validation.ErrRequiredUserId
	}

	if !userIdPattern.MatchString(value) {
		return nil, validation.ErrInvalidUserId
	}

	return &UserId{value: value}, nil
}

func (id *UserId) Value() string {
	return id.value
}
