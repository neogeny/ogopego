package json

type JsonError struct {
	Message string
}

func (e *JsonError) Error() string {
	return e.Message
}

func newJsonError(msg string) *JsonError {
	return &JsonError{Message: msg}
}
