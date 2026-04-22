package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"hris/internal/dashboard/service"
	"hris/shared/constants"
	"hris/shared/helper"
)

type DashboardHandler struct {
	service service.DashboardService
}

func NewDashboardHandler(service service.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

func (h *DashboardHandler) GetSummary(c *fiber.Ctx) error {
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	summary, err := h.service.GetSummary(c.Context(), companyID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get dashboard summary", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Dashboard summary retrieved", summary)
}

func (h *DashboardHandler) GetAttendanceStats(c *fiber.Ctx) error {
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "start_date and end_date are required", nil)
	}

	stats, err := h.service.GetAttendanceStats(c.Context(), companyID, startDate, endDate)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get attendance stats", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Attendance stats retrieved", stats)
}

func (h *DashboardHandler) GetPayrollStats(c *fiber.Ctx) error {
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	monthStr := c.Query("month")
	yearStr := c.Query("year")

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid month parameter", "Month must be between 1 and 12")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2000 || year > 2100 {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid year parameter", "Year must be between 2000 and 2100")
	}

	stats, err := h.service.GetPayrollStats(c.Context(), companyID, month, year)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get payroll stats", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Payroll stats retrieved", stats)
}

func (h *DashboardHandler) GetEmployeeStats(c *fiber.Ctx) error {
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	stats, err := h.service.GetEmployeeStats(c.Context(), companyID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get employee stats", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Employee stats retrieved", stats)
}

func (h *DashboardHandler) GetRecentActivities(c *fiber.Ctx) error {
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	limitStr := c.Query("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	activities, err := h.service.GetRecentActivities(c.Context(), companyID, limit)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get recent activities", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Recent activities retrieved", activities)
}

func (h *DashboardHandler) GetSuperUserSummary(c *fiber.Ctx) error {
	summary, err := h.service.GetSuperUserSummary(c.Context())
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get superuser dashboard", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Superuser dashboard retrieved", summary)
}
