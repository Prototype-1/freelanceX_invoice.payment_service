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
	ProcessSimulatedPayment(ctx context.Context, invoiceID, milestoneID, payerID, receiverID string, amount float64) (*model.Payment, error)
}

type paymentService struct {
	paymentRepo repository.PaymentRepository
	invoiceRepo repository.InvoiceRepository
	milestoneRepo repository.MilestoneRuleRepository
}

func NewPaymentUsecase(pR repository.PaymentRepository, iR repository.InvoiceRepository, mR repository.MilestoneRuleRepository) PaymentService {
	return &paymentService{
		paymentRepo:  pR,
		invoiceRepo:  iR,
		milestoneRepo: mR,
	}
}

func (u *paymentService) ProcessSimulatedPayment(ctx context.Context, invoiceID, milestoneID, payerID, receiverID string, amount float64) (*model.Payment, error) {
	const platformFeePercentage = 0.10

	invoiceUUID := uuid.MustParse(invoiceID)
	payerUUID := uuid.MustParse(payerID)
	receiverUUID := uuid.MustParse(receiverID)

	platformFee := amount * platformFeePercentage
	amountCredited := amount - platformFee

	rzClient := pkg.NewRazorpayClient()
	order, err := rzClient.CreateOrder(amount, "receipt_"+invoiceID)
	if err != nil {
		return nil, err
	}

	payment := &model.Payment{
		ID:             uuid.New(),
		InvoiceID:      invoiceUUID,
		PayerID:        payerUUID,
		ReceiverID:     receiverUUID,
		AmountPaid:     amount,
		PlatformFee:    platformFee,
		AmountCredited: amountCredited,
		Status:         "completed",
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
		return nil, err
	}

	err = u.invoiceRepo.MarkPaid(ctx, invoiceUUID)
	if err != nil {
		return nil, err
	}

	return payment, nil
}
