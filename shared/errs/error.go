package errs

import (
	"fmt"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common errors
func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:    400,
		Message: message,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:    401,
		Message: message,
	}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:    403,
		Message: message,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:    404,
		Message: message,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Code:    409,
		Message: message,
	}
}

func NewInternalServerError(message string, err error) *AppError {
	return &AppError{
		Code:    500,
		Message: message,
		Err:     err,
	}
}

func NewValidationError(message string) *AppError {
	return &AppError{
		Code:    422,
		Message: message,
	}
}