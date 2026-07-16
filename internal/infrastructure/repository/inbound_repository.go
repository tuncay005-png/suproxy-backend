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
	if inbound == nil {
		return errors.New("inbound cannot be nil")
	}

	model := toInboundModel(inbound)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return xray.ErrInboundAlreadyExists
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("xray instance not found")
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
	if inbound == nil {
		return errors.New("inbound cannot be nil")
	}

	model := toInboundModel(inbound)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}
	return nil
}

func (r *inboundRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&InboundModel{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return xray.ErrInboundNotFound
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

func (r *inboundRepository) ListWithFilters(ctx context.Context, filters xray.InboundFilters) ([]*xray.Inbound, int64, error) {
	query := r.db.WithContext(ctx).Model(&InboundModel{})

	// Apply filters
	if filters.InstanceID != nil {
		query = query.Where("xray_instance_id = ?", *filters.InstanceID)
	}

	if filters.Protocol != nil {
		query = query.Where("protocol = ?", string(*filters.Protocol))
	}

	if filters.Enabled != nil {
		query = query.Where("enabled = ?", *filters.Enabled)
	}

	// Count total before pagination
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortField := "created_at"
	if filters.SortBy != "" {
		sortField = filters.SortBy
	}

	sortOrder := "DESC"
	if filters.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	query = query.Order(sortField + " " + sortOrder)

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	// Execute query
	var models []InboundModel
	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	// Convert to domain inbounds
	inbounds := make([]*xray.Inbound, 0, len(models))
	for _, model := range models {
		inbounds = append(inbounds, toDomainInbound(&model))
	}

	return inbounds, total, nil
}
