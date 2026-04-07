package handler

import (
	"regexp"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"hris/internal/schedule/dto"
	"hris/internal/schedule/service"
	"hris/shared/constants"
	"hris/shared/helper"
	customValidator "hris/shared/validator"
)

type ScheduleHandler struct {
	service service.ScheduleService
}

func NewScheduleHandler(service service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{
		service: service,
	}
}

func (h *ScheduleHandler) CreateSchedule(c *fiber.Ctx) error {
	var req dto.CreateScheduleRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Manual validation for time format
	if !isValidTimeFormat(req.TimeIn) {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid timeIn format", "Use HH:MM format (e.g., 09:00, 17:30)")
	}
	if !isValidTimeFormat(req.TimeOut) {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid timeOut format", "Use HH:MM format (e.g., 09:00, 17:30)")
	}

	if validationErrors := customValidator.ValidateStruct(&req); len(validationErrors) > 0 {
		return helper.ValidationErrorResponse(c, validationErrors)
	}

	// Get company_id from context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	schedule, err := h.service.CreateSchedule(c.Context(), &req, companyID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create schedule", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusCreated, "Schedule created successfully", schedule)
}

func isValidTimeFormat(time string) bool {
	matched, _ := regexp.MatchString(`^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$`, time)
	return matched
}

func (h *ScheduleHandler) GetAllSchedules(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", strconv.Itoa(constants.DefaultPage)))
	perPage, _ := strconv.Atoi(c.Query("per_page", strconv.Itoa(constants.DefaultPerPage)))

	if page < 1 {
		page = constants.DefaultPage
	}
	if perPage < 1 || perPage > constants.MaxPerPage {
		perPage = constants.DefaultPerPage
	}

	path := c.Protocol() + "://" + c.Hostname() + c.Path()

	// Get company_id from context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	pagination, err := h.service.GetAllSchedules(c.Context(), page, perPage, path, companyID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch schedules", err.Error())
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

	return helper.SuccessResponseWithPagination(c, fiber.StatusOK, "Schedules retrieved successfully", pagination.Data, paginationMeta)
}

func (h *ScheduleHandler) GetScheduleByID(c *fiber.Ctx) error {
	id := c.Params("id")

	// Get company_id from context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	schedule, err := h.service.GetScheduleByID(c.Context(), id, companyID)
	if err != nil {
		if err == service.ErrScheduleNotFound {
			return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch schedule", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Schedule fetched successfully", schedule)
}

func (h *ScheduleHandler) UpdateSchedule(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateScheduleRequest
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

	schedule, err := h.service.UpdateSchedule(c.Context(), id, companyID, &req)
	if err != nil {
		if err == service.ErrScheduleNotFound {
			return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update schedule", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Schedule updated successfully", schedule)
}

func (h *ScheduleHandler) DeleteSchedule(c *fiber.Ctx) error {
	id := c.Params("id")

	// Get company_id from context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	if err := h.service.DeleteSchedule(c.Context(), id, companyID); err != nil {
		if err == service.ErrScheduleNotFound {
			return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		if err == service.ErrScheduleInUse {
			return helper.ErrorResponse(c, fiber.StatusConflict, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete schedule", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Schedule deleted successfully", nil)
}
