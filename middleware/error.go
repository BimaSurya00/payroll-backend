package middleware

import (
	"github.com/gofiber/fiber/v2"
	"hris/shared/errs"
	"hris/shared/helper"
	"go.uber.org/zap"
)

func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Log the error
		zap.L().Error("unhandled error",
			zap.Error(err),
			zap.String("path", c.Path()),
			zap.String("method", c.Method()),
		)

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
