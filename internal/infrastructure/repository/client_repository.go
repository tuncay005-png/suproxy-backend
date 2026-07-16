package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/xray"
	"gorm.io/gorm"
)

type clientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) xray.ClientRepository {
	return &clientRepository{db: db}
}

func (r *clientRepository) Create(ctx context.Context, client *xray.Client) error {
	if client == nil {
		return errors.New("client cannot be nil")
	}

	model := toClientModel(client)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return xray.ErrClientAlreadyExists
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("inbound or user not found")
		}
		return err
	}
	return nil
}

func (r *clientRepository) FindByID(ctx context.Context, id uuid.UUID) (*xray.Client, error) {
	var model ClientModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xray.ErrClientNotFound
		}
		return nil, err
	}
	return toDomainClient(&model), nil
}

func (r *clientRepository) FindByInboundID(ctx context.Context, inboundID uuid.UUID) ([]*xray.Client, error) {
	var models []ClientModel
	if err := r.db.WithContext(ctx).Where("inbound_id = ?", inboundID).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	clients := make([]*xray.Client, 0, len(models))
	for _, model := range models {
		clients = append(clients, toDomainClient(&model))
	}
	return clients, nil
}

func (r *clientRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*xray.Client, error) {
	var models []ClientModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	clients := make([]*xray.Client, 0, len(models))
	for _, model := range models {
		clients = append(clients, toDomainClient(&model))
	}
	return clients, nil
}

func (r *clientRepository) FindByUUID(ctx context.Context, clientUUID string) (*xray.Client, error) {
	var model ClientModel
	if err := r.db.WithContext(ctx).Where("uuid = ?", clientUUID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xray.ErrClientNotFound
		}
		return nil, err
	}
	return toDomainClient(&model), nil
}

func (r *clientRepository) FindEnabledByInboundID(ctx context.Context, inboundID uuid.UUID) ([]*xray.Client, error) {
	var models []ClientModel
	if err := r.db.WithContext(ctx).
		Where("inbound_id = ? AND enabled = ?", inboundID, true).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	clients := make([]*xray.Client, 0, len(models))
	for _, model := range models {
		clients = append(clients, toDomainClient(&model))
	}
	return clients, nil
}

func (r *clientRepository) Update(ctx context.Context, client *xray.Client) error {
	model := toClientModel(client)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}
	return nil
}

func (r *clientRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&ClientModel{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return xray.ErrClientNotFound
	}
	return nil
}

func (r *clientRepository) List(ctx context.Context, offset, limit int) ([]*xray.Client, error) {
	var models []ClientModel
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}

	clients := make([]*xray.Client, 0, len(models))
	for _, model := range models {
		clients = append(clients, toDomainClient(&model))
	}
	return clients, nil
}

func (r *clientRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&ClientModel{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *clientRepository) ListWithFilters(ctx context.Context, filters xray.ClientFilters) ([]*xray.Client, int64, error) {
	query := r.db.WithContext(ctx).Model(&ClientModel{})

	// Apply filters
	if filters.InboundID != nil {
		query = query.Where("inbound_id = ?", *filters.InboundID)
	}

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
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
	var models []ClientModel
	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	// Convert to domain clients
	clients := make([]*xray.Client, 0, len(models))
	for _, model := range models {
		clients = append(clients, toDomainClient(&model))
	}

	return clients, total, nil
}
