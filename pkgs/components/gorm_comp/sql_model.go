package gorm_comp

import (
	"time"
)

type SQLModel struct {
	ID        string    `gorm:"column:id;primaryKey"`
	CreatedAt *time.Time `gorm:"column:created_at;type:timestamp without time zone;default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:timestamp without time zone;default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:timestamp without time zone"`
	CreatedBy *string    `gorm:"column:created_by"`
	UpdatedBy *string    `gorm:"column:updated_by"`
	DeletedBy *string    `gorm:"column:deleted_by"`
}
