package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/audit"
	"github.com/suproxy/backend/internal/domain/session"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

type LogoutCommand struct {
	refreshTokenRepo session.RefreshTokenRepository
	auditRepo        audit.Repository
	logger           *logger.Logger
}

func NewLogoutCommand(
	refreshTokenRepo session.RefreshTokenRepository,
	auditRepo audit.Repository,
	logger *logger.Logger,
) *LogoutCommand {
	return &LogoutCommand{
		refreshTokenRepo: refreshTokenRepo,
		auditRepo:        auditRepo,
		logger:           logger,
	}
}

func (c *LogoutCommand) ExecuteSingle(ctx context.Context, userID uuid.UUID, tokenID uuid.UUID, ipAddress, userAgent string) error {
	// Revoke specific token
	if err := c.refreshTokenRepo.RevokeByID(ctx, tokenID); err != nil {
		c.logger.Error("Failed to revoke token", "error", err, "token_id", tokenID)
		return err
	}

	// Audit log
	auditLog := audit.NewLog(userID, audit.ActionLogout, "refresh_token", tokenID, ipAddress, userAgent)
	auditLog.AddMetadata("type", "single")
	c.auditRepo.Create(ctx, auditLog)

	c.logger.Info("User logged out from single device", "user_id", userID)
	return nil
}

func (c *LogoutCommand) ExecuteAll(ctx context.Context, userID uuid.UUID, ipAddress, userAgent string) error {
	// Revoke all tokens
	if err := c.refreshTokenRepo.RevokeAllByUserID(ctx, userID); err != nil {
		c.logger.Error("Failed to revoke all tokens", "error", err, "user_id", userID)
		return err
	}

	// Audit log
	auditLog := audit.NewLog(userID, audit.ActionLogout, "user", userID, ipAddress, userAgent)
	auditLog.AddMetadata("type", "all")
	c.auditRepo.Create(ctx, auditLog)

	c.logger.Info("User logged out from all devices", "user_id", userID)
	return nil
}
