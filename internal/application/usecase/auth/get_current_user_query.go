package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	"github.com/suproxy/backend/internal/domain/user"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type GetCurrentUserQuery struct {
	userRepo user.Repository
	logger   *logger.Logger
}

func NewGetCurrentUserQuery(userRepo user.Repository, logger *logger.Logger) *GetCurrentUserQuery {
	return &GetCurrentUserQuery{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (q *GetCurrentUserQuery) Execute(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	// Find user
	foundUser, err := q.userRepo.FindByID(ctx, userID)
	if err != nil {
		q.logger.Warn("User not found", "user_id", userID)
		return nil, user.ErrUserNotFound
	}

	// Check if user is active
	if !foundUser.IsActive() {
		q.logger.Warn("Inactive user attempted to get profile", "user_id", userID)
		return nil, user.ErrUserNotActive
	}

	// Map to response
	return mapper.ToUserResponse(foundUser), nil
}
