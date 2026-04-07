package entity

import "time"

type AuditLog struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"userId" db:"user_id"`
	UserName     string    `json:"userName" db:"user_name"`
	Action       string    `json:"action" db:"action"`           // CREATE, UPDATE, DELETE, APPROVE, REJECT, GENERATE
	ResourceType string    `json:"resourceType" db:"resource_type"` // payroll, leave, employee, etc.
	ResourceID   string    `json:"resourceId,omitempty" db:"resource_id"`
	OldData      *string   `json:"oldData,omitempty" db:"old_data"` // JSON string
	NewData      *string   `json:"newData,omitempty" db:"new_data"` // JSON string
	Metadata     *string   `json:"metadata,omitempty" db:"metadata"`
	IPAddress    string    `json:"ipAddress" db:"ip_address"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}
