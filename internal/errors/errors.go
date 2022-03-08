package errors

type ExternalError struct {
	msg      string
	notFound bool
}

func (e *ExternalError) Error() string {
	return e.msg
}

func (e *ExternalError) External() bool {
	return true
}

func (e *ExternalError) NotFound() bool {
	return e.notFound
}

func IsExternal(err error) bool {
	external, ok := err.(*ExternalError)
	return ok && external.External()
}

func IsNotFound(err error) bool {
	external, ok := err.(*ExternalError)
	return ok && external.External() && external.NotFound()
}

func NewExternalError(message string) error {
	return &ExternalError{
		msg: message,
	}
}

func NewNotFoundError(message string) error {
	return &ExternalError{
		msg:      message,
		notFound: true,
	}
}
