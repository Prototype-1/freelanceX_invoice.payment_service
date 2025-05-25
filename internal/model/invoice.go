package model

import (
	"time"
	"github.com/google/uuid"
)

type Invoice struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProjectID     uuid.UUID `gorm:"type:uuid;not null"`
	ClientID      uuid.UUID `gorm:"type:uuid;not null"`
	FreelancerID  uuid.UUID `gorm:"type:uuid;not null"`
	Type          string    `gorm:"type:varchar(50);not null"` 
	Status        string    `gorm:"type:varchar(50);default:'draft'"`
	Amount        float64   `gorm:"type:numeric;not null"`
	Description   string    `gorm:"type:text"`
	DueDate       *time.Time
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}
