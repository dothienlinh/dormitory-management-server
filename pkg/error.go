package pkg

type HttpError struct {
	Message    string
	Errors     any
	StatusCode int
	Code       int
}

func (e *HttpError) Error() string {
	return e.Message
}

func NewHttpError(statusCode int, msg string) *HttpError {
	return &HttpError{Message: msg, StatusCode: statusCode}
}

func NewCustomHttpError(statusCode int, msg string, errors any) *HttpError {
	return &HttpError{Message: msg, StatusCode: statusCode, Errors: errors}
}

func NewHttpCodeError(statusCode, code int, msg string) *HttpError {
	return &HttpError{StatusCode: statusCode, Code: code, Message: msg}
}
