package repository

import (
	"context"
	"gorm.io/gorm"
	"github.com/google/uuid"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/model"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *model.Payment) error
	MarkPaid(ctx context.Context, paymentID uuid.UUID) error
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

func (r *paymentRepository) MarkPaid(ctx context.Context, paymentID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.Payment{}).
		Where("id = ?", paymentID).
		Update("status", "paid").Error
}