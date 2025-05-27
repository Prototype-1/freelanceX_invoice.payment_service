package service

import (
	"os"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
VerifyPayment(ctx context.Context, razorpayPaymentID, razorpayOrderID, razorpaySignature, invoiceID string) (bool, string, error)
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

func (u *paymentService) VerifyPayment(
	ctx context.Context,
	razorpayPaymentID, razorpayOrderID, razorpaySignature, invoiceID string,
) (bool, string, error) {
	secret := os.Getenv("RAZORPAY_SECRET")

	data := razorpayOrderID + "|" + razorpayPaymentID
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	if expectedSignature != razorpaySignature {
		return false, "Invalid signature", nil
	}

	invoiceUUID := uuid.MustParse(invoiceID)

	if err := u.invoiceRepo.MarkPaid(ctx, invoiceUUID); err != nil {
		return false, "Signature valid but DB update failed", err
	}

	return true, "Payment verified and invoice marked paid", nil
}
