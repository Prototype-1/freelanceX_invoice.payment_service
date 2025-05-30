package service

import (
	"context"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/model"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/repository"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/pkg"
)

type InvoiceUsecase interface {
	CreateInvoice(ctx context.Context, invoice *model.Invoice) error
	GetInvoiceByID(ctx context.Context, id string) (*model.Invoice, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	ListInvoices(ctx context.Context, filter *repository.InvoiceFilter) ([]*model.Invoice, error)
	GetInvoicePDF(ctx context.Context, invoiceID string) ([]byte, error) 
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

func (s *invoiceService) GetInvoicePDF(ctx context.Context, invoiceID string) ([]byte, error) {
	invoice, err := s.repo.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return nil, err
	}

	pdfBytes, err := s.generateInvoicePDF(invoice)
	if err != nil {
		return nil, err
	}

	return pdfBytes, nil
}

func (s *invoiceService) generateInvoicePDF(inv *model.Invoice) ([]byte, error) {
	return pkg.GenerateInvoicePDF(inv)
}