package repository

import (
	"context"
	"gorm.io/gorm"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/model"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *model.Payment) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}
