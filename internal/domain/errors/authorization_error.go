package errors

type AuthorizationError struct {
	message string
}

func NewAuthorizationError(message string) *AuthorizationError {
	return &AuthorizationError{message: message}
}

func (e *AuthorizationError) Error() string {
	return e.message
}
