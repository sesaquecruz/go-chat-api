package validation

type InternalError struct {
	err error
}

func NewInternalError(err error) *InternalError {
	return &InternalError{
		err: err,
	}
}

func (e *InternalError) Error() string {
	return e.err.Error()
}
