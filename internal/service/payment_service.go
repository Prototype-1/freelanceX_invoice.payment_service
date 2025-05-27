package service

import (
	"context"
	"time"
	"github.com/google/uuid"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/model"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/pkg"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/repository"
)

type PaymentService interface {
	CreatePaymentOrder(
	ctx context.Context,
	invoiceID, milestoneID, payerID, receiverID string,
	amount float64,
) (*model.Payment, *pkg.CreateOrderResponse, error)
}

type paymentService struct {
	paymentRepo repository.PaymentRepository
	invoiceRepo repository.InvoiceRepository
	milestoneRepo repository.MilestoneRuleRepository
}

func (u *paymentService) CreatePaymentOrder(
	ctx context.Context,
	invoiceID, milestoneID, payerID, receiverID string,
	amount float64,
) (*model.Payment, *pkg.CreateOrderResponse, error) {

	invoiceUUID := uuid.MustParse(invoiceID)
	payerUUID := uuid.MustParse(payerID)
	receiverUUID := uuid.MustParse(receiverID)

	rzClient := pkg.NewRazorpayClient()
	order, err := rzClient.CreateOrder(amount, "receipt_"+invoiceID)
	if err != nil {
		return nil, nil, err
	}

	payment := &model.Payment{
		ID:             uuid.New(),
		InvoiceID:      invoiceUUID,
		PayerID:        payerUUID,
		ReceiverID:     receiverUUID,
		AmountPaid:     amount,
		Status:         "pending",
		OrderID:        order.ID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if milestoneID != "" {
		mid := uuid.MustParse(milestoneID)
		payment.MilestoneID = &mid
	}

	err = u.paymentRepo.Create(ctx, payment)
	if err != nil {
		return nil, nil, err
	}

	return payment, order, nil
}

func NewPaymentService(
	paymentRepo repository.PaymentRepository,
	invoiceRepo repository.InvoiceRepository,
	milestoneRepo repository.MilestoneRuleRepository,
) PaymentService {
	return &paymentService{
		paymentRepo:   paymentRepo,
		invoiceRepo:   invoiceRepo,
		milestoneRepo: milestoneRepo,
	}
}
