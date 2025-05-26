package handler

import (
	"context"
	"fmt"
	"time"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/model"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/service"
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
	MilestoneService   *service.MilestoneRuleService
	ProfileClient     profilepb.ProfileServiceClient
	TimeTrackerClient timepb.TimeLogServiceClient
}

func NewInvoiceHandler(repo repository.InvoiceRepository, milestoneSvc *service.MilestoneRuleService) *InvoiceHandler {
	return &InvoiceHandler{
		Repo:             repo,
		MilestoneService: milestoneSvc,
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

	if req.GetDateTo() != nil {
		t := req.GetDateTo().AsTime()
		dueDate = &t
	}

	switch req.GetType() {
	case invoicepb.InvoiceType_FIXED:
		amount = req.GetFixedAmount()

	case invoicepb.InvoiceType_HOURLY:
		// 1. Get hourly rate from profile
		profileResp, err := h.ProfileClient.GetProfile(ctx, &profilepb.GetProfileRequest{
			UserId: req.GetFreelancerId(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get freelancer profile: %v", err)
		}
		hourlyRate = float64(profileResp.HourlyRate)

		timeResp, err := h.TimeTrackerClient.GetTimeLogsByUser(ctx, &timepb.GetTimeLogsByUserRequest{
			UserId:    req.GetFreelancerId(),
			ProjectId: req.GetProjectId(),
			DateFrom:  req.GetDateFrom(),
			DateTo:    req.GetDateTo(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get tracked time: %v", err)
		}
		for _, log := range timeResp.Logs {
			hoursWorked += float64(log.Duration) / 60.0
		}

		amount = hourlyRate * hoursWorked

	case invoicepb.InvoiceType_MILESTONE:
		phase := req.GetMilestonePhase()
		if phase == "" {
			return nil, fmt.Errorf("milestone phase is required")
		}

		existing, err := h.MilestoneService.GetMilestoneByProjectIDAndPhase(projectID, phase)
		if err == nil && existing != nil {
			return nil, fmt.Errorf("invoice already exists for this milestone phase")
		}

		rule, err := h.MilestoneService.GetMilestoneByProjectIDAndPhase(projectID, phase)
		if err != nil {
			return nil, fmt.Errorf("milestone rule for phase '%s' not found: %v", phase, err)
		}
		amount = rule.Amount
	}

	invoice := &model.Invoice{
		ProjectID:      projectID,
		ClientID:       clientID,
		FreelancerID:   freelancerID,
		Type:           invoiceType,
		Amount:         amount,
		DueDate:        dueDate,
		HoursWorked:    hoursWorked,
		HourlyRate:     hourlyRate,
		Status:         "PENDING",
		MilestonePhase: req.GetMilestonePhase(), 
	}

	if err := h.Repo.CreateInvoice(ctx, invoice); err != nil {
		return nil, fmt.Errorf("failed to create invoice: %v", err)
	}

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
			MilestonePhase: invoice.MilestonePhase,
		},
	}
	if dueDate != nil {
		resp.Invoice.DueDate = timestamppb.New(*dueDate)
	}
	return resp, nil
}
