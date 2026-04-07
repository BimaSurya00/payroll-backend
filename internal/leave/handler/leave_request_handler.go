package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"hris/internal/leave/dto"
	"hris/internal/leave/service"
	"hris/shared/constants"
	"hris/shared/helper"
	customValidator "hris/shared/validator"
	"time"
)

type LeaveRequestHandler struct {
	service service.LeaveService
}

func NewLeaveRequestHandler(service service.LeaveService) *LeaveRequestHandler {
	return &LeaveRequestHandler{service: service}
}

func (h *LeaveRequestHandler) CreateLeaveRequest(c *fiber.Ctx) error {
	// Get employee ID from JWT context
	userID := c.Locals(constants.ContextKeyUserID).(string)

	var req dto.CreateLeaveRequestRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if validationErrors := customValidator.ValidateStruct(&req); len(validationErrors) > 0 {
		return helper.ValidationErrorResponse(c, validationErrors)
	}

	leaveRequest, err := h.service.CreateLeaveRequest(c.Context(), userID, &req)
	if err != nil {
		if err == service.ErrInvalidDateRange {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid date range", nil)
		}
		if err == service.ErrInsufficientBalance {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "Insufficient leave balance", nil)
		}
		if err == service.ErrOverlappingLeave {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "Overlapping leave request exists", nil)
		}
		if err == service.ErrEmployeeNotFound {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Employee record not found. Please contact HR to set up your employee profile.", nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create leave request", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusCreated, "Leave request created successfully", leaveRequest)
}

func (h *LeaveRequestHandler) GetMyLeaveRequests(c *fiber.Ctx) error {
	userID := c.Locals(constants.ContextKeyUserID).(string)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "15"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 15
	}

	requests, total, err := h.service.GetMyLeaveRequests(c.Context(), userID, page, perPage)
	if err != nil {
		if err == service.ErrEmployeeNotFound {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Employee record not found. Please contact HR to set up your employee profile.", nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get leave requests", err.Error())
	}

	return helper.SuccessResponseWithPagination(c, fiber.StatusOK, "Leave requests retrieved successfully", requests, helper.PaginationMeta{
		CurrentPage: page,
		PerPage:     perPage,
		Total:       total,
		LastPage:    int((total + int64(perPage) - 1) / int64(perPage)),
	})
}

func (h *LeaveRequestHandler) GetLeaveRequestByID(c *fiber.Ctx) error {
	id := c.Params("id")

	// Get company_id from context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	leaveRequest, err := h.service.GetLeaveRequestByID(c.Context(), id, companyID)
	if err != nil {
		if err == service.ErrLeaveRequestNotFound {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Leave request not found", nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get leave request", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Leave request retrieved successfully", leaveRequest)
}

func (h *LeaveRequestHandler) GetAllLeaveRequests(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	// Get company_id from context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	requests, total, err := h.service.GetAllLeaveRequests(c.Context(), page, perPage, companyID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get leave requests", err.Error())
	}

	return helper.SuccessResponseWithPagination(c, fiber.StatusOK, "Leave requests retrieved successfully", requests, helper.PaginationMeta{
		CurrentPage: page,
		PerPage:     perPage,
		Total:       total,
		LastPage:    int((total + int64(perPage) - 1) / int64(perPage)),
	})
}

func (h *LeaveRequestHandler) GetPendingLeaveRequests(c *fiber.Ctx) error {
	// Get company_id from context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	requests, err := h.service.GetPendingLeaveRequests(c.Context(), companyID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get pending requests", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Pending leave requests retrieved successfully", requests)
}

func (h *LeaveRequestHandler) ApproveLeaveRequest(c *fiber.Ctx) error {
	id := c.Params("id")
	approverID := c.Locals(constants.ContextKeyUserID).(string)

	var req dto.ApproveLeaveRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Get company_id from context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	leaveRequest, err := h.service.ApproveLeaveRequest(c.Context(), id, approverID, companyID, &req)
	if err != nil {
		if err == service.ErrLeaveRequestNotFound {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Leave request not found", nil)
		}
		if err == service.ErrInvalidRequestStatus {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request status", nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to approve leave request", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Leave request approved successfully", leaveRequest)
}

func (h *LeaveRequestHandler) RejectLeaveRequest(c *fiber.Ctx) error {
	id := c.Params("id")
	approverID := c.Locals(constants.ContextKeyUserID).(string)

	var req dto.RejectLeaveRequest
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

	err := h.service.RejectLeaveRequest(c.Context(), id, approverID, companyID, &req)
	if err != nil {
		if err == service.ErrLeaveRequestNotFound {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Leave request not found", nil)
		}
		if err == service.ErrInvalidRequestStatus {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request status", nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to reject leave request", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Leave request rejected successfully", nil)
}

func (h *LeaveRequestHandler) GetMyLeaveBalances(c *fiber.Ctx) error {
	userID := c.Locals(constants.ContextKeyUserID).(string)
	year, _ := strconv.Atoi(c.Query("year", strconv.Itoa(time.Now().Year())))

	balances, err := h.service.GetMyLeaveBalances(c.Context(), userID, year)
	if err != nil {
		if err == service.ErrEmployeeNotFound {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Employee record not found. Please contact HR to set up your employee profile.", nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get leave balances", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Leave balances retrieved successfully", balances)
}
