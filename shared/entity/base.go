package entity

import (
	"time"

	"github.com/google/uuid"
)

type BaseEntity struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func NewBaseEntity() BaseEntity {
	now := time.Now()
	return BaseEntity{
		ID:        uuid.New().String(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (b *BaseEntity) UpdateTimestamp() {
	b.UpdatedAt = time.Now()
}