package repository

import (
	"context"
	"gorm.io/gorm"
	"github.com/google/uuid"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/model"
)

type InvoiceRepository interface {
	CreateInvoice(ctx context.Context, invoice *model.Invoice) error
	GetInvoiceByID(ctx context.Context, id string) (*model.Invoice, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	ListInvoices(ctx context.Context, filter *InvoiceFilter) ([]*model.Invoice, error)
	FindInvoiceByProjectAndPhase(ctx context.Context, projectID uuid.UUID, phase string) (*model.Invoice, error)
}

type InvoiceFilter struct {
	ClientID     *string
	FreelancerID *string
	ProjectID    *string
	Status       *string
	Type         *string
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

func (r *invoiceRepo) GetInvoiceByID(ctx context.Context, id string) (*model.Invoice, error) {
	var invoice model.Invoice
	if err := r.db.WithContext(ctx).First(&invoice, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *invoiceRepo) UpdateStatus(ctx context.Context, id string, status string) error {
	return r.db.WithContext(ctx).Model(&model.Invoice{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *invoiceRepo) ListInvoices(ctx context.Context, filter *InvoiceFilter) ([]*model.Invoice, error) {
	var invoices []*model.Invoice
	query := r.db.WithContext(ctx).Model(&model.Invoice{})

	if filter.ClientID != nil {
		query = query.Where("client_id = ?", *filter.ClientID)
	}
	if filter.FreelancerID != nil {
		query = query.Where("freelancer_id = ?", *filter.FreelancerID)
	}
	if filter.ProjectID != nil {
		query = query.Where("project_id = ?", *filter.ProjectID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Type != nil {
		query = query.Where("type = ?", *filter.Type)
	}

	if err := query.Find(&invoices).Error; err != nil {
		return nil, err
	}

	return invoices, nil
}

func (r *invoiceRepo) FindInvoiceByProjectAndPhase(ctx context.Context, projectID uuid.UUID, phase string) (*model.Invoice, error) {
	var invoice model.Invoice
	err := r.db.WithContext(ctx).Where("project_id = ? AND milestone_phase = ?", projectID, phase).First(&invoice).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}
