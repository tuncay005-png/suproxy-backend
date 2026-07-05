package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/session"
	"gorm.io/gorm"
)

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) session.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(ctx context.Context, token *session.RefreshToken) error {
	model := toRefreshTokenModel(token)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *refreshTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*session.RefreshToken, error) {
	var model RefreshTokenModel
	if err := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, session.ErrSessionNotFound
		}
		return nil, err
	}
	return toDomainRefreshToken(&model), nil
}

func (r *refreshTokenRepository) FindActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*session.RefreshToken, error) {
	var models []RefreshTokenModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_revoked = ? AND expires_at > ?", userID, false, time.Now().UTC()).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	tokens := make([]*session.RefreshToken, 0, len(models))
	for _, model := range models {
		tokens = append(tokens, toDomainRefreshToken(&model))
	}
	return tokens, nil
}

func (r *refreshTokenRepository) RevokeByID(ctx context.Context, id uuid.UUID) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).
		Model(&RefreshTokenModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_revoked": true,
			"revoked_at": now,
		}).Error
}

func (r *refreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).
		Model(&RefreshTokenModel{}).
		Where("user_id = ? AND is_revoked = ?", userID, false).
		Updates(map[string]interface{}{
			"is_revoked": true,
			"revoked_at": now,
		}).Error
}

func (r *refreshTokenRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now().UTC()).
		Delete(&RefreshTokenModel{}).Error
}
