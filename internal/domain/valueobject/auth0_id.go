package valueobject

import (
	"regexp"

	"github.com/sesaquecruz/go-chat-api/internal/domain"
)

var auth0IDPattern = regexp.MustCompile(`^auth0|[a-fA-F0-9]{24}$`)

type Auth0ID struct {
	value string
}

func NewAuth0IDWith(value string) (*Auth0ID, error) {
	if value == "" {
		return nil, domain.NewValidationError(ErrRequiredId)
	}

	if !auth0IDPattern.MatchString(value) {
		return nil, domain.NewValidationError(ErrInvalidId)
	}

	return &Auth0ID{value: value}, nil
}

func (id *Auth0ID) Value() string {
	return id.value
}
