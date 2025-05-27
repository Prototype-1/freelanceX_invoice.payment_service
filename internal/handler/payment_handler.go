package handler

import (
"context"
"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/service"
pb "github.com/Prototype-1/freelanceX_invoice.payment_service/proto/payment"
)

type PaymentServiceServer struct {
	pb.UnimplementedPaymentServiceServer
	usecase service.PaymentService
}

func NewPaymentServiceServer(u service.PaymentService) *PaymentServiceServer {
	return &PaymentServiceServer{usecase: u}
}

func (s *PaymentServiceServer) SimulatePayment(ctx context.Context, req *pb.SimulatePaymentRequest) (*pb.SimulatePaymentResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	roles := md.Get("role")
	if len(roles) == 0 || roles[0] != "client" {
		return nil, status.Error(codes.PermissionDenied, "only clients can initiate payment")
	}

	
	payment, err := s.usecase.ProcessSimulatedPayment(
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

	return &pb.SimulatePaymentResponse{
		PaymentId:      payment.ID.String(),
		AmountPaid:     payment.AmountPaid,
		PlatformFee:    payment.PlatformFee,
		AmountCredited: payment.AmountCredited,
		Status:         payment.Status,
	}, nil
}
