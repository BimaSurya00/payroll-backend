package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"hris/internal/payroll/dto"
	"hris/internal/payroll/service"
	"hris/shared/constants"
	"hris/shared/helper"
	customValidator "hris/shared/validator"
)

type PayrollHandler struct {
	service service.PayrollService
}

func NewPayrollHandler(service service.PayrollService) *PayrollHandler {
	return &PayrollHandler{
		service: service,
	}
}

func (h *PayrollHandler) GenerateBulk(c *fiber.Ctx) error {
	var req dto.GeneratePayrollRequest

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

	result, err := h.service.GenerateBulk(c.Context(), &req, companyID)
	if err != nil {
		if err == service.ErrPayrollAlreadyExists {
			return helper.ErrorResponse(c, fiber.StatusConflict, err.Error(), nil)
		}
		if err == service.ErrNoEmployeesFound {
			return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate payroll", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusCreated, "Payroll generated successfully", result)
}

func (h *PayrollHandler) GetAllPayrolls(c *fiber.Ctx) error {
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

	pagination, err := h.service.GetAll(c.Context(), page, perPage, path, companyID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch payrolls", err.Error())
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

	return helper.SuccessResponseWithPagination(c, fiber.StatusOK, "Payrolls retrieved successfully", pagination.Data, paginationMeta)
}

func (h *PayrollHandler) GetPayrollByID(c *fiber.Ctx) error {
	id := c.Params("id")

	// Get company_id from context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	payroll, err := h.service.GetByID(c.Context(), id, companyID)
	if err != nil {
		if err == service.ErrPayrollNotFound {
			return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch payroll", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Payroll fetched successfully", payroll)
}

func (h *PayrollHandler) UpdateStatus(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdatePayrollStatusRequest
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

	if err := h.service.UpdateStatus(c.Context(), id, companyID, &req); err != nil {
		if err == service.ErrPayrollNotFound {
			return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		if err == service.ErrInvalidStatus {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update status", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Status updated successfully", nil)
}

func (h *PayrollHandler) ExportCSV(c *fiber.Ctx) error {
	monthStr := c.Query("month")
	yearStr := c.Query("year")

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

	csvData, filename, err := h.service.ExportCSV(c.Context(), month, year, companyID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to export CSV", err.Error())
	}

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename="+filename)
	return c.Send(csvData)
}

func (h *PayrollHandler) GetMyPayrolls(c *fiber.Ctx) error {
	userID := c.Locals(constants.ContextKeyUserID).(string)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "15"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 15
	}

	path := c.Protocol() + "://" + c.Hostname() + c.Path()

	result, err := h.service.GetMyPayrolls(c.Context(), userID, page, perPage, path)
	if err != nil {
		if err == service.ErrEmployeeNotFound {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Employee record not found. Please contact HR to set up your employee profile.", nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch payrolls", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "My payrolls retrieved", result)
}

func (h *PayrollHandler) GetMyPayrollByID(c *fiber.Ctx) error {
	userID := c.Locals(constants.ContextKeyUserID).(string)
	payrollID := c.Params("id")

	result, err := h.service.GetMyPayrollByID(c.Context(), userID, payrollID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Payroll not found", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Payroll detail retrieved", result)
}
