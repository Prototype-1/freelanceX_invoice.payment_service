package model

import (
	"time"
	"github.com/google/uuid"
)

type TimeEntry struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProjectID    uuid.UUID `gorm:"type:uuid;not null"`
	FreelancerID uuid.UUID `gorm:"type:uuid;not null"`
	HoursWorked  float64   `gorm:"type:numeric;not null"`
	WorkDate     time.Time `gorm:"type:date;not null"`
}
