package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/server"
	"gorm.io/gorm"
)

type serverRepository struct {
	db *gorm.DB
}

func NewServerRepository(db *gorm.DB) server.Repository {
	return &serverRepository{db: db}
}

func (r *serverRepository) Create(ctx context.Context, srv *server.Server) error {
	if srv == nil {
		return errors.New("server cannot be nil")
	}

	model := toServerModel(srv)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return server.ErrServerAlreadyExists
		}
		return err
	}
	return nil
}

func (r *serverRepository) FindByID(ctx context.Context, id uuid.UUID) (*server.Server, error) {
	var model ServerModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, server.ErrServerNotFound
		}
		return nil, err
	}
	return toDomainServer(&model), nil
}

func (r *serverRepository) FindByHostname(ctx context.Context, hostname string) (*server.Server, error) {
	var model ServerModel
	if err := r.db.WithContext(ctx).Where("hostname = ?", hostname).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, server.ErrServerNotFound
		}
		return nil, err
	}
	return toDomainServer(&model), nil
}

func (r *serverRepository) FindByCountry(ctx context.Context, country string) ([]*server.Server, error) {
	var models []ServerModel
	if err := r.db.WithContext(ctx).Where("country = ?", country).Order("name ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	servers := make([]*server.Server, 0, len(models))
	for _, model := range models {
		servers = append(servers, toDomainServer(&model))
	}
	return servers, nil
}

func (r *serverRepository) FindActive(ctx context.Context) ([]*server.Server, error) {
	var models []ServerModel
	if err := r.db.WithContext(ctx).Where("status = ?", "active").Order("country ASC, name ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	servers := make([]*server.Server, 0, len(models))
	for _, model := range models {
		servers = append(servers, toDomainServer(&model))
	}
	return servers, nil
}

func (r *serverRepository) FindPublicActive(ctx context.Context) ([]*server.Server, error) {
	var models []ServerModel
	if err := r.db.WithContext(ctx).
		Where("status = ? AND is_public = ?", "active", true).
		Order("country ASC, name ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	servers := make([]*server.Server, 0, len(models))
	for _, model := range models {
		servers = append(servers, toDomainServer(&model))
	}
	return servers, nil
}

func (r *serverRepository) Update(ctx context.Context, srv *server.Server) error {
	model := toServerModel(srv)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}
	return nil
}

func (r *serverRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&ServerModel{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *serverRepository) List(ctx context.Context, offset, limit int) ([]*server.Server, error) {
	var models []ServerModel
	if err := r.db.WithContext(ctx).
		Order("country ASC, name ASC").
		Offset(offset).
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}

	servers := make([]*server.Server, 0, len(models))
	for _, model := range models {
		servers = append(servers, toDomainServer(&model))
	}
	return servers, nil
}

func (r *serverRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&ServerModel{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
