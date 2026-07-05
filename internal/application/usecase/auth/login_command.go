package auth

import (
	"context"

	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	"github.com/suproxy/backend/internal/domain/user"
	"github.com/suproxy/backend/internal/infrastructure/jwt"
	"github.com/suproxy/backend/internal/infrastructure/logger"
	"github.com/suproxy/backend/internal/infrastructure/security"
)

type LoginCommand struct {
	userRepo   user.Repository
	jwtManager *jwt.Manager
	logger     *logger.Logger
}

func NewLoginCommand(userRepo user.Repository, jwtManager *jwt.Manager, logger *logger.Logger) *LoginCommand {
	return &LoginCommand{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		logger:     logger,
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
		return nil, user.ErrInvalidCredentials
	}

	// Check if user is active
	if !foundUser.IsActive() {
		c.logger.Warn("Inactive user attempted login", "user_id", foundUser.ID)
		return nil, user.ErrUserNotActive
	}

	// Verify password
	if err := security.CheckPassword(foundUser.Password.Hash(), req.Password); err != nil {
		c.logger.Warn("Invalid password attempt", "user_id", foundUser.ID)
		return nil, user.ErrInvalidCredentials
	}

	// Generate JWT tokens
	accessToken, refreshToken, err := c.jwtManager.GenerateTokenPair(
		foundUser.ID.String(),
		foundUser.Email.String(),
	)
	if err != nil {
		c.logger.Error("Failed to generate tokens", "error", err, "user_id", foundUser.ID)
		return nil, err
	}

	c.logger.Info("User logged in successfully", "user_id", foundUser.ID, "email", email.String())

	// Map to response
	userResponse := mapper.ToUserResponse(foundUser)

	return &dto.AuthResponse{
		User:         userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
