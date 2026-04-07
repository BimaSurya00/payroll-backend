package dto

type AuditLogResponse struct {
	ID           string    `json:"id"`
	UserID       string    `json:"userId"`
	UserName     string    `json:"userName"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resourceType"`
	ResourceID   string    `json:"resourceId,omitempty"`
	OldData      map[string]interface{} `json:"oldData,omitempty"`
	NewData      map[string]interface{} `json:"newData,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	IPAddress    string    `json:"ipAddress"`
	CreatedAt    string    `json:"createdAt"`
}

type AuditLogListResponse struct {
	ID           string    `json:"id"`
	UserID       string    `json:"userId"`
	UserName     string    `json:"userName"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resourceType"`
	ResourceID   string    `json:"resourceId,omitempty"`
	IPAddress    string    `json:"ipAddress"`
	CreatedAt    string    `json:"createdAt"`
}

type AuditLogPagination struct {
	Data       []AuditLogListResponse `json:"data"`
	Page       int                    `json:"page"`
	PerPage    int                    `json:"perPage"`
	Total      int64                  `json:"total"`
	TotalPages int                    `json:"totalPages"`
}
