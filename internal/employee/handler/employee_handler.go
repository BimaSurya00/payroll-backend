package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"hris/internal/employee/dto"
	"hris/internal/employee/service"
	"hris/shared/constants"
	"hris/shared/helper"
	customValidator "hris/shared/validator"
)

type EmployeeHandler struct {
	service service.EmployeeService
}

func NewEmployeeHandler(service service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		service: service,
	}
}

func (h *EmployeeHandler) CreateEmployee(c *fiber.Ctx) error {
	var req dto.CreateEmployeeRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Debug: Log received phone number
	println("DEBUG: Received phoneNumber = ", req.PhoneNumber)

	// Validate request
	if validationErrors := customValidator.ValidateStruct(&req); len(validationErrors) > 0 {
		return helper.ValidationErrorResponse(c, validationErrors)
	}

	// Get company_id from context (set by JWTAuth middleware)
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	employee, err := h.service.CreateEmployee(c.Context(), &req, companyID)
	if err != nil {
		if errors.Is(err, service.ErrUserCreationFailed) || errors.Is(err, service.ErrEmployeeCreateFailed) {
			return helper.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create employee", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusCreated, "Employee created successfully", employee)
}

func (h *EmployeeHandler) GetAllEmployees(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", strconv.Itoa(constants.DefaultPage)))
	perPage, _ := strconv.Atoi(c.Query("per_page", strconv.Itoa(constants.DefaultPerPage)))
	search := c.Query("search", "")

	if page < 1 {
		page = constants.DefaultPage
	}
	if perPage < 1 || perPage > constants.MaxPerPage {
		perPage = constants.DefaultPerPage
	}

	path := c.Protocol() + "://" + c.Hostname() + c.Path()

	// Get company_id from context
	companyID, _ := c.Locals(constants.ContextKeyCompanyID).(string)

	pagination, err := h.service.GetAllEmployees(c.Context(), page, perPage, path, search, companyID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch employees", err.Error())
	}

	// Extract pagination metadata
	paginationMeta := helper.PaginationMeta{
		CurrentPage:  pagination.CurrentPage,
		PerPage:      pagination.PerPage,
		Total:        pagination.Total,
		LastPage:     pagination.LastPage,
		FirstPageUrl: pagination.FirstPageURL,
		LastPageUrl:  pagination.LastPageURL,
	}

	if pagination.NextPageURL != nil {
		paginationMeta.NextPageUrl = *pagination.NextPageURL
	}

	if pagination.PrevPageURL != nil {
		paginationMeta.PrevPageUrl = *pagination.PrevPageURL
	}

	return helper.SuccessResponseWithPagination(c, fiber.StatusOK, "Employees retrieved successfully", pagination.Data, paginationMeta)
}

func (h *EmployeeHandler) GetEmployeeByID(c *fiber.Ctx) error {
	id := c.Params("id")

	// Get company_id from context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	employee, err := h.service.GetEmployeeByID(c.Context(), id, companyID)
	if err != nil {
		if errors.Is(err, service.ErrEmployeeNotFound) {
			return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch employee", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Employee fetched successfully", employee)
}

func (h *EmployeeHandler) UpdateEmployee(c *fiber.Ctx) error {
	id := c.Params("id")

	// Get company_id from context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	var req dto.UpdateEmployeeRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Debug: Log received phone number
	if req.PhoneNumber != nil {
		println("DEBUG: Received phoneNumber = ", *req.PhoneNumber)
	} else {
		println("DEBUG: phoneNumber is NIL")
	}

	// Validate request
	if validationErrors := customValidator.ValidateStruct(&req); len(validationErrors) > 0 {
		return helper.ValidationErrorResponse(c, validationErrors)
	}

	employee, err := h.service.UpdateEmployee(c.Context(), id, companyID, &req)
	if err != nil {
		if errors.Is(err, service.ErrEmployeeNotFound) {
			return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update employee", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Employee updated successfully", employee)
}

func (h *EmployeeHandler) DeleteEmployee(c *fiber.Ctx) error {
	id := c.Params("id")

	// Get company_id from context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	if err := h.service.DeleteEmployee(c.Context(), id, companyID); err != nil {
		if errors.Is(err, service.ErrEmployeeNotFound) {
			return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete employee", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Employee deleted successfully", nil)
}

func (h *EmployeeHandler) GetMyProfile(c *fiber.Ctx) error {
	userID := c.Locals(constants.ContextKeyUserID).(string)

	employee, err := h.service.GetMyProfile(c.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrEmployeeNotFound) {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Employee profile not found", err.Error())
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve profile", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Profile retrieved successfully", employee)
}

func (h *EmployeeHandler) UpdateMyProfile(c *fiber.Ctx) error {
	userID := c.Locals(constants.ContextKeyUserID).(string)

	var req dto.UpdateMyProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if validationErrors := customValidator.ValidateStruct(&req); len(validationErrors) > 0 {
		return helper.ValidationErrorResponse(c, validationErrors)
	}

	employee, err := h.service.UpdateMyProfile(c.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, service.ErrEmployeeNotFound) {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Employee profile not found", err.Error())
		}
		if err.Error() == "no fields to update" {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update profile", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Profile updated successfully", employee)
}
