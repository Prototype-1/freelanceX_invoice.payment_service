package repository

import (
	"context"
	"gorm.io/gorm"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/model"
)

type InvoiceRepository interface {
	CreateInvoice(ctx context.Context, invoice *model.Invoice) error
	// future: GetInvoiceByID, UpdateStatus, ListInvoices etc.
}

type invoiceRepo struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return &invoiceRepo{db: db}
}

func (r *invoiceRepo) CreateInvoice(ctx context.Context, invoice *model.Invoice) error {
	return r.db.WithContext(ctx).Create(invoice).Error
}