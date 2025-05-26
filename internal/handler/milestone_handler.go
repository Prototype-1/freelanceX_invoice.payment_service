package handler

import (
	"context"
	"errors"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/service"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/model"
	milestonePb "github.com/Prototype-1/freelanceX_invoice.payment_service/proto/milestone"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
		"google.golang.org/grpc/metadata"
	"github.com/google/uuid"
)

type MilestoneRuleHandler struct {
	milestonePb.UnimplementedMilestoneRuleServiceServer
	svc *service.MilestoneRuleService
}

func NewMilestoneRuleHandler(svc *service.MilestoneRuleService) *MilestoneRuleHandler {
	return &MilestoneRuleHandler{svc: svc}
}

func extractRoles(ctx context.Context) (string, error) {
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

func (h *MilestoneRuleHandler) CreateMilestoneRule(ctx context.Context, req *milestonePb.CreateMilestoneRuleRequest) (*milestonePb.MilestoneRule, error) {
	role, err := extractRoles(ctx)
		if err != nil {
		return nil, err
	}
	if role != "client" {
		return nil, status.Error(codes.PermissionDenied, "only clients can create milestone rules")
	}

	projectID, err := uuid.Parse(req.GetProjectId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid project ID")
	}

	rule := &model.MilestoneRule{
		ProjectID: projectID,
		Phase:     req.GetPhase(),
		Amount:    req.GetAmount(),
	}

	if req.DueDate != nil {
		due := req.DueDate.AsTime()
		rule.DueDate = &due
	}

	if err := h.svc.CreateMilestoneRule(rule); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toProto(rule), nil
}

func (h *MilestoneRuleHandler) UpdateMilestoneRule(ctx context.Context, req *milestonePb.UpdateMilestoneRuleRequest) (*milestonePb.MilestoneRule, error) {
	role, err := extractRoles(ctx)
		if err != nil {
		return nil, err
	}
	if role != "client" {
		return nil, status.Error(codes.PermissionDenied, "only clients can update milestone rules")
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid ID")
	}

	rule := &model.MilestoneRule{
		ID:     id,
		Phase:  req.GetPhase(),
		Amount: req.GetAmount(),
	}

	if req.DueDate != nil {
		due := req.DueDate.AsTime()
		rule.DueDate = &due
	}

	if err := h.svc.UpdateMilestoneRule(rule); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toProto(rule), nil
}

func (h *MilestoneRuleHandler) GetMilestonesByProjectID(ctx context.Context, req *milestonePb.GetMilestonesByProjectIDRequest) (*milestonePb.GetMilestonesByProjectIDResponse, error) {
	projectID, err := uuid.Parse(req.GetProjectId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid project ID")
	}

	rules, err := h.svc.GetMilestonesByProjectID(projectID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &milestonePb.GetMilestonesByProjectIDResponse{}
	for _, r := range rules {
		resp.Milestones = append(resp.Milestones, toProto(&r))
	}
	return resp, nil
}

func (h *MilestoneRuleHandler) GetMilestoneByProjectIDAndPhase(ctx context.Context, req *milestonePb.GetMilestoneByProjectIDAndPhaseRequest) (*milestonePb.GetMilestoneByProjectIDAndPhaseResponse, error) {
	projectID, err := uuid.Parse(req.GetProjectId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid project ID")
	}

	rule, err := h.svc.GetMilestoneByProjectIDAndPhase(projectID, req.GetPhase())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &milestonePb.GetMilestoneByProjectIDAndPhaseResponse{
		Milestone: toProto(rule),
	}, nil
}

func toProto(m *model.MilestoneRule) *milestonePb.MilestoneRule {
	pb := &milestonePb.MilestoneRule{
		Id:        m.ID.String(),
		ProjectId: m.ProjectID.String(),
		Phase:     m.Phase,
		Amount:    m.Amount,
		CreatedAt: timestamppb.New(m.CreatedAt),
	}
	if m.DueDate != nil {
		pb.DueDate = timestamppb.New(*m.DueDate)
	}
	return pb
}
