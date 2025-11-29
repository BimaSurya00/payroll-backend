package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/itsahyarr/go-fiber-boilerplate/shared/errs"
	"github.com/itsahyarr/go-fiber-boilerplate/shared/helper"
)

func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Log the error
		log.Printf("Error: %v", err)

		// Check if it's a custom app error
		if appErr, ok := err.(*errs.AppError); ok {
			return helper.ErrorResponse(c, appErr.Code, appErr.Message, appErr.Err)
		}

		// Check if it's a Fiber error
		if fiberErr, ok := err.(*fiber.Error); ok {
			return helper.ErrorResponse(c, fiberErr.Code, fiberErr.Message, nil)
		}

		// Default to internal server error
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Internal server error", err.Error())
	}
}
