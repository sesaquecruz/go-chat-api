package domain

type ValidationError struct {
	message string
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{message: message}
}

func (ve *ValidationError) Error() string {
	return ve.message
}

type GatewayError struct {
	message string
}

func NewGatewayError(message string) *GatewayError {
	return &GatewayError{message: message}
}

func (ge *GatewayError) Error() string {
	return ge.message
}
