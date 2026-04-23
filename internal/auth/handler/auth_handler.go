package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"hris/internal/auth/dto"
	"hris/internal/auth/service"
	"hris/shared/constants"
	"hris/shared/helper"
	customValidator "hris/shared/validator"
)

var (
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrAccountDeactivated  = errors.New("account is deactivated")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if validationErrors := customValidator.ValidateStruct(&req); len(validationErrors) > 0 {
		return helper.ValidationErrorResponse(c, validationErrors)
	}

	result, err := h.service.Register(c.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			return helper.ErrorResponse(c, fiber.StatusConflict, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Registration failed", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusCreated, "Registration successful", result)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if validationErrors := customValidator.ValidateStruct(&req); len(validationErrors) > 0 {
		return helper.ValidationErrorResponse(c, validationErrors)
	}

	result, err := h.service.Login(c.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) ||
			errors.Is(err, service.ErrAccountDeactivated) {
			return helper.ErrorResponse(c, fiber.StatusUnauthorized, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Login failed", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Login successful", result)
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if validationErrors := customValidator.ValidateStruct(&req); len(validationErrors) > 0 {
		return helper.ValidationErrorResponse(c, validationErrors)
	}

	result, err := h.service.RefreshToken(c.Context(), &req)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusUnauthorized, "Token refresh failed", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Token refreshed successfully", result)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID := c.Locals(constants.ContextKeyUserID).(string)

	// Get refresh token ID from request body
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := h.service.Logout(c.Context(), userID, req.RefreshToken); err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Logout failed", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Logout successful", nil)
}

func (h *AuthHandler) LogoutAll(c *fiber.Ctx) error {
	userID := c.Locals(constants.ContextKeyUserID).(string)

	if err := h.service.LogoutAll(c.Context(), userID); err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Logout from all devices failed", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Logged out from all devices successfully", nil)
}

func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	userID := c.Locals(constants.ContextKeyUserID).(string)

	var req dto.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if validationErrors := customValidator.ValidateStruct(&req); len(validationErrors) > 0 {
		return helper.ValidationErrorResponse(c, validationErrors)
	}

	err := h.service.ChangePassword(c.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidOldPassword) {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}
		if errors.Is(err, service.ErrSamePassword) {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to change password", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Password changed successfully. Please login again.", nil)
}

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var req dto.ForgotPasswordRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if validationErrors := customValidator.ValidateStruct(&req); len(validationErrors) > 0 {
		return helper.ValidationErrorResponse(c, validationErrors)
	}

	token, err := h.service.ForgotPassword(c.Context(), req.Email)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) || errors.Is(err, service.ErrAccountDeactivated) {
			return helper.SuccessResponse(c, fiber.StatusOK, "If an account with that email exists, a reset token has been generated.", nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to process request", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Reset token generated successfully", fiber.Map{
		"token":   token,
		"message": "Use this token to reset your password. In production, this would be sent via email.",
	})
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var req dto.ResetPasswordRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if validationErrors := customValidator.ValidateStruct(&req); len(validationErrors) > 0 {
		return helper.ValidationErrorResponse(c, validationErrors)
	}

	err := h.service.ResetPassword(c.Context(), req.Token, req.NewPassword)
	if err != nil {
		if errors.Is(err, service.ErrInvalidResetToken) {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid or expired reset token", nil)
		}
		if errors.Is(err, service.ErrUserNotFound) {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "User not found", nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to reset password", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Password reset successfully. Please login with your new password.", nil)
}
