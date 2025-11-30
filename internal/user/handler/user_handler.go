package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/minio"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/user/dto"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/user/service"
	"github.com/itsahyarr/go-fiber-boilerplate/shared/constants"
	"github.com/itsahyarr/go-fiber-boilerplate/shared/helper"
	customValidator "github.com/itsahyarr/go-fiber-boilerplate/shared/validator"
)

type UserHandler struct {
	service      service.UserService
	minioService minio.MinioService
}

func NewUserHandler(service service.UserService, minioService minio.MinioService) *UserHandler {
	return &UserHandler{
		service:      service,
		minioService: minioService,
	}
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
		if errors.Is(err, service.ErrEmailAlreadyExists) {
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
		if errors.Is(err, service.ErrUserNotFound) {
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
		if errors.Is(err, service.ErrUserNotFound) {
			return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			return helper.ErrorResponse(c, fiber.StatusConflict, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update user", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "User updated successfully", user)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.service.DeleteUser(c.Context(), id); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete user", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "User deleted successfully", nil)
}

// GetOwnProfile returns the authenticated user's profile
func (h *UserHandler) GetOwnProfile(c *fiber.Ctx) error {
	userID := c.Locals(constants.ContextKeyUserID).(string)

	user, err := h.service.GetUserByID(c.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch profile", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Profile retrieved successfully", user)
}

// UploadProfileImage uploads a profile image for a user
func (h *UserHandler) UploadProfileImage(c *fiber.Ctx) error {
	userID := c.Params("id")

	// Get current user from context (for authorization check)
	currentUserID := c.Locals(constants.ContextKeyUserID).(string)
	currentUserRole := c.Locals(constants.ContextKeyUserRole).(string)

	// Authorization: USER can only upload their own image
	if currentUserRole == constants.RoleUser && currentUserID != userID {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "You can only upload your own profile image", nil)
	}

	// Parse multipart form
	file, err := c.FormFile("image")
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Image file is required", err.Error())
	}

	// Open the file
	fileHandle, err := file.Open()
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to open file", err.Error())
	}
	defer fileHandle.Close()

	// Upload image
	fileURL, err := h.minioService.UploadUserImage(c.Context(), userID, fileHandle, file)
	if err != nil {
		// Check if it's a validation error
		if validationErr, ok := err.(*minio.ValidationError); ok {
			return helper.ErrorResponse(c, fiber.StatusUnprocessableEntity, validationErr.Message, nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to upload image", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Profile image uploaded successfully", map[string]string{
		"imageUrl": fileURL,
	})
}

// UpdateProfileImage updates/replaces a user's profile image
func (h *UserHandler) UpdateProfileImage(c *fiber.Ctx) error {
	userID := c.Params("id")

	// Get current user from context (for authorization check)
	currentUserID := c.Locals(constants.ContextKeyUserID).(string)
	currentUserRole := c.Locals(constants.ContextKeyUserRole).(string)

	// Authorization: USER can only update their own image
	if currentUserRole == constants.RoleUser && currentUserID != userID {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "You can only update your own profile image", nil)
	}

	// Parse multipart form
	file, err := c.FormFile("image")
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Image file is required", err.Error())
	}

	// Open the file
	fileHandle, err := file.Open()
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to open file", err.Error())
	}
	defer fileHandle.Close()

	// Update image
	fileURL, err := h.minioService.UpdateUserImage(c.Context(), userID, fileHandle, file)
	if err != nil {
		// Check if it's a validation error
		if validationErr, ok := err.(*minio.ValidationError); ok {
			return helper.ErrorResponse(c, fiber.StatusUnprocessableEntity, validationErr.Message, nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update image", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Profile image updated successfully", map[string]string{
		"imageUrl": fileURL,
	})
}

// DeleteProfileImage deletes a user's profile image
func (h *UserHandler) DeleteProfileImage(c *fiber.Ctx) error {
	userID := c.Params("id")

	// Get current user from context (for authorization check)
	currentUserID := c.Locals(constants.ContextKeyUserID).(string)
	currentUserRole := c.Locals(constants.ContextKeyUserRole).(string)

	// Authorization: USER can only delete their own image
	if currentUserRole == constants.RoleUser && currentUserID != userID {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "You can only delete your own profile image", nil)
	}

	// Delete image
	if err := h.minioService.DeleteUserImage(c.Context(), userID); err != nil {
		if errors.Is(err, minio.ErrUserHasNoProfileImageToDelete) {
			return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete image", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Profile image deleted successfully", nil)
}
