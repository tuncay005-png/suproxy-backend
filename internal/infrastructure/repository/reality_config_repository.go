package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/xray"
	"gorm.io/gorm"
)

type realityConfigRepository struct {
	db *gorm.DB
}

func NewRealityConfigRepository(db *gorm.DB) xray.RealityConfigRepository {
	return &realityConfigRepository{db: db}
}

func (r *realityConfigRepository) Create(ctx context.Context, config *xray.RealityConfig) error {
	model := toRealityConfigModel(config)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return xray.ErrRealityConfigExists
		}
		return err
	}
	return nil
}

func (r *realityConfigRepository) FindByID(ctx context.Context, id uuid.UUID) (*xray.RealityConfig, error) {
	var model RealityConfigModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xray.ErrRealityConfigNotFound
		}
		return nil, err
	}
	return toDomainRealityConfig(&model), nil
}

func (r *realityConfigRepository) FindByInboundID(ctx context.Context, inboundID uuid.UUID) (*xray.RealityConfig, error) {
	var model RealityConfigModel
	if err := r.db.WithContext(ctx).Where("inbound_id = ?", inboundID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xray.ErrRealityConfigNotFound
		}
		return nil, err
	}
	return toDomainRealityConfig(&model), nil
}

func (r *realityConfigRepository) Update(ctx context.Context, config *xray.RealityConfig) error {
	model := toRealityConfigModel(config)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}
	return nil
}

func (r *realityConfigRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&RealityConfigModel{}, id).Error; err != nil {
		return err
	}
	return nil
}
