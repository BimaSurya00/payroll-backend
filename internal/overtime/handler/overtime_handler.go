package handler

import (
	"github.com/gofiber/fiber/v2"
	"hris/internal/overtime/dto"
	"hris/internal/overtime/service"
	"hris/shared/constants"
	"hris/shared/helper"
	customValidator "hris/shared/validator"
)

type OvertimeHandler struct {
	service service.OvertimeService
}

func NewOvertimeHandler(service service.OvertimeService) *OvertimeHandler {
	return &OvertimeHandler{service: service}
}

// CreateOvertimeRequest handles overtime request creation
func (h *OvertimeHandler) CreateOvertimeRequest(c *fiber.Ctx) error {
	userID := c.Locals(constants.ContextKeyUserID).(string)

	var req dto.CreateOvertimeRequestRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if validationErrors := customValidator.ValidateStruct(&req); len(validationErrors) > 0 {
		return helper.ValidationErrorResponse(c, validationErrors)
	}

	overtimeRequest, err := h.service.CreateOvertimeRequest(c.Context(), userID, &req)
	if err != nil {
		if err == service.ErrInvalidTimeRange {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid time range", nil)
		}
		if err == service.ErrDuplicateOvertimeRequest {
			return helper.ErrorResponse(c, fiber.StatusConflict, "Overtime request for this date already exists", nil)
		}
		if err == service.ErrOvertimeTooLong {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "Overtime hours exceed maximum allowed", nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create overtime request", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusCreated, "Overtime request created successfully", overtimeRequest)
}

// GetMyOvertimeRequests handles get my overtime requests
func (h *OvertimeHandler) GetMyOvertimeRequests(c *fiber.Ctx) error {
	userID := c.Locals(constants.ContextKeyUserID).(string)

	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 15)

	requests, total, err := h.service.GetMyOvertimeRequests(c.Context(), userID, page, perPage)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get overtime requests", err.Error())
	}

	// Calculate pagination
	lastPage := int(total) / perPage
	if int(total)%perPage != 0 {
		lastPage++
	}

	pagination := map[string]interface{}{
		"currentPage": page,
		"perPage":     perPage,
		"total":       total,
		"lastPage":    lastPage,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":    true,
		"statusCode": fiber.StatusOK,
		"message":    "Overtime requests retrieved successfully",
		"data":       requests,
		"pagination": pagination,
	})
}

// GetAllOvertimeRequests handles get all overtime requests (admin only)
func (h *OvertimeHandler) GetAllOvertimeRequests(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 15)
	status := c.Query("status", "")
	employeeID := c.Query("employee_id", "")

	// Get company_id from context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	requests, total, err := h.service.GetAllOvertimeRequests(c.Context(), page, perPage, status, employeeID, companyID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get overtime requests", err.Error())
	}

	// Calculate pagination
	lastPage := int(total) / perPage
	if int(total)%perPage != 0 {
		lastPage++
	}

	pagination := map[string]interface{}{
		"currentPage": page,
		"perPage":     perPage,
		"total":       total,
		"lastPage":    lastPage,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":    true,
		"statusCode": fiber.StatusOK,
		"message":    "Overtime requests retrieved successfully",
		"data":       requests,
		"pagination": pagination,
	})
}

// GetOvertimeRequestByID handles get overtime request by ID
func (h *OvertimeHandler) GetOvertimeRequestByID(c *fiber.Ctx) error {
	id := c.Params("id")

	// Get company_id from context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	request, err := h.service.GetOvertimeRequestByID(c.Context(), id, companyID)
	if err != nil {
		if err == service.ErrOvertimeRequestNotFound {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Overtime request not found", nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get overtime request", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Overtime request retrieved successfully", request)
}

// GetPendingOvertimeRequests handles get pending overtime requests (admin only)
func (h *OvertimeHandler) GetPendingOvertimeRequests(c *fiber.Ctx) error {
	// Get company_id from context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	requests, err := h.service.GetPendingOvertimeRequests(c.Context(), companyID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get pending requests", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Pending requests retrieved successfully", requests)
}

// ApproveOvertimeRequest handles approve overtime request
func (h *OvertimeHandler) ApproveOvertimeRequest(c *fiber.Ctx) error {
	id := c.Params("id")
	approverID := c.Locals(constants.ContextKeyUserID).(string)

	// Get company_id from context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	request, err := h.service.ApproveOvertimeRequest(c.Context(), id, approverID, companyID)
	if err != nil {
		if err == service.ErrOvertimeRequestNotFound {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Overtime request not found", nil)
		}
		if err == service.ErrInvalidRequestStatus {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "Cannot approve request with current status", nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to approve request", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Overtime request approved successfully", request)
}

// RejectOvertimeRequest handles reject overtime request
func (h *OvertimeHandler) RejectOvertimeRequest(c *fiber.Ctx) error {
	id := c.Params("id")
	approverID := c.Locals(constants.ContextKeyUserID).(string)

	var req dto.RejectOvertimeRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if validationErrors := customValidator.ValidateStruct(&req); len(validationErrors) > 0 {
		return helper.ValidationErrorResponse(c, validationErrors)
	}

	// Get company_id from context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	request, err := h.service.RejectOvertimeRequest(c.Context(), id, approverID, companyID, req.RejectionReason)
	if err != nil {
		if err == service.ErrOvertimeRequestNotFound {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Overtime request not found", nil)
		}
		if err == service.ErrInvalidRequestStatus {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "Cannot reject request with current status", nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to reject request", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Overtime request rejected successfully", request)
}

// ClockIn handles clock in for overtime
func (h *OvertimeHandler) ClockIn(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.ClockInRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	attendance, err := h.service.ClockIn(c.Context(), id, req.Notes)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to clock in", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Clocked in successfully", attendance)
}

// ClockOut handles clock out for overtime
func (h *OvertimeHandler) ClockOut(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.ClockOutRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	attendance, err := h.service.ClockOut(c.Context(), id, req.Notes)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to clock out", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Clocked out successfully", attendance)
}

// GetActivePolicies handles get active overtime policies
func (h *OvertimeHandler) GetActivePolicies(c *fiber.Ctx) error {
	policies, err := h.service.GetActivePolicies(c.Context())
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get policies", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Policies retrieved successfully", policies)
}

// CalculateOvertimePay handles calculate overtime pay
func (h *OvertimeHandler) CalculateOvertimePay(c *fiber.Ctx) error {
	employeeID := c.Params("employeeId")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "start_date and end_date are required", nil)
	}

	calculation, err := h.service.CalculateOvertimePay(c.Context(), employeeID, startDate, endDate)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to calculate overtime pay", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Overtime pay calculated successfully", calculation)
}
