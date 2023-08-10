package errors

type GatewayError struct {
	message string
}

func NewGatewayError(message string) *GatewayError {
	return &GatewayError{message: message}
}

func (e *GatewayError) Error() string {
	return e.message
}
