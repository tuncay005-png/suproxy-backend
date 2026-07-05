package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/domain/user"
	"github.com/suproxy/backend/internal/infrastructure/jwt"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type RefreshTokenCommand struct {
	userRepo   user.Repository
	jwtManager *jwt.Manager
	logger     *logger.Logger
}

func NewRefreshTokenCommand(userRepo user.Repository, jwtManager *jwt.Manager, logger *logger.Logger) *RefreshTokenCommand {
	return &RefreshTokenCommand{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		logger:     logger,
	}
}

func (c *RefreshTokenCommand) Execute(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.TokenPair, error) {
	// Validate refresh token
	claims, err := c.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.logger.Warn("Invalid refresh token", "error", err)
		return nil, jwt.ErrInvalidToken
	}

	// Parse user ID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		c.logger.Error("Invalid user ID in token", "user_id", claims.UserID)
		return nil, jwt.ErrInvalidToken
	}

	// Verify user still exists and is active
	foundUser, err := c.userRepo.FindByID(ctx, userID)
	if err != nil {
		c.logger.Warn("User not found for refresh token", "user_id", userID)
		return nil, user.ErrUserNotFound
	}

	if !foundUser.IsActive() {
		c.logger.Warn("Inactive user attempted token refresh", "user_id", userID)
		return nil, user.ErrUserNotActive
	}

	// Generate new token pair
	accessToken, refreshToken, err := c.jwtManager.GenerateTokenPair(
		foundUser.ID.String(),
		foundUser.Email.String(),
	)
	if err != nil {
		c.logger.Error("Failed to generate new tokens", "error", err, "user_id", userID)
		return nil, err
	}

	c.logger.Info("Token refreshed successfully", "user_id", userID)

	return &dto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
