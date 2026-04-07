package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	auditrepo "hris/internal/audit/repository"
	"hris/internal/audit/service"
	"hris/shared/helper"
)

type AuditHandler interface {
	GetAll(c *fiber.Ctx) error
	GetByResource(c *fiber.Ctx) error
}

type auditHandler struct {
	auditService service.AuditService
}

func NewAuditHandler(auditService service.AuditService) AuditHandler {
	return &auditHandler{
		auditService: auditService,
	}
}

func (h *auditHandler) GetAll(c *fiber.Ctx) error {
	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "50"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 50
	}

	// Build filter
	filter := auditrepo.AuditFilter{}

	if action := c.Query("action"); action != "" {
		filter.Action = &action
	}

	if resourceType := c.Query("resource_type"); resourceType != "" {
		filter.ResourceType = &resourceType
	}

	if userID := c.Query("user_id"); userID != "" {
		filter.UserID = &userID
	}

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		t, err := time.Parse("2006-01-02", dateFrom)
		if err == nil {
			filter.DateFrom = &t
		}
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		t, err := time.Parse("2006-01-02", dateTo)
		if err == nil {
			filter.DateTo = &t
		}
	}

	pagination, err := h.auditService.GetAll(c.Context(), filter, page, perPage, c.Path())
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get audit logs", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Audit logs retrieved successfully",
		"data":    pagination,
	})
}

func (h *auditHandler) GetByResource(c *fiber.Ctx) error {
	resourceType := c.Params("resourceType")
	resourceID := c.Params("resourceId")

	if resourceType == "" || resourceID == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Resource type and ID are required", "")
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "20"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	pagination, err := h.auditService.GetByResource(c.Context(), resourceType, resourceID, page, perPage)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get audit logs", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Audit logs retrieved successfully",
		"data":    pagination,
	})
}
