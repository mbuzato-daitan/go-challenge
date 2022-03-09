package errors

type ExternalError struct {
	msg            string
	httpStatusCode int
}

func (e *ExternalError) Error() string {
	return e.msg
}

func (e *ExternalError) External() bool {
	return true
}

func (e *ExternalError) IsHTTP() bool {
	return e.httpStatusCode != 0
}

func (e *ExternalError) GetHTTPStatusCode() int {
	return e.httpStatusCode
}

func IsExternal(err error) bool {
	external, ok := err.(*ExternalError)
	return ok && external.External()
}

func IsHTTP(err error) bool {
	external, ok := err.(*ExternalError)
	return ok && external.External() && external.IsHTTP()
}

func GetHTTPStatusCode(err error) int {
	external, ok := err.(*ExternalError)
	if ok {
		return external.GetHTTPStatusCode()
	}

	return 0
}

func NewExternalError(message string) error {
	return &ExternalError{
		msg: message,
	}
}

func NewHTTPError(message string, code int) error {
	return &ExternalError{
		msg:            message,
		httpStatusCode: code,
	}
}
