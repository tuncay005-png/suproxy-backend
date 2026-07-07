package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	"github.com/suproxy/backend/internal/domain/audit"
	"github.com/suproxy/backend/internal/domain/session"
	"github.com/suproxy/backend/internal/domain/user"
	"github.com/suproxy/backend/internal/infrastructure/jwt"
	"github.com/suproxy/backend/internal/infrastructure/logger"
	"github.com/suproxy/backend/internal/infrastructure/metrics"
	"github.com/suproxy/backend/internal/infrastructure/security"
)

type LoginCommand struct {
	userRepo         user.Repository
	refreshTokenRepo session.RefreshTokenRepository
	auditRepo        audit.Repository
	jwtManager       *jwt.Manager
	logger           *logger.Logger
}

func NewLoginCommand(
	userRepo user.Repository,
	refreshTokenRepo session.RefreshTokenRepository,
	auditRepo audit.Repository,
	jwtManager *jwt.Manager,
	logger *logger.Logger,
) *LoginCommand {
	return &LoginCommand{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		auditRepo:        auditRepo,
		jwtManager:       jwtManager,
		logger:           logger,
	}
}

func (c *LoginCommand) Execute(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Create email value object
	email, err := user.NewEmail(req.Email)
	if err != nil {
		return nil, user.ErrInvalidEmail
	}

	// Find user by email
	foundUser, err := c.userRepo.FindByEmail(ctx, email)
	if err != nil {
		c.logger.Warn("User not found", "email", email.String())

		// Audit log for failed login
		c.auditRepo.Create(ctx, audit.NewLog(
			uuid.Nil,
			audit.ActionLogin,
			"user",
			uuid.Nil,
			req.IPAddress,
			req.UserAgent,
		))

		// Record failed login metric
		metrics.IncUserLoginFailures()

		return nil, user.ErrInvalidCredentials
	}

	// Check if account is locked
	if foundUser.IsLocked() {
		c.logger.Warn("Locked account attempted login", "user_id", foundUser.ID)

		// Audit log
		c.auditRepo.Create(ctx, audit.NewLog(
			foundUser.ID,
			audit.ActionLogin,
			"user",
			foundUser.ID,
			req.IPAddress,
			req.UserAgent,
		))

		return nil, user.ErrUserLocked
	}

	// Check if user is active
	if !foundUser.IsActive() {
		c.logger.Warn("Inactive user attempted login", "user_id", foundUser.ID)
		return nil, user.ErrUserNotActive
	}

	// Verify password
	if err := security.CheckPassword(foundUser.Password.Hash(), req.Password); err != nil {
		c.logger.Warn("Invalid password attempt", "user_id", foundUser.ID)

		// Record failed login
		foundUser.RecordFailedLogin()
		c.userRepo.Update(ctx, foundUser)

		// Audit log
		c.auditRepo.Create(ctx, audit.NewLog(
			foundUser.ID,
			audit.ActionLogin,
			"user",
			foundUser.ID,
			req.IPAddress,
			req.UserAgent,
		))

		// Record failed login metric
		metrics.IncUserLoginFailures()

		return nil, user.ErrInvalidCredentials
	}

	// Record successful login
	foundUser.RecordSuccessfulLogin(req.IPAddress)
	if err := c.userRepo.Update(ctx, foundUser); err != nil {
		c.logger.Error("Failed to update user login info", "error", err)
	}

	// Generate JWT tokens
	accessToken, refreshToken, err := c.jwtManager.GenerateTokenPair(
		foundUser.ID.String(),
		foundUser.Email.String(),
		string(foundUser.Role),
	)
	if err != nil {
		c.logger.Error("Failed to generate tokens", "error", err, "user_id", foundUser.ID)
		return nil, err
	}

	// Store refresh token
	tokenHash := hashRefreshToken(refreshToken)
	refreshTokenEntity := session.NewRefreshToken(
		foundUser.ID,
		tokenHash,
		req.DeviceName,
		req.Platform,
		req.IPAddress,
		req.UserAgent,
		time.Now().UTC().Add(time.Duration(c.jwtManager.Config().RefreshTokenExpiry)*time.Hour),
	)

	if err := c.refreshTokenRepo.Create(ctx, refreshTokenEntity); err != nil {
		c.logger.Error("Failed to store refresh token", "error", err)
		// Continue anyway, user can still use access token
	}

	// Audit log for successful login
	auditLog := audit.NewLog(
		foundUser.ID,
		audit.ActionLogin,
		"user",
		foundUser.ID,
		req.IPAddress,
		req.UserAgent,
	)
	auditLog.AddMetadata("status", "success")
	c.auditRepo.Create(ctx, auditLog)
	
	// Record successful login metric
	metrics.IncUserLogins()

	c.logger.Info("User logged in successfully", "user_id", foundUser.ID, "email", email.String())

	// Map to response
	userResponse := mapper.ToUserResponse(foundUser)

	return &dto.AuthResponse{
		User:         userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func hashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
