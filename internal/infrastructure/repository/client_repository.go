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
	model := toClientModel(client)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return xray.ErrClientAlreadyExists
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
	if err := r.db.WithContext(ctx).Delete(&ClientModel{}, id).Error; err != nil {
		return err
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
