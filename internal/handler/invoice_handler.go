package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/model"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/repository"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/service"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/kafka"
	invoicepb "github.com/Prototype-1/freelanceX_invoice.payment_service/proto/invoice_service"
	timepb "github.com/Prototype-1/freelanceX_invoice.payment_service/proto/timeTracker_service"
	profilepb "github.com/Prototype-1/freelanceX_invoice.payment_service/proto/user_service"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type InvoiceHandler struct {
	invoicepb.UnimplementedInvoiceServiceServer
	Service           service.InvoiceUsecase
	Repo              repository.InvoiceRepository
	MilestoneService  *service.MilestoneRuleService
	ProfileClient     profilepb.ProfileServiceClient
	TimeTrackerClient timepb.TimeLogServiceClient
	KafkaBroker       string
	KafkaTopic        string
}

func NewInvoiceHandler(
	repo repository.InvoiceRepository,
	svc service.InvoiceUsecase,
	milestoneSvc *service.MilestoneRuleService,
	kafkaBroker, kafkaTopic string,
) *InvoiceHandler {
	return &InvoiceHandler{
		Service:          svc, 
		Repo:             repo,
		MilestoneService: milestoneSvc,
		KafkaBroker:      kafkaBroker,
		KafkaTopic:       kafkaTopic,
	}
}

func extractRole(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("missing metadata")
	}
	roles := md.Get("role")
	if len(roles) == 0 {
		return "", errors.New("role not found in metadata")
	}
	return roles[0], nil
}

func (h *InvoiceHandler) CreateInvoice(ctx context.Context, req *invoicepb.CreateInvoiceRequest) (*invoicepb.InvoiceResponse, error) {

	role, err := extractRole(ctx)
	if err != nil {
		return nil, err
	}
	if role != "freelancer" {
		return nil, errors.New("unauthorized: only freelancer can create invoice status")
	}

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
		md := metadata.Pairs(
			"role", role,
			"user_id", req.GetFreelancerId(),
		)
		outgoingCtx := metadata.NewOutgoingContext(ctx, md)
		profileResp, err := h.ProfileClient.GetProfile(outgoingCtx, &profilepb.GetProfileRequest{
			UserId: req.GetFreelancerId(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get freelancer profile: %v", err)
		}
		hourlyRate = float64(profileResp.HourlyRate)

		timeResp, err := h.TimeTrackerClient.GetTimeLogsByUser(outgoingCtx, &timepb.GetTimeLogsByUserRequest{
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

	event := kafka.InvoiceEvent{
		InvoiceID: invoice.ID.String(),
		ClientID:  invoice.ClientID.String(),
		EventType: "invoice_created",
	}

	if err := kafka.ProduceInvoiceEvent("kafka:9092", "invoice-events", event); err != nil {
		fmt.Printf("failed to send Kafka event: %v\n", err)
	}

	resp := &invoicepb.InvoiceResponse{
		Invoice: &invoicepb.Invoice{
			InvoiceId:      invoice.ID.String(),
			FreelancerId:   invoice.FreelancerID.String(),
			ClientId:       invoice.ClientID.String(),
			ProjectId:      invoice.ProjectID.String(),
			Type:           req.GetType(),
			Amount:         invoice.Amount,
			HourlyRate:     invoice.HourlyRate,
			HoursWorked:    invoice.HoursWorked,
			Status:         invoicepb.InvoiceStatus_PENDING,
			IssuedAt:       timestamppb.New(invoice.CreatedAt),
			MilestonePhase: invoice.MilestonePhase,
		},
	}
	if dueDate != nil {
		resp.Invoice.DueDate = timestamppb.New(*dueDate)
	}
	return resp, nil
}

func (h *InvoiceHandler) GetInvoice(ctx context.Context, req *invoicepb.GetInvoiceRequest) (*invoicepb.InvoiceResponse, error) {
	invoice, err := h.Repo.GetInvoiceByID(ctx, req.GetInvoiceId())
	if err != nil {
		return nil, err
	}
	return &invoicepb.InvoiceResponse{
		Invoice: model.ToProto(invoice),
	}, nil
}

func (h *InvoiceHandler) GetInvoicesByUser(ctx context.Context, req *invoicepb.GetInvoicesByUserRequest) (*invoicepb.InvoicesResponse, error) {
	role, err := extractRole(ctx)
	if err != nil {
		return nil, err
	}
	if role != req.GetRole() {
		return nil, errors.New("unauthorized: role mismatch")
	}

	filter := &repository.InvoiceFilter{}
	if role == "freelancer" {
		id := req.GetUserId()
		filter.FreelancerID = &id
	} else if role == "client" {
		id := req.GetUserId()
		filter.ClientID = &id
	} else {
		return nil, errors.New("invalid role")
	}

	invoices, err := h.Repo.ListInvoices(ctx, filter)
	if err != nil {
		return nil, err
	}
	var protoInvoices []*invoicepb.Invoice
	for _, inv := range invoices {
		protoInvoices = append(protoInvoices, model.ToProto(inv))
	}
	return &invoicepb.InvoicesResponse{Invoices: protoInvoices}, nil
}

func (h *InvoiceHandler) GetInvoicesByProject(ctx context.Context, req *invoicepb.GetInvoicesByProjectRequest) (*invoicepb.InvoicesResponse, error) {
	id := req.GetProjectId()
	filter := &repository.InvoiceFilter{ProjectID: &id}
	invoices, err := h.Repo.ListInvoices(ctx, filter)
	if err != nil {
		return nil, err
	}
	var protoInvoices []*invoicepb.Invoice
	for _, inv := range invoices {
		protoInvoices = append(protoInvoices, model.ToProto(inv))
	}
	return &invoicepb.InvoicesResponse{Invoices: protoInvoices}, nil
}

func (h *InvoiceHandler) UpdateInvoiceStatus(ctx context.Context, req *invoicepb.UpdateInvoiceStatusRequest) (*invoicepb.InvoiceResponse, error) {
	role, err := extractRole(ctx)
	if err != nil {
		return nil, err
	}
	if role != "client" {
		return nil, errors.New("unauthorized: only clients can update invoice status")
	}

	statusStr := req.GetStatus().String()
	err = h.Repo.UpdateStatus(ctx, req.GetInvoiceId(), statusStr)
	if err != nil {
		return nil, err
	}

	invoice, err := h.Repo.GetInvoiceByID(ctx, req.GetInvoiceId())
	if err != nil {
		return nil, err
	}
	return &invoicepb.InvoiceResponse{
		Invoice: model.ToProto(invoice),
	}, nil
}

func (h *InvoiceHandler) GetInvoicePDF(ctx context.Context, req *invoicepb.GetInvoicePDFRequest) (*invoicepb.GetInvoicePDFResponse, error) {
	pdfData, err := h.Service.GetInvoicePDF(ctx, req.GetInvoiceId())
	if err != nil {
		return nil, err
	}
	return &invoicepb.GetInvoicePDFResponse{
		PdfData: pdfData,
	}, nil
}
