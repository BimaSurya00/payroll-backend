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
