package handler

import (
	"log"
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
	log.Println("CreatePaymentOrder called") 

md, ok := metadata.FromIncomingContext(ctx)
log.Printf("metadata ok: %v, md: %+v\n", ok, md)

if !ok {
	return nil, status.Error(codes.Unauthenticated, "missing metadata")
}

roles := md.Get("role")
log.Printf("Roles: %+v\n", roles)

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
		Currency:        "INR", 
		InvoiceId:       req.InvoiceId,
	}, nil
}

func (s *PaymentServiceServer) VerifyPayment(
	ctx context.Context,
	req *paymentpb.VerifyPaymentRequest,
) (*paymentpb.VerifyPaymentResponse, error) {
	valid, msg, err := s.service.VerifyPayment(
		ctx,
		req.RazorpayPaymentId,
		req.RazorpayOrderId,
		req.RazorpaySignature,
		req.InvoiceId,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "verification failed: %v", err)
	}

	return &paymentpb.VerifyPaymentResponse{
		Valid:   valid,
		Message: msg,
	}, nil
}
