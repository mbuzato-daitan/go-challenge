package errors

type ExternalError struct {
	Message string
}

func (e *ExternalError) Error() string {
	return e.Message
}

func (e *ExternalError) External() bool {
	return true
}

func IsExternal(err error) bool {
	external, ok := err.(*ExternalError)
	return ok && external.External()
}

func NewExternalError(message string) error {
	return &ExternalError{
		Message: message,
	}
}
