package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/model"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/repository"
)

type MilestoneRuleService struct {
	repo repository.MilestoneRuleRepository
}

func NewMilestoneRuleService(repo repository.MilestoneRuleRepository) *MilestoneRuleService {
	return &MilestoneRuleService{repo}
}

func (s *MilestoneRuleService) CreateMilestoneRule(rule *model.MilestoneRule) error {
	if rule.Amount < 0 {
		return errors.New("milestone amount must be non-negative")
	}
	return s.repo.Create(rule)
}

func (s *MilestoneRuleService) UpdateMilestoneRule(rule *model.MilestoneRule) error {
	if rule.Amount < 0 {
		return errors.New("milestone amount must be non-negative")
	}
existing, err := s.repo.GetByID(rule.ID)
	if err != nil {
		return err
	}
	existing.Phase = rule.Phase
	existing.Amount = rule.Amount
	if rule.DueDate != nil {
		existing.DueDate = rule.DueDate
	}

	return s.repo.Update(existing)
}

func (s *MilestoneRuleService) GetMilestonesByProjectID(projectID uuid.UUID) ([]model.MilestoneRule, error) {
	return s.repo.GetByProjectID(projectID)
}

func (s *MilestoneRuleService) GetMilestoneByProjectIDAndPhase(projectID uuid.UUID, phase string) (*model.MilestoneRule, error) {
	return s.repo.GetByProjectIDAndPhase(projectID, phase)
}

