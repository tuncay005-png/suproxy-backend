package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/node"
	"gorm.io/gorm"
)

type nodeRepository struct {
	db *gorm.DB
}

func NewNodeRepository(db *gorm.DB) node.Repository {
	return &nodeRepository{db: db}
}

func (r *nodeRepository) Create(ctx context.Context, n *node.Node) error {
	model := toNodeModel(n)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	return nil
}

func (r *nodeRepository) FindByID(ctx context.Context, id uuid.UUID) (*node.Node, error) {
	var model NodeModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, node.ErrNodeNotFound
		}
		return nil, err
	}
	return toDomainNode(&model), nil
}

func (r *nodeRepository) FindByServerID(ctx context.Context, serverID uuid.UUID) ([]*node.Node, error) {
	var models []NodeModel
	if err := r.db.WithContext(ctx).Where("server_id = ?", serverID).Order("port ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	nodes := make([]*node.Node, 0, len(models))
	for _, model := range models {
		nodes = append(nodes, toDomainNode(&model))
	}
	return nodes, nil
}

func (r *nodeRepository) FindHealthyByServerID(ctx context.Context, serverID uuid.UUID) ([]*node.Node, error) {
	var models []NodeModel
	if err := r.db.WithContext(ctx).
		Where("server_id = ? AND health_status = ?", serverID, "healthy").
		Order("current_users ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	nodes := make([]*node.Node, 0, len(models))
	for _, model := range models {
		nodes = append(nodes, toDomainNode(&model))
	}
	return nodes, nil
}

func (r *nodeRepository) FindAvailableNodes(ctx context.Context) ([]*node.Node, error) {
	var models []NodeModel
	if err := r.db.WithContext(ctx).
		Where("health_status = ? AND current_users < max_users", "healthy").
		Order("current_users ASC, latency_ms ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	nodes := make([]*node.Node, 0, len(models))
	for _, model := range models {
		n := toDomainNode(&model)
		if n.CanAcceptUser() {
			nodes = append(nodes, n)
		}
	}
	return nodes, nil
}

func (r *nodeRepository) Update(ctx context.Context, n *node.Node) error {
	model := toNodeModel(n)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}
	return nil
}

func (r *nodeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&NodeModel{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *nodeRepository) List(ctx context.Context, offset, limit int) ([]*node.Node, error) {
	var models []NodeModel
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}

	nodes := make([]*node.Node, 0, len(models))
	for _, model := range models {
		nodes = append(nodes, toDomainNode(&model))
	}
	return nodes, nil
}

func (r *nodeRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&NodeModel{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *nodeRepository) CountByServerID(ctx context.Context, serverID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&NodeModel{}).Where("server_id = ?", serverID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
