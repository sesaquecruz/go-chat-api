package valueobject

import (
	"regexp"

	"github.com/sesaquecruz/go-chat-api/internal/domain/errors"
)

var auth0IDPattern = regexp.MustCompile(`^auth0|[a-fA-F0-9]{24}$`)

type Auth0ID struct {
	value string
}

func NewAuth0IDWith(value string) (*Auth0ID, error) {
	if value == "" {
		return nil, errors.NewValidationError(ErrRequiredId)
	}

	if !auth0IDPattern.MatchString(value) {
		return nil, errors.NewValidationError(ErrInvalidId)
	}

	return &Auth0ID{value: value}, nil
}

func (id *Auth0ID) Value() string {
	return id.value
}
