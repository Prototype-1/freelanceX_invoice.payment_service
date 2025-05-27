package handler

import (
"context"
"google.golang.org/grpc/metadata"
"google.golang.org/grpc/status"	
"google.golang.org/grpc/codes"
"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/service"
paymentpb "github.com/Prototype-1/freelanceX_invoice.payment_service/proto/payment"
)

type PaymentServiceServer struct {
	paymentpb.UnimplementedPaymentServiceServer
	service service.PaymentService
}

func NewPaymentServiceServer(u service.PaymentService) *PaymentServiceServer {
	return &PaymentServiceServer{service: u}
}

func (s *PaymentServiceServer) CreatePaymentOrder(
	ctx context.Context,
	req *paymentpb.CreatePaymentOrderRequest,
) (*paymentpb.CreatePaymentOrderResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	roles := md.Get("role")
	if len(roles) == 0 || roles[0] != "client" {
		return nil, status.Error(codes.PermissionDenied, "only clients can initiate payment")
	}

	payment, order, err := s.service.CreatePaymentOrder(
		ctx,
		req.InvoiceId,
		req.MilestoneId,
		req.PayerId,
		req.ReceiverId,
		req.Amount,
	)
	if err != nil {
		return nil, err
	}

	return &paymentpb.CreatePaymentOrderResponse{
		PaymentId:       payment.ID.String(),
		RazorpayOrderId: order.ID,
		Amount:          req.Amount,
		Currency:        "INR", // or order.Currency
		InvoiceId:       req.InvoiceId,
	}, nil
}