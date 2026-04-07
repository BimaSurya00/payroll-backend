package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"hris/internal/company/dto"
	"hris/internal/company/service"
	"hris/shared/constants"
	"hris/shared/helper"
	customValidator "hris/shared/validator"
)

type CompanyHandler struct {
	service service.CompanyService
}

func NewCompanyHandler(svc service.CompanyService) *CompanyHandler {
	return &CompanyHandler{service: svc}
}

func (h *CompanyHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateCompanyRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if validationErrors := customValidator.ValidateStruct(&req); len(validationErrors) > 0 {
		return helper.ValidationErrorResponse(c, validationErrors)
	}

	company, err := h.service.Create(c.Context(), &req)
	if err != nil {
		if err == service.ErrSlugAlreadyExists {
			return helper.ErrorResponse(c, fiber.StatusConflict, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create company", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusCreated, "Company created successfully", company)
}

func (h *CompanyHandler) GetCurrent(c *fiber.Ctx) error {
	// Get company ID from JWT context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", nil)
	}

	company, err := h.service.GetByID(c.Context(), companyID)
	if err != nil {
		if err == service.ErrCompanyNotFound {
			return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get company", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Company retrieved", company)
}

func (h *CompanyHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Company ID is required", nil)
	}

	company, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		if err == service.ErrCompanyNotFound {
			return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get company", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Company retrieved", company)
}

func (h *CompanyHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", strconv.Itoa(constants.DefaultPage)))
	perPage, _ := strconv.Atoi(c.Query("per_page", strconv.Itoa(constants.DefaultPerPage)))

	if page < 1 {
		page = constants.DefaultPage
	}
	if perPage < 1 || perPage > constants.MaxPerPage {
		perPage = constants.DefaultPerPage
	}

	companies, total, err := h.service.GetAll(c.Context(), page, perPage)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get companies", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Companies retrieved",
		"data":    companies,
		"meta": fiber.Map{
			"page":    page,
			"perPage": perPage,
			"total":   total,
		},
	})
}

func (h *CompanyHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Company ID is required", nil)
	}

	var req dto.UpdateCompanyRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	company, err := h.service.Update(c.Context(), id, &req)
	if err != nil {
		if err == service.ErrCompanyNotFound {
			return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		if err == service.ErrSlugAlreadyExists {
			return helper.ErrorResponse(c, fiber.StatusConflict, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update company", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Company updated successfully", company)
}

func (h *CompanyHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Company ID is required", nil)
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		if err == service.ErrCompanyNotFound {
			return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete company", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Company deleted successfully", nil)
}
