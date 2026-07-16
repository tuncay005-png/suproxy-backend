package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/subscription"
	"gorm.io/gorm"
)

type planRepository struct {
	db *gorm.DB
}

func NewPlanRepository(db *gorm.DB) subscription.PlanRepository {
	return &planRepository{db: db}
}

func (r *planRepository) Create(ctx context.Context, plan *subscription.Plan) error {
	model := toPlanModel(plan)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return subscription.ErrPlanAlreadyExists
		}
		return err
	}
	return nil
}

func (r *planRepository) FindByID(ctx context.Context, id uuid.UUID) (*subscription.Plan, error) {
	var model PlanModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, subscription.ErrPlanNotFound
		}
		return nil, err
	}
	return toDomainPlan(&model)
}

func (r *planRepository) FindByName(ctx context.Context, name string) (*subscription.Plan, error) {
	var model PlanModel
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, subscription.ErrPlanNotFound
		}
		return nil, err
	}
	return toDomainPlan(&model)
}

func (r *planRepository) FindActive(ctx context.Context) ([]*subscription.Plan, error) {
	var models []PlanModel
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("price ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	plans := make([]*subscription.Plan, 0, len(models))
	for _, model := range models {
		plan, err := toDomainPlan(&model)
		if err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}
	return plans, nil
}

func (r *planRepository) Update(ctx context.Context, plan *subscription.Plan) error {
	model := toPlanModel(plan)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}
	return nil
}

func (r *planRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&PlanModel{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return subscription.ErrPlanNotFound
	}
	return nil
}

func (r *planRepository) List(ctx context.Context) ([]*subscription.Plan, error) {
	var models []PlanModel
	if err := r.db.WithContext(ctx).Order("price ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	plans := make([]*subscription.Plan, 0, len(models))
	for _, model := range models {
		plan, err := toDomainPlan(&model)
		if err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}
	return plans, nil
}

func (r *planRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&PlanModel{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
