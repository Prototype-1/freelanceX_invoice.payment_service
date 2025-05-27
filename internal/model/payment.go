package model

import (
	"time"
	"github.com/google/uuid"
)

type Payment struct {
	ID             uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	InvoiceID      uuid.UUID  `gorm:"type:uuid;not null"`
	MilestoneID    *uuid.UUID `gorm:"type:uuid"` 
	PayerID        uuid.UUID  `gorm:"type:uuid;not null"` 
	ReceiverID     uuid.UUID  `gorm:"type:uuid;not null"`
	AmountPaid     float64    `gorm:"type:numeric;not null"`
	PlatformFee    float64    `gorm:"type:numeric;not null"`
	AmountCredited float64    `gorm:"type:numeric;not null"` 
	Status         string     `gorm:"type:varchar(50);default:'completed'"`
	CreatedAt      time.Time  `gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime"`
}
