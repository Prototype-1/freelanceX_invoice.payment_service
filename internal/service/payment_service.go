package service

import (
	"context"
	"time"
	"github.com/google/uuid"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/model"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/repository"
)

type PaymentService interface {
	ProcessSimulatedPayment(ctx context.Context, invoiceID, milestoneID, payerID, receiverID string, amount float64) (*model.Payment, error)
}

type paymentUsecase struct {
	paymentRepo repository.PaymentRepository
	invoiceRepo repository.InvoiceRepository
	milestoneRepo repository.MilestoneRuleRepository
}

func NewPaymentUsecase(pR repository.PaymentRepository, iR repository.InvoiceRepository, mR repository.MilestoneRuleRepository) PaymentService {
	return &paymentUsecase{
		paymentRepo:  pR,
		invoiceRepo:  iR,
		milestoneRepo: mR,
	}
}

func (u *paymentUsecase) ProcessSimulatedPayment(ctx context.Context, invoiceID, milestoneID, payerID, receiverID string, amount float64) (*model.Payment, error) {
	const platformFeePercentage = 0.10

	invoiceUUID := uuid.MustParse(invoiceID)
	payerUUID := uuid.MustParse(payerID)
	receiverUUID := uuid.MustParse(receiverID)

	platformFee := amount * platformFeePercentage
	amountCredited := amount - platformFee

	payment := &model.Payment{
		ID:             uuid.New(),
		InvoiceID:      invoiceUUID,
		PayerID:        payerUUID,
		ReceiverID:     receiverUUID,
		AmountPaid:     amount,
		PlatformFee:    platformFee,
		AmountCredited: amountCredited,
		Status:         "completed",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if milestoneID != "" {
		mid := uuid.MustParse(milestoneID)
		payment.MilestoneID = &mid
	}

	err := u.paymentRepo.Create(ctx, payment)
	if err != nil {
		return nil, err
	}

	err = u.invoiceRepo.MarkPaid(ctx, invoiceUUID)
	if err != nil {
		return nil, err
	}

	return payment, nil
}
