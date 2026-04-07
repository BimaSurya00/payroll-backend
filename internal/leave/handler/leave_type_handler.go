package handler

import (
	"github.com/gofiber/fiber/v2"
	"hris/internal/leave/dto"
	"hris/internal/leave/service"
	"hris/shared/helper"
	customValidator "hris/shared/validator"
)

type LeaveTypeHandler struct {
	service service.LeaveTypeService
}

func NewLeaveTypeHandler(service service.LeaveTypeService) *LeaveTypeHandler {
	return &LeaveTypeHandler{service: service}
}

func (h *LeaveTypeHandler) GetActiveLeaveTypes(c *fiber.Ctx) error {
	types, err := h.service.GetActiveLeaveTypes(c.Context())
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get leave types", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Leave types retrieved successfully", types)
}

func (h *LeaveTypeHandler) GetAllLeaveTypes(c *fiber.Ctx) error {
	types, err := h.service.GetAllLeaveTypes(c.Context())
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get leave types", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Leave types retrieved successfully", types)
}

func (h *LeaveTypeHandler) GetLeaveTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")

	leaveType, err := h.service.GetLeaveTypeByID(c.Context(), id)
	if err != nil {
		if err == service.ErrLeaveTypeNotFound {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Leave type not found", nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get leave type", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Leave type retrieved successfully", leaveType)
}

func (h *LeaveTypeHandler) CreateLeaveType(c *fiber.Ctx) error {
	var req dto.CreateLeaveTypeRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if validationErrors := customValidator.ValidateStruct(&req); len(validationErrors) > 0 {
		return helper.ValidationErrorResponse(c, validationErrors)
	}

	leaveType, err := h.service.CreateLeaveType(c.Context(), &req)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create leave type", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusCreated, "Leave type created successfully", leaveType)
}

func (h *LeaveTypeHandler) UpdateLeaveType(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateLeaveTypeRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if validationErrors := customValidator.ValidateStruct(&req); len(validationErrors) > 0 {
		return helper.ValidationErrorResponse(c, validationErrors)
	}

	leaveType, err := h.service.UpdateLeaveType(c.Context(), id, &req)
	if err != nil {
		if err == service.ErrLeaveTypeNotFound {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Leave type not found", nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update leave type", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Leave type updated successfully", leaveType)
}

func (h *LeaveTypeHandler) DeleteLeaveType(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.service.DeleteLeaveType(c.Context(), id)
	if err != nil {
		if err == service.ErrLeaveTypeNotFound {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Leave type not found", nil)
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete leave type", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, "Leave type deleted successfully", nil)
}
