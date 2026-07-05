package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/subscription"
	"gorm.io/gorm"
)

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) subscription.SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) Create(ctx context.Context, sub *subscription.Subscription) error {
	model := toSubscriptionModel(sub)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return subscription.ErrSubscriptionAlreadyExists
		}
		return err
	}
	return nil
}

func (r *subscriptionRepository) FindByID(ctx context.Context, id uuid.UUID) (*subscription.Subscription, error) {
	var model SubscriptionModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, subscription.ErrSubscriptionNotFound
		}
		return nil, err
	}
	return toDomainSubscription(&model), nil
}

func (r *subscriptionRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*subscription.Subscription, error) {
	var model SubscriptionModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, subscription.ErrSubscriptionNotFound
		}
		return nil, err
	}
	return toDomainSubscription(&model), nil
}

func (r *subscriptionRepository) FindActiveByUserID(ctx context.Context, userID uuid.UUID) (*subscription.Subscription, error) {
	var model SubscriptionModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, "active").
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, subscription.ErrSubscriptionNotFound
		}
		return nil, err
	}
	return toDomainSubscription(&model), nil
}

func (r *subscriptionRepository) Update(ctx context.Context, sub *subscription.Subscription) error {
	model := toSubscriptionModel(sub)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}
	return nil
}

func (r *subscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&SubscriptionModel{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *subscriptionRepository) List(ctx context.Context, offset, limit int) ([]*subscription.Subscription, error) {
	var models []SubscriptionModel
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}

	subs := make([]*subscription.Subscription, 0, len(models))
	for _, model := range models {
		subs = append(subs, toDomainSubscription(&model))
	}
	return subs, nil
}

func (r *subscriptionRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&SubscriptionModel{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
