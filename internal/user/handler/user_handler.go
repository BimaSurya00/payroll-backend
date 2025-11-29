package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/user/dto"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/user/service"
	"github.com/itsahyarr/go-fiber-boilerplate/shared/constants"
	"github.com/itsahyarr/go-fiber-boilerplate/shared/helper"
	customValidator "github.com/itsahyarr/go-fiber-boilerplate/shared/validator"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req dto.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if validationErrors := customValidator.ValidateStruct(&req); len(validationErrors) > 0 {
		return helper.ValidationErrorResponse(c, validationErrors)
	}

	user, err := h.service.CreateUser(c.Context(), &req)
	if err != nil {
		if err.Error() == "email already exists" {
			return helper.ErrorResponse(c, fiber.StatusConflict, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create user", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusCreated, "User created successfully", user)
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", strconv.Itoa(constants.DefaultPage)))
	perPage, _ := strconv.Atoi(c.Query("per_page", strconv.Itoa(constants.DefaultPerPage)))

	if page < 1 {
		page = constants.DefaultPage
	}
	if perPage < 1 || perPage > constants.MaxPerPage {
		perPage = constants.DefaultPerPage
	}

	path := c.Protocol() + "://" + c.Hostname() + c.Path()

	pagination, err := h.service.GetUsers(c.Context(), page, perPage, path)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch users", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(pagination)
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.service.GetUserByID(c.Context(), id)
	if err != nil {
		if err.Error() == "user not found" {
			return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch user", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "User retrieved successfully", user)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if validationErrors := customValidator.ValidateStruct(&req); len(validationErrors) > 0 {
		return helper.ValidationErrorResponse(c, validationErrors)
	}

	if !req.HasUpdates() {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "No fields to update", nil)
	}

	user, err := h.service.UpdateUser(c.Context(), id, &req)
	if err != nil {
		if err.Error() == "user not found" {
			return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		if err.Error() == "email already exists" {
			return helper.ErrorResponse(c, fiber.StatusConflict, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update user", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "User updated successfully", user)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.service.DeleteUser(c.Context(), id); err != nil {
		if err.Error() == "user not found" {
			return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete user", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "User deleted successfully", nil)
}
