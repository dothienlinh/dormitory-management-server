package response

import (
	"net/http"
)

// Response is the standard API response format
type Response struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    interface{}       `json:"data,omitempty"`
	Total   int64             `json:"total,omitempty"`
	Error   string            `json:"error,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

// SuccessResponse returns a success response
func SuccessResponse(data interface{}, total int64) Response {
	return Response{
		Success: true,
		Message: "Success",
		Data:    data,
		Total:   total,
	}
}

// ErrorResponse returns an error response
func ErrorResponse(message string, err string) Response {
	return Response{
		Success: false,
		Message: message,
		Error:   err,
	}
}

// StatusResponse contains HTTP status code and response
type StatusResponse struct {
	Status   int
	Response Response
}

// NewStatusResponse creates a new StatusResponse
func NewStatusResponse(status int, response Response) StatusResponse {
	return StatusResponse{
		Status:   status,
		Response: response,
	}
}

// Success returns a success status response
func Success(data interface{}, total int64) StatusResponse {
	return NewStatusResponse(
		http.StatusOK,
		SuccessResponse(data, total),
	)
}

// Created returns a created status response
func Created(data interface{}) StatusResponse {
	return NewStatusResponse(
		http.StatusCreated,
		SuccessResponse(data, 0),
	)
}

// BadRequest returns a bad request status response
func BadRequest(message string) StatusResponse {
	return NewStatusResponse(
		http.StatusBadRequest,
		ErrorResponse(message, "Bad Request"),
	)
}

// Unauthorized returns an unauthorized status response
func Unauthorized(message string) StatusResponse {
	return NewStatusResponse(
		http.StatusUnauthorized,
		ErrorResponse(message, "Unauthorized"),
	)
}

// Forbidden returns a forbidden status response
func Forbidden(message string) StatusResponse {
	return NewStatusResponse(
		http.StatusForbidden,
		ErrorResponse(message, "Forbidden"),
	)
}

// NotFound returns a not found status response
func NotFound(message string) StatusResponse {
	return NewStatusResponse(
		http.StatusNotFound,
		ErrorResponse(message, "Not Found"),
	)
}

// InternalServerError returns an internal server error status response
func InternalServerError(message string) StatusResponse {
	return NewStatusResponse(
		http.StatusInternalServerError,
		ErrorResponse(message, "Internal Server Error"),
	)
}

// ValidationError returns a validation error status response
func ValidationError(message string) StatusResponse {
	return NewStatusResponse(
		http.StatusUnprocessableEntity,
		ErrorResponse(message, "Validation Error"),
	)
}

// DBError returns a database error status response
func DBError(message string) StatusResponse {
	return NewStatusResponse(
		http.StatusInternalServerError,
		ErrorResponse(message, "Database Error"),
	)
}

// Validation returns a validation status response
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
