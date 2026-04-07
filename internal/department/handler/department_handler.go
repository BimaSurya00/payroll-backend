package handler

import (
	"github.com/gofiber/fiber/v2"
	"hris/internal/department/dto"
	"hris/internal/department/service"
	"hris/shared/constants"
	"hris/shared/validator"
)

type DepartmentHandler interface {
	Create(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type departmentHandler struct {
	departmentService service.DepartmentService
}

func NewDepartmentHandler(departmentService service.DepartmentService) DepartmentHandler {
	return &departmentHandler{
		departmentService: departmentService,
	}
}

func (h *departmentHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateDepartmentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
	}

	if err := validator.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  err,
		})
	}

	// Get company_id from context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Company context not found",
		})
	}

	department, err := h.departmentService.Create(c.Context(), &req, companyID)
	if err != nil {
		errMsg := err.Error()
		statusCode := fiber.StatusInternalServerError

		// Handle specific errors with appropriate status codes
		switch errMsg {
		case "department code already exists":
			statusCode = fiber.StatusConflict // 409
		case "head employee not found":
			statusCode = fiber.StatusBadRequest // 400
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create department",
			"errors":  errMsg,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Department created successfully",
		"data":    department,
	})
}

func (h *departmentHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Department ID is required",
		})
	}

	// Get company_id from context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Company context not found",
		})
	}

	department, err := h.departmentService.GetByID(c.Context(), id, companyID)
	if err != nil {
		if err.Error() == "department not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Department not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get department",
			"errors":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Department retrieved successfully",
		"data":    department,
	})
}

func (h *departmentHandler) GetAll(c *fiber.Ctx) error {
	// Get company_id from context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Company context not found",
		})
	}

	departments, err := h.departmentService.GetAll(c.Context(), companyID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get departments",
			"errors":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Departments retrieved successfully",
		"data":    departments,
	})
}

func (h *departmentHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Department ID is required",
		})
	}

	var req dto.UpdateDepartmentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
	}

	if err := validator.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  err,
		})
	}

	// Get company_id from context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Company context not found",
		})
	}

	department, err := h.departmentService.Update(c.Context(), id, companyID, &req)
	if err != nil {
		if err.Error() == "department not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Department not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update department",
			"errors":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Department updated successfully",
		"data":    department,
	})
}

func (h *departmentHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Department ID is required",
		})
	}

	// Get company_id from context
	companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
	if !ok || companyID == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Company context not found",
		})
	}

	if err := h.departmentService.Delete(c.Context(), id, companyID); err != nil {
		if err.Error() == "department not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Department not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete department",
			"errors":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Department deleted successfully",
	})
}
