package service

import (
	"context"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/model"
)

type InvoiceUsecase interface {
	CreateInvoice(ctx context.Context, invoice *model.Invoice) error
}

