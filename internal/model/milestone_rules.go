package model

import (
	"time"
	"github.com/google/uuid"
)

type MilestoneRule struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProjectID uuid.UUID `gorm:"type:uuid;not null;index"` 
	Phase     string    `gorm:"type:varchar(50);not null"` 
	Amount    float64   `gorm:"type:numeric;not null"`    
	DueDate   *time.Time
	CreatedAt time.Time
}

