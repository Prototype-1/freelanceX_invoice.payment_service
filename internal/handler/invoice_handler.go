package handler

import (
	"context"
	"fmt"
	"time"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/model"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/repository"
	invoicepb "github.com/Prototype-1/freelanceX_invoice.payment_service/proto/invoice_service"
	profilepb "github.com/Prototype-1/freelanceX_invoice.payment_service/proto/user_service"
	timepb "github.com/Prototype-1/freelanceX_invoice.payment_service/proto/timeTracker_service"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type InvoiceHandler struct {
	invoicepb.UnimplementedInvoiceServiceServer
	Repo repository.InvoiceRepository
	ProfileClient     profilepb.ProfileServiceClient
	TimeTrackerClient timepb.TimeLogServiceClient
}

func NewInvoiceHandler(repo repository.InvoiceRepository) *InvoiceHandler {
	return &InvoiceHandler{
		Repo: repo,
	}
}

func (h *InvoiceHandler) CreateInvoice(ctx context.Context, req *invoicepb.CreateInvoiceRequest) (*invoicepb.InvoiceResponse, error) {
	projectID, _ := uuid.Parse(req.GetProjectId())
	clientID, _ := uuid.Parse(req.GetClientId())
	freelancerID, _ := uuid.Parse(req.GetFreelancerId())

	invoiceType := req.GetType().String()
	var amount float64
	var dueDate *time.Time
	var hoursWorked float64
	var hourlyRate float64

	// === Handle DueDate from date_to if provided ===
	if req.GetDateTo() != nil {
		t := req.GetDateTo().AsTime()
		dueDate = &t
	}

	// === Compute amount ===
	switch req.GetType() {
	case invoicepb.InvoiceType_FIXED:
		amount = req.GetFixedAmount()

	case invoicepb.InvoiceType_HOURLY:
		// Get hourly rate from ProfileService
		profileResp, err := h.ProfileClient.GetProfile(ctx, &profilepb.GetProfileRequest{
			UserId: req.GetFreelancerId(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get freelancer profile: %v", err)
		}
		hourlyRate = float64(profileResp.HourlyRate)

		// Get hours worked from TimeTrackerService
		timeResp, err := h.TimeTrackerClient.GetTimeLogsByUser(ctx, &timepb.GetTimeLogsByUserRequest{
			UserId: req.GetFreelancerId(),
			ProjectId:    req.GetProjectId(),
			DateFrom:     req.GetDateFrom(),
			DateTo:       req.GetDateTo(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get tracked time: %v", err)
		}
	for _, log := range timeResp.Logs {
		hoursWorked += float64(log.Duration) / 60.0 // assuming Duration is in minutes
	}

		// Compute total
		amount = hourlyRate * hoursWorked
	}

	// === Create Invoice ===
	invoice := &model.Invoice{
		ProjectID:     projectID,
		ClientID:      clientID,
		FreelancerID:  freelancerID,
		Type:          invoiceType,
		Amount:        amount,
		DueDate:       dueDate,
		HoursWorked:   hoursWorked,
		HourlyRate:    hourlyRate,
		Status:        "PENDING",
	}

	if err := h.Repo.CreateInvoice(ctx, invoice); err != nil {
		return nil, fmt.Errorf("failed to create invoice: %v", err)
	}

	// === Return Response ===
	resp := &invoicepb.InvoiceResponse{
		Invoice: &invoicepb.Invoice{
			InvoiceId:     invoice.ID.String(),
			FreelancerId:  invoice.FreelancerID.String(),
			ClientId:      invoice.ClientID.String(),
			ProjectId:     invoice.ProjectID.String(),
			Type:          req.GetType(),
			Amount:        invoice.Amount,
			HourlyRate:    invoice.HourlyRate,
			HoursWorked:   invoice.HoursWorked,
			Status:        invoicepb.InvoiceStatus_PENDING,
			IssuedAt:      timestamppb.New(invoice.CreatedAt),
		},
	}
	if dueDate != nil {
		resp.Invoice.DueDate = timestamppb.New(*dueDate)
	}
	return resp, nil
}

