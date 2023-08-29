package valueobject

import "github.com/sesaquecruz/go-chat-api/internal/domain/validation"

const (
	ErrRequiredUserName = validation.ValidationError("user name is required")
	ErrInvalidUserName  = validation.ValidationError("user name must not have more than 50 characters")
)

type UserName struct {
	value string
}

func NewUserNameWith(value string) (*UserName, error) {
	if value == "" {
		return nil, ErrRequiredUserName
	}
	if len(value) > 50 {
		return nil, ErrInvalidUserName
	}

	return &UserName{value: value}, nil
}

func (n *UserName) Value() string {
	return n.value
}
