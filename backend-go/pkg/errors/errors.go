package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error with HTTP status code
type AppError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	Err        error  `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// Common error constructors

func NewBadRequestError(message string, err error) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Err:        err,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

func NewRateLimitError(message string) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: http.StatusTooManyRequests,
	}
}
