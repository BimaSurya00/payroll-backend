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

type PaginationMeta struct {
	CurrentPage  int    `json:"currentPage"`
	PerPage      int    `json:"perPage"`
	Total        int64  `json:"total"`
	LastPage     int    `json:"lastPage"`
	FirstPageUrl string `json:"firstPageUrl,omitempty"`
	LastPageUrl  string `json:"lastPageUrl,omitempty"`
	NextPageUrl  string `json:"nextPageUrl,omitempty"`
	PrevPageUrl  string `json:"prevPageUrl,omitempty"`
}

type ResponseWithPagination struct {
	Success    bool           `json:"success"`
	StatusCode int            `json:"statusCode"`
	Message    string         `json:"message"`
	Data       any            `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

func SuccessResponse(c *fiber.Ctx, statusCode int, message string, data any) error {
	return c.Status(statusCode).JSON(Response{
		Success:    true,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}

func SuccessResponseWithPagination(c *fiber.Ctx, statusCode int, message string, data any, pagination PaginationMeta) error {
	return c.Status(statusCode).JSON(ResponseWithPagination{
		Success:    true,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Pagination: pagination,
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
