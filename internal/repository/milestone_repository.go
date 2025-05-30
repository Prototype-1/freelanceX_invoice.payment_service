package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/Prototype-1/freelanceX_invoice.payment_service/internal/model"
)

// interface
type MilestoneRuleRepository interface {
	Create(rule *model.MilestoneRule) error
	Update(rule *model.MilestoneRule) error
	GetByID(id uuid.UUID) (*model.MilestoneRule, error)
	GetByProjectID(projectID uuid.UUID) ([]model.MilestoneRule, error)
	GetByProjectIDAndPhase(projectID uuid.UUID, phase string) (*model.MilestoneRule, error)
}

type milestoneRuleRepo struct {
	db *gorm.DB
}

func NewMilestoneRuleRepository(db *gorm.DB) MilestoneRuleRepository {
	return &milestoneRuleRepo{db}
}

func (r *milestoneRuleRepo) Create(rule *model.MilestoneRule) error {
	return r.db.Create(rule).Error
}

func (r *milestoneRuleRepo) Update(rule *model.MilestoneRule) error {
	return r.db.Save(rule).Error
}

func (r *milestoneRuleRepo) GetByProjectID(projectID uuid.UUID) ([]model.MilestoneRule, error) {
	var rules []model.MilestoneRule
	if err := r.db.Where("project_id = ?", projectID).Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

func (r *milestoneRuleRepo) GetByProjectIDAndPhase(projectID uuid.UUID, phase string) (*model.MilestoneRule, error) {
	var rule model.MilestoneRule
	if err := r.db.Where("project_id = ? AND phase = ?", projectID, phase).First(&rule).Error; err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *milestoneRuleRepo) GetByID(id uuid.UUID) (*model.MilestoneRule, error) {
	var rule model.MilestoneRule
	if err := r.db.First(&rule, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &rule, nil
}
