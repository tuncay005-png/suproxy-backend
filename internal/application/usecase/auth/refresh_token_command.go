package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/domain/audit"
	"github.com/suproxy/backend/internal/domain/session"
	"github.com/suproxy/backend/internal/domain/user"
	"github.com/suproxy/backend/internal/infrastructure/jwt"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type RefreshTokenCommand struct {
	userRepo         user.Repository
	refreshTokenRepo session.RefreshTokenRepository
	auditRepo        audit.Repository
	jwtManager       *jwt.Manager
	logger           *logger.Logger
}

func NewRefreshTokenCommand(
	userRepo user.Repository,
	refreshTokenRepo session.RefreshTokenRepository,
	auditRepo audit.Repository,
	jwtManager *jwt.Manager,
	logger *logger.Logger,
) *RefreshTokenCommand {
	return &RefreshTokenCommand{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		auditRepo:        auditRepo,
		jwtManager:       jwtManager,
		logger:           logger,
	}
}

func (c *RefreshTokenCommand) Execute(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.TokenPair, error) {
	// Hash the refresh token
	tokenHash := hashToken(req.RefreshToken)

	// Find token in database
	storedToken, err := c.refreshTokenRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		c.logger.Warn("Refresh token not found in database", "error", err)
		return nil, jwt.ErrInvalidToken
	}

	// Validate token
	if !storedToken.IsValid() {
		c.logger.Warn("Invalid or revoked refresh token", "token_id", storedToken.ID)
		return nil, jwt.ErrInvalidToken
	}

	// Verify user still exists and is active
	foundUser, err := c.userRepo.FindByID(ctx, storedToken.UserID)
	if err != nil {
		c.logger.Warn("User not found for refresh token", "user_id", storedToken.UserID)
		return nil, user.ErrUserNotFound
	}

	if !foundUser.IsActive() {
		c.logger.Warn("Inactive user attempted token refresh", "user_id", storedToken.UserID)
		return nil, user.ErrUserNotActive
	}

	// Revoke old refresh token (rotation)
	if err := c.refreshTokenRepo.RevokeByID(ctx, storedToken.ID); err != nil {
		c.logger.Error("Failed to revoke old token", "error", err)
		return nil, err
	}

	// Generate new token pair
	accessToken, newRefreshToken, err := c.jwtManager.GenerateTokenPair(
		foundUser.ID.String(),
		foundUser.Email.String(),
		string(foundUser.Role),
	)
	if err != nil {
		c.logger.Error("Failed to generate new tokens", "error", err, "user_id", foundUser.ID)
		return nil, err
	}

	// Store new refresh token
	newTokenHash := hashToken(newRefreshToken)
	newTokenEntity := session.NewRefreshToken(
		foundUser.ID,
		newTokenHash,
		storedToken.DeviceName,
		storedToken.Platform,
		storedToken.IPAddress,
		storedToken.UserAgent,
		time.Now().UTC().Add(time.Duration(c.jwtManager.Config().RefreshTokenExpiry)*time.Hour),
	)

	if err := c.refreshTokenRepo.Create(ctx, newTokenEntity); err != nil {
		c.logger.Error("Failed to store new refresh token", "error", err)
		return nil, err
	}

	// Audit log
	c.auditRepo.Create(ctx, audit.NewLog(
		foundUser.ID,
		audit.ActionAccess,
		"refresh_token",
		newTokenEntity.ID,
		storedToken.IPAddress,
		storedToken.UserAgent,
	))

	c.logger.Info("Token refreshed successfully", "user_id", foundUser.ID)

	return &dto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
