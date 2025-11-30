package helper

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
	Error      any    `json:"error,omitempty"`
}

func SuccessResponse(c *fiber.Ctx, statusCode int, message string, data any) error {
	return c.Status(statusCode).JSON(Response{
		Success:    true,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}

func ErrorResponse(c *fiber.Ctx, statusCode int, message string, err any) error {
	return c.Status(statusCode).JSON(Response{
		Success:    false,
		StatusCode: statusCode,
		Message:    message,
		Error:      err,
	})
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidationErrorResponse(c *fiber.Ctx, errors []ValidationError) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(Response{
		Success:    false,
		StatusCode: fiber.StatusUnprocessableEntity,
		Message:    "Validation failed",
		Error:      errors,
	})
}
