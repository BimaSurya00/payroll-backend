package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"hris/internal/audit/dto"
	auditentity "hris/internal/audit/entity"
	auditrepo "hris/internal/audit/repository"
)

type AuditService interface {
	Log(ctx context.Context, entry AuditEntry) error
	GetByResource(ctx context.Context, resourceType, resourceID string, page, perPage int) (*dto.AuditLogPagination, error)
	GetAll(ctx context.Context, filter auditrepo.AuditFilter, page, perPage int, path string) (*dto.AuditLogPagination, error)
}

type AuditEntry struct {
	UserID       string
	UserName     string
	Action       string
	ResourceType string
	ResourceID   string
	OldData      interface{}
	NewData      interface{}
	Metadata     map[string]interface{}
	IPAddress    string
}

type auditService struct {
	auditRepo auditrepo.AuditRepository
}

func NewAuditService(auditRepo auditrepo.AuditRepository) AuditService {
	return &auditService{
		auditRepo: auditRepo,
	}
}

func (s *auditService) Log(ctx context.Context, entry AuditEntry) error {
	now := time.Now()

	// Marshal old and new data to JSON strings for JSONB columns
	var oldDataJSON, newDataJSON, metadataJSON *string
	if entry.OldData != nil {
		bytes, err := json.Marshal(entry.OldData)
		if err == nil {
			str := string(bytes)
			oldDataJSON = &str
		}
	}
	if entry.NewData != nil {
		bytes, err := json.Marshal(entry.NewData)
		if err == nil {
			str := string(bytes)
			newDataJSON = &str
		}
	}
	if entry.Metadata != nil {
		bytes, err := json.Marshal(entry.Metadata)
		if err == nil {
			str := string(bytes)
			metadataJSON = &str
		}
	}

	log := &auditentity.AuditLog{
		ID:           uuid.New().String(),
		UserID:       entry.UserID,
		UserName:     entry.UserName,
		Action:       entry.Action,
		ResourceType: entry.ResourceType,
		ResourceID:   entry.ResourceID,
		OldData:      oldDataJSON,
		NewData:      newDataJSON,
		Metadata:     metadataJSON,
		IPAddress:    entry.IPAddress,
		CreatedAt:    now,
	}

	return s.auditRepo.Create(ctx, log)
}

func (s *auditService) GetByResource(ctx context.Context, resourceType, resourceID string, page, perPage int) (*dto.AuditLogPagination, error) {
	offset := (page - 1) * perPage

	logs, err := s.auditRepo.FindByResource(ctx, resourceType, resourceID, perPage, offset)
	if err != nil {
		return nil, err
	}

	// Get total count for this resource
	resType := resourceType
	total, err := s.auditRepo.Count(ctx, auditrepo.AuditFilter{
		ResourceType: &resType,
	})
	if err != nil {
		return nil, err
	}

	data := make([]dto.AuditLogListResponse, len(logs))
	for i, log := range logs {
		data[i] = dto.AuditLogListResponse{
			ID:           log.ID,
			UserID:       log.UserID,
			UserName:     log.UserName,
			Action:       log.Action,
			ResourceType: log.ResourceType,
			ResourceID:   log.ResourceID,
			IPAddress:    log.IPAddress,
			CreatedAt:    log.CreatedAt.Format(time.RFC3339),
		}
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &dto.AuditLogPagination{
		Data:       data,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *auditService) GetAll(ctx context.Context, filter auditrepo.AuditFilter, page, perPage int, path string) (*dto.AuditLogPagination, error) {
	offset := (page - 1) * perPage

	logs, err := s.auditRepo.FindAll(ctx, filter, perPage, offset)
	if err != nil {
		return nil, err
	}

	total, err := s.auditRepo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	data := make([]dto.AuditLogListResponse, len(logs))
	for i, log := range logs {
		data[i] = dto.AuditLogListResponse{
			ID:           log.ID,
			UserID:       log.UserID,
			UserName:     log.UserName,
			Action:       log.Action,
			ResourceType: log.ResourceType,
			ResourceID:   log.ResourceID,
			IPAddress:    log.IPAddress,
			CreatedAt:    log.CreatedAt.Format(time.RFC3339),
		}
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &dto.AuditLogPagination{
		Data:       data,
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}
