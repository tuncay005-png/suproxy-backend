package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// BaseRepository provides common database operations
type BaseRepository struct {
	db *gorm.DB
}

func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

func (r *BaseRepository) Create(ctx context.Context, entity interface{}) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return fmt.Errorf("failed to create entity: %w", err)
	}
	return nil
}

func (r *BaseRepository) FindByID(ctx context.Context, id interface{}, entity interface{}) error {
	if err := r.db.WithContext(ctx).First(entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return fmt.Errorf("failed to find entity: %w", err)
	}
	return nil
}

func (r *BaseRepository) Update(ctx context.Context, entity interface{}) error {
	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		return fmt.Errorf("failed to update entity: %w", err)
	}
	return nil
}

func (r *BaseRepository) Delete(ctx context.Context, entity interface{}) error {
	if err := r.db.WithContext(ctx).Delete(entity).Error; err != nil {
		return fmt.Errorf("failed to delete entity: %w", err)
	}
	return nil
}

func (r *BaseRepository) FindAll(ctx context.Context, entities interface{}, conditions ...interface{}) error {
	query := r.db.WithContext(ctx)
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	if err := query.Find(entities).Error; err != nil {
		return fmt.Errorf("failed to find entities: %w", err)
	}
	return nil
}

func (r *BaseRepository) Count(ctx context.Context, model interface{}, conditions ...interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(model)
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count entities: %w", err)
	}
	return count, nil
}

func (r *BaseRepository) Exists(ctx context.Context, model interface{}, conditions ...interface{}) (bool, error) {
	count, err := r.Count(ctx, model, conditions...)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// DB returns the underlying gorm.DB instance for custom queries
func (r *BaseRepository) DB() *gorm.DB {
	return r.db
}
