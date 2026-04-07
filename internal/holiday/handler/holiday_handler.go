package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"hris/internal/holiday/dto"
	"hris/internal/holiday/service"
	"hris/shared/validator"
)

type HolidayHandler interface {
	Create(c *fiber.Ctx) error
	GetAllByYear(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type holidayHandler struct {
	holidayService service.HolidayService
}

func NewHolidayHandler(holidayService service.HolidayService) HolidayHandler {
	return &holidayHandler{
		holidayService: holidayService,
	}
}

func (h *holidayHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateHolidayRequest
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

	holiday, err := h.holidayService.Create(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create holiday",
			"errors":  err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Holiday created successfully",
		"data":    holiday,
	})
}

func (h *holidayHandler) GetAllByYear(c *fiber.Ctx) error {
	yearStr := c.Query("year", strconv.Itoa(time.Now().Year()))
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid year parameter",
		})
	}

	holidays, err := h.holidayService.GetAllByYear(c.Context(), year)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get holidays",
			"errors":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Holidays retrieved successfully",
		"data":    holidays,
	})
}

func (h *holidayHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Holiday ID is required",
		})
	}

	holiday, err := h.holidayService.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Holiday not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Holiday retrieved successfully",
		"data":    holiday,
	})
}

func (h *holidayHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Holiday ID is required",
		})
	}

	var req dto.UpdateHolidayRequest
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

	holiday, err := h.holidayService.Update(c.Context(), id, &req)
	if err != nil {
		if err.Error() == "holiday not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Holiday not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update holiday",
			"errors":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Holiday updated successfully",
		"data":    holiday,
	})
}

func (h *holidayHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Holiday ID is required",
		})
	}

	if err := h.holidayService.Delete(c.Context(), id); err != nil {
		if err.Error() == "holiday not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Holiday not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete holiday",
			"errors":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Holiday deleted successfully",
	})
}
