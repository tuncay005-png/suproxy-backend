package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/xray"
	"gorm.io/gorm"
)

type xrayInstanceRepository struct {
	db *gorm.DB
}

func NewXrayInstanceRepository(db *gorm.DB) xray.XrayInstanceRepository {
	return &xrayInstanceRepository{db: db}
}

func (r *xrayInstanceRepository) Create(ctx context.Context, instance *xray.XrayInstance) error {
	model := toXrayInstanceModel(instance)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return xray.ErrInstanceAlreadyExists
		}
		return err
	}
	return nil
}

func (r *xrayInstanceRepository) FindByID(ctx context.Context, id uuid.UUID) (*xray.XrayInstance, error) {
	var model XrayInstanceModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xray.ErrInstanceNotFound
		}
		return nil, err
	}
	return toDomainXrayInstance(&model), nil
}

func (r *xrayInstanceRepository) FindByNodeID(ctx context.Context, nodeID uuid.UUID) (*xray.XrayInstance, error) {
	var model XrayInstanceModel
	if err := r.db.WithContext(ctx).Where("node_id = ?", nodeID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xray.ErrInstanceNotFound
		}
		return nil, err
	}
	return toDomainXrayInstance(&model), nil
}

func (r *xrayInstanceRepository) FindRunning(ctx context.Context) ([]*xray.XrayInstance, error) {
	var models []XrayInstanceModel
	if err := r.db.WithContext(ctx).Where("status = ?", "running").Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	instances := make([]*xray.XrayInstance, 0, len(models))
	for _, model := range models {
		instances = append(instances, toDomainXrayInstance(&model))
	}
	return instances, nil
}

func (r *xrayInstanceRepository) Update(ctx context.Context, instance *xray.XrayInstance) error {
	model := toXrayInstanceModel(instance)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}
	return nil
}

func (r *xrayInstanceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&XrayInstanceModel{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *xrayInstanceRepository) List(ctx context.Context, offset, limit int) ([]*xray.XrayInstance, error) {
	var models []XrayInstanceModel
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}

	instances := make([]*xray.XrayInstance, 0, len(models))
	for _, model := range models {
		instances = append(instances, toDomainXrayInstance(&model))
	}
	return instances, nil
}

func (r *xrayInstanceRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&XrayInstanceModel{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *xrayInstanceRepository) ListWithFilters(ctx context.Context, filters xray.XrayInstanceFilters) ([]*xray.XrayInstance, int64, error) {
	query := r.db.WithContext(ctx).Model(&XrayInstanceModel{})

	// Apply filters
	if filters.NodeID != nil {
		query = query.Where("node_id = ?", *filters.NodeID)
	}

	if filters.Status != nil {
		query = query.Where("status = ?", string(*filters.Status))
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
	var models []XrayInstanceModel
	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	// Convert to domain instances
	instances := make([]*xray.XrayInstance, 0, len(models))
	for _, model := range models {
		instances = append(instances, toDomainXrayInstance(&model))
	}

	return instances, total, nil
}
