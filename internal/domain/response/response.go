package response

import (
	"net/http"
)

type Response struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    interface{}       `json:"data,omitempty"`
	Total   int64             `json:"total,omitempty"`
	Error   string            `json:"error,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func SuccessResponse(data interface{}, total int64) Response {
	return Response{
		Success: true,
		Message: "Success",
		Data:    data,
		Total:   total,
	}
}

func ErrorResponse(message string, err string) Response {
	return Response{
		Success: false,
		Message: message,
		Error:   err,
	}
}

type StatusResponse struct {
	Status   int
	Response Response
}

func NewStatusResponse(status int, response Response) StatusResponse {
	return StatusResponse{
		Status:   status,
		Response: response,
	}
}

func Success(data interface{}, total int64) StatusResponse {
	return NewStatusResponse(
		http.StatusOK,
		SuccessResponse(data, total),
	)
}

func Created(data interface{}) StatusResponse {
	return NewStatusResponse(
		http.StatusCreated,
		SuccessResponse(data, 0),
	)
}

func BadRequest(message string) StatusResponse {
	return NewStatusResponse(
		http.StatusBadRequest,
		ErrorResponse(message, "Bad Request"),
	)
}

func Unauthorized(message string) StatusResponse {
	return NewStatusResponse(
		http.StatusUnauthorized,
		ErrorResponse(message, "Unauthorized"),
	)
}

func Forbidden(message string) StatusResponse {
	return NewStatusResponse(
		http.StatusForbidden,
		ErrorResponse(message, "Forbidden"),
	)
}

func NotFound(message string) StatusResponse {
	return NewStatusResponse(
		http.StatusNotFound,
		ErrorResponse(message, "Not Found"),
	)
}

func InternalServerError(message string) StatusResponse {
	return NewStatusResponse(
		http.StatusInternalServerError,
		ErrorResponse(message, "Internal Server Error"),
	)
}

func ValidationError(message string) StatusResponse {
	return NewStatusResponse(
		http.StatusUnprocessableEntity,
		ErrorResponse(message, "Validation Error"),
	)
}

func DBError(message string) StatusResponse {
	return NewStatusResponse(
		http.StatusInternalServerError,
		ErrorResponse(message, "Database Error"),
	)
}

func Validation(message string, errors map[string]string) StatusResponse {
	return NewStatusResponse(
		http.StatusBadRequest,
		Response{
			Success: false,
			Message: message,
			Errors:  errors,
		},
	)
}
