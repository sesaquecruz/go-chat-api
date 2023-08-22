package valueobject

import "github.com/sesaquecruz/go-chat-api/internal/domain/validation"

type UserName struct {
	value string
}

func NewUserNameWith(value string) (*UserName, error) {
	if value == "" {
		return nil, validation.ErrRequiredUserName
	}
	if len(value) > 50 {
		return nil, validation.ErrSizeUserName
	}

	return &UserName{value: value}, nil
}

func (n *UserName) Value() string {
	return n.value
}
