package base

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
	CreatedBy *string    `json:"createdBy,omitempty"`
	UpdatedBy *string    `json:"updatedBy,omitempty"`
	DeletedBy *string    `json:"deletedBy,omitempty"`
}

func GenerateBaseModel() *BaseModel {
	now := time.Now()
	return &BaseModel{
		ID:        uuid.New(),
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}
