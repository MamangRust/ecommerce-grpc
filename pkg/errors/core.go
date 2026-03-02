package errors

import (
	"fmt"
	"net/http"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type AppError struct {
	Code        int               `json:"-"`
	Message     string            `json:"message"`
	Retryable   bool              `json:"retryable,omitempty"`
	Validations []ValidationError `json:"validations,omitempty"`
	Internal    error             `json:"-"`
}

func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Internal)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Internal
}

func (e *AppError) WithInternal(err error) *AppError {
	return &AppError{
		Code:        e.Code,
		Message:     e.Message,
		Retryable:   e.Retryable,
		Validations: e.Validations,
		Internal:    err,
	}
}

func (e *AppError) WithMessage(message string) *AppError {
	return &AppError{
		Code:        e.Code,
		Message:     message,
		Retryable:   e.Retryable,
		Validations: e.Validations,
		Internal:    e.Internal,
	}
}

func (e *AppError) WithValidations(validations []ValidationError) *AppError {
	return &AppError{
		Code:        e.Code,
		Message:     e.Message,
		Retryable:   e.Retryable,
		Validations: validations,
		Internal:    e.Internal,
	}
}

func (e *AppError) AsRetryable() *AppError {
	return &AppError{
		Code:        e.Code,
		Message:     e.Message,
		Retryable:   true,
		Validations: e.Validations,
		Internal:    e.Internal,
	}
}

var (
	ErrBadRequest = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Bad request",
	}

	ErrValidationFailed = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Validation failed",
	}

	ErrUnauthorized = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "Unauthorized",
	}

	ErrForbidden = &AppError{
		Code:    http.StatusForbidden,
		Message: "Forbidden",
	}

	ErrNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "Resource not found",
	}

	ErrConflict = &AppError{
		Code:    http.StatusConflict,
		Message: "Resource conflict",
	}

	ErrTooManyRequests = &AppError{
		Code:      http.StatusTooManyRequests,
		Message:   "Too many requests",
		Retryable: true,
	}

	ErrInternal = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "Internal server error",
	}

	ErrServiceUnavailable = &AppError{
		Code:      http.StatusServiceUnavailable,
		Message:   "Service unavailable",
		Retryable: true,
	}

	ErrTimeout = &AppError{
		Code:      http.StatusGatewayTimeout,
		Message:   "Request timeout",
		Retryable: true,
	}
)
