package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"hris/internal/attendance/dto"
	"hris/internal/attendance/service"
	"hris/shared/constants"
	"hris/shared/helper"
	customValidator "hris/shared/validator"
)

type AttendanceHandler struct {
	service service.AttendanceService
}

func NewAttendanceHandler(service service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{
		service: service,
	}
}

func (h *AttendanceHandler) ClockIn(c *fiber.Ctx) error {
	var req dto.ClockInRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if validationErrors := customValidator.ValidateStruct(&req); len(validationErrors) > 0 {
		return helper.ValidationErrorResponse(c, validationErrors)
	}

	// Get current user role from JWT context
	userRole := c.Locals(constants.ContextKeyUserRole).(string)

	// Admin and Super User do not need to clock in
	if userRole == constants.RoleAdmin || userRole == constants.RoleSuperUser {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Admin and Super User do not need to clock in", nil)
	}

	// Get current user ID from JWT context
	userID := c.Locals(constants.ContextKeyUserID).(string)

	resp, err := h.service.ClockIn(c.Context(), userID, &req)
	if err != nil {
		if err == service.ErrAlreadyClockedIn {
			return helper.ErrorResponse(c, fiber.StatusConflict, err.Error(), nil)
		}
		if err == service.ErrOutOfOfficeRange {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}
		if err == service.ErrEmployeeNotFound || err == service.ErrScheduleNotFound {
			return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to clock in", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusCreated, "Clock in successful", resp)
}

func (h *AttendanceHandler) ClockOut(c *fiber.Ctx) error {
	var req dto.ClockOutRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if validationErrors := customValidator.ValidateStruct(&req); len(validationErrors) > 0 {
		return helper.ValidationErrorResponse(c, validationErrors)
	}

	// Get current user role from JWT context
	userRole := c.Locals(constants.ContextKeyUserRole).(string)

	// Admin and Super User do not need to clock out
	if userRole == constants.RoleAdmin || userRole == constants.RoleSuperUser {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Admin and Super User do not need to clock out", nil)
	}

	// Get current user ID from JWT context
	userID := c.Locals(constants.ContextKeyUserID).(string)

	resp, err := h.service.ClockOut(c.Context(), userID, &req)
	if err != nil {
		if err == service.ErrNotClockedIn {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}
		if err == service.ErrEmployeeNotFound {
			return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to clock out", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Clock out successful", resp)
}

func (h *AttendanceHandler) GetHistory(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", strconv.Itoa(constants.DefaultPage)))
	perPage, _ := strconv.Atoi(c.Query("per_page", strconv.Itoa(constants.DefaultPerPage)))

	if page < 1 {
		page = constants.DefaultPage
	}
	if perPage < 1 || perPage > constants.MaxPerPage {
		perPage = constants.DefaultPerPage
	}

	path := c.Protocol() + "://" + c.Hostname() + c.Path()

	// Get current user ID from JWT context
	userID := c.Locals(constants.ContextKeyUserID).(string)

	pagination, err := h.service.GetHistory(c.Context(), userID, page, perPage, path)
	if err != nil {
		if err == service.ErrEmployeeNotFound {
			return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch attendance history", err.Error())
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

	return helper.SuccessResponseWithPagination(c, fiber.StatusOK, "Attendance history retrieved successfully", pagination.Data, paginationMeta)
}

func (h *AttendanceHandler) GetAllAttendances(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", strconv.Itoa(constants.DefaultPage)))
	perPage, _ := strconv.Atoi(c.Query("per_page", strconv.Itoa(constants.DefaultPerPage)))

	if page < 1 {
		page = constants.DefaultPage
	}
	if perPage < 1 || perPage > constants.MaxPerPage {
		perPage = constants.DefaultPerPage
	}

	path := c.Protocol() + "://" + c.Hostname() + c.Path()

	// Parse filters
	var filter service.GetAllAttendanceFilter

	// Filter by employee_id (optional)
	if employeeIDStr := c.Query("employee_id"); employeeIDStr != "" {
		employeeUUID, err := uuid.Parse(employeeIDStr)
		if err != nil {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid employee_id format", err.Error())
		}
		filter.EmployeeID = &employeeUUID
	}

	// Filter by schedule_id (optional)
	if scheduleIDStr := c.Query("schedule_id"); scheduleIDStr != "" {
		scheduleUUID, err := uuid.Parse(scheduleIDStr)
		if err != nil {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid schedule_id format", err.Error())
		}
		filter.ScheduleID = &scheduleUUID
	}

	// Filter by status (optional)
	if statusStr := c.Query("status"); statusStr != "" {
		filter.Status = &statusStr
	}

	// Filter by date_from (optional)
	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		dateFrom, err := time.Parse("2006-01-02", dateFromStr)
		if err != nil {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid date_from format, use YYYY-MM-DD", err.Error())
		}
		filter.DateFrom = &dateFrom
	}

	// Filter by date_to (optional)
	if dateToStr := c.Query("date_to"); dateToStr != "" {
		dateTo, err := time.Parse("2006-01-02", dateToStr)
		if err != nil {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid date_to format, use YYYY-MM-DD", err.Error())
		}
		filter.DateTo = &dateTo
	}

	// Get company_id from context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	pagination, err := h.service.GetAllAttendances(c.Context(), filter, page, perPage, path, companyID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch attendances", err.Error())
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

	return helper.SuccessResponseWithPagination(c, fiber.StatusOK, "Attendances retrieved successfully", pagination.Data, paginationMeta)
}

func (h *AttendanceHandler) GetMonthlyReport(c *fiber.Ctx) error {
	// Parse query parameters
	monthStr := c.Query("month", strconv.Itoa(int(time.Now().Month())))
	yearStr := c.Query("year", strconv.Itoa(time.Now().Year()))

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid month parameter", "Month must be between 1 and 12")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2020 || year > 2100 {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid year parameter", "Year must be between 2020 and 2100")
	}

	// Get company_id from context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	// Get report
	report, err := h.service.GetMonthlyReport(c.Context(), month, year, companyID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get monthly report", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Monthly attendance report retrieved successfully",
		"data":    report,
	})
}

func (h *AttendanceHandler) GetMyMonthlySummary(c *fiber.Ctx) error {
	// Get user ID from context
	userID := c.Locals(constants.ContextKeyUserID)
	if userID == nil {
		return helper.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", "User ID not found in context")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Invalid user ID", "User ID is not a string")
	}

	// Parse query parameters
	monthStr := c.Query("month", strconv.Itoa(int(time.Now().Month())))
	yearStr := c.Query("year", strconv.Itoa(time.Now().Year()))

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid month parameter", "Month must be between 1 and 12")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2020 || year > 2100 {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid year parameter", "Year must be between 2020 and 2100")
	}

	// Get summary
	summary, err := h.service.GetMyMonthlySummary(c.Context(), userIDStr, month, year)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get monthly summary", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Monthly attendance summary retrieved successfully",
		"data":    summary,
	})
}

func (h *AttendanceHandler) CreateCorrection(c *fiber.Ctx) error {
	var req dto.CreateCorrectionRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := customValidator.ValidateStruct(req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	// Get admin ID from context
	adminID := c.Locals(constants.ContextKeyUserID)
	if adminID == nil {
		return helper.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", "User ID not found in context")
	}

	adminIDStr, ok := adminID.(string)
	if !ok {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Invalid user ID", "User ID is not a string")
	}

	// Get company_id from context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	attendance, err := h.service.CreateCorrection(c.Context(), adminIDStr, &req, companyID)
	if err != nil {
		if err.Error() == "attendance already exists for this date, use update correction instead" {
			return helper.ErrorResponse(c, fiber.StatusConflict, "Attendance already exists", "Use update correction endpoint instead")
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create correction", err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Attendance correction created successfully",
		"data":    attendance,
	})
}

func (h *AttendanceHandler) UpdateCorrection(c *fiber.Ctx) error {
	attendanceID := c.Params("id")
	if attendanceID == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Attendance ID is required", "")
	}

	var req dto.UpdateCorrectionRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := customValidator.ValidateStruct(req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	// Get admin ID from context
	adminID := c.Locals(constants.ContextKeyUserID)
	if adminID == nil {
		return helper.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", "User ID not found in context")
	}

	adminIDStr, ok := adminID.(string)
	if !ok {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Invalid user ID", "User ID is not a string")
	}

	// Get company_id from context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	attendance, err := h.service.UpdateCorrection(c.Context(), adminIDStr, attendanceID, &req, companyID)
	if err != nil {
		if err.Error() == "attendance not found" {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Attendance not found", "")
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update correction", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Attendance correction updated successfully",
		"data":    attendance,
	})
}
