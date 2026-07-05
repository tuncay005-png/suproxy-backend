package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/xray"
	"gorm.io/gorm"
)

type inboundRepository struct {
	db *gorm.DB
}

func NewInboundRepository(db *gorm.DB) xray.InboundRepository {
	return &inboundRepository{db: db}
}

func (r *inboundRepository) Create(ctx context.Context, inbound *xray.Inbound) error {
	model := toInboundModel(inbound)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return xray.ErrInboundAlreadyExists
		}
		return err
	}
	return nil
}

func (r *inboundRepository) FindByID(ctx context.Context, id uuid.UUID) (*xray.Inbound, error) {
	var model InboundModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xray.ErrInboundNotFound
		}
		return nil, err
	}
	return toDomainInbound(&model), nil
}

func (r *inboundRepository) FindByInstanceID(ctx context.Context, instanceID uuid.UUID) ([]*xray.Inbound, error) {
	var models []InboundModel
	if err := r.db.WithContext(ctx).Where("xray_instance_id = ?", instanceID).Order("port ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	inbounds := make([]*xray.Inbound, 0, len(models))
	for _, model := range models {
		inbounds = append(inbounds, toDomainInbound(&model))
	}
	return inbounds, nil
}

func (r *inboundRepository) FindEnabledByInstanceID(ctx context.Context, instanceID uuid.UUID) ([]*xray.Inbound, error) {
	var models []InboundModel
	if err := r.db.WithContext(ctx).
		Where("xray_instance_id = ? AND enabled = ?", instanceID, true).
		Order("port ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	inbounds := make([]*xray.Inbound, 0, len(models))
	for _, model := range models {
		inbounds = append(inbounds, toDomainInbound(&model))
	}
	return inbounds, nil
}

func (r *inboundRepository) Update(ctx context.Context, inbound *xray.Inbound) error {
	model := toInboundModel(inbound)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}
	return nil
}

func (r *inboundRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&InboundModel{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *inboundRepository) List(ctx context.Context, offset, limit int) ([]*xray.Inbound, error) {
	var models []InboundModel
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}

	inbounds := make([]*xray.Inbound, 0, len(models))
	for _, model := range models {
		inbounds = append(inbounds, toDomainInbound(&model))
	}
	return inbounds, nil
}

func (r *inboundRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&InboundModel{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
