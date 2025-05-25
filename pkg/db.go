package db

import (
	"fmt"
	"log"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dsn string) (*gorm.DB, error) {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = AutoMigrate()
	if err != nil {
		return nil, err
	}

	return DB, nil
}

func AutoMigrate() error {
	err := DB.AutoMigrate(
		&model.Invoice{},
		&model.MilestoneRule{},	
	)
	if err != nil {
		log.Printf("AutoMigrate error: %v", err)
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	log.Println("Database migration successful")
	return nil
}
