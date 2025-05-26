package service

import (
	"context"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/model"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/repository"
)

type InvoiceUsecase interface {
	CreateInvoice(ctx context.Context, invoice *model.Invoice) error
	GetInvoiceByID(ctx context.Context, id string) (*model.Invoice, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	ListInvoices(ctx context.Context, filter *repository.InvoiceFilter) ([]*model.Invoice, error)
}

type invoiceService struct {
	repo repository.InvoiceRepository
}

func NewInvoiceService(repo repository.InvoiceRepository) InvoiceUsecase {
	return &invoiceService{repo: repo}
}

func (s *invoiceService) CreateInvoice(ctx context.Context, invoice *model.Invoice) error {
	return s.repo.CreateInvoice(ctx, invoice)
}

func (s *invoiceService) GetInvoiceByID(ctx context.Context, id string) (*model.Invoice, error) {
	return s.repo.GetInvoiceByID(ctx, id)
}

func (s *invoiceService) UpdateStatus(ctx context.Context, id string, status string) error {
	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *invoiceService) ListInvoices(ctx context.Context, filter *repository.InvoiceFilter) ([]*model.Invoice, error) {
	return s.repo.ListInvoices(ctx, filter)
}

