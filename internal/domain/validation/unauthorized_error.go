package validation

type UnauthorizedError string

func (e UnauthorizedError) Error() string {
	return string(e)
}
