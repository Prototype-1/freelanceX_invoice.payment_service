package model

import (
	"time"
	"github.com/google/uuid"
)

type MilestoneRule struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProjectID      uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"` // one rule per project
	AdvancePercent float64   `gorm:"type:numeric;check:advance_percent >= 0 AND advance_percent <= 100"`
	MidPercent     float64   `gorm:"type:numeric;check:mid_percent >= 0 AND mid_percent <= 100"`
	FinalPercent   float64   `gorm:"type:numeric;check:final_percent >= 0 AND final_percent <= 100"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}
