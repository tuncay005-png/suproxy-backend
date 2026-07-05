package auth

import (
	"context"

	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/application/mapper"
	"github.com/suproxy/backend/internal/domain/user"
	"github.com/suproxy/backend/internal/infrastructure/logger"
	"github.com/suproxy/backend/internal/infrastructure/security"
)

type RegisterCommand struct {
	userRepo user.Repository
	logger   *logger.Logger
}

func NewRegisterCommand(userRepo user.Repository, logger *logger.Logger) *RegisterCommand {
	return &RegisterCommand{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (c *RegisterCommand) Execute(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Validate password strength
	if err := security.ValidatePasswordStrength(req.Password); err != nil {
		return nil, err
	}

	// Create email value object
	email, err := user.NewEmail(req.Email)
	if err != nil {
		return nil, user.ErrInvalidEmail
	}

	// Check if user already exists
	exists, err := c.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		c.logger.Error("Failed to check user existence", "error", err)
		return nil, err
	}
	if exists {
		return nil, user.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		c.logger.Error("Failed to hash password", "error", err)
		return nil, err
	}

	// Create password value object
	password, err := user.NewPassword(hashedPassword)
	if err != nil {
		return nil, err
	}

	// Create profile
	profile := user.NewProfile(req.FirstName, req.LastName, "", "")

	// Create user entity
	newUser, err := user.NewUser(email, password, profile)
	if err != nil {
		c.logger.Error("Failed to create user entity", "error", err)
		return nil, err
	}

	// Save user
	if err := c.userRepo.Create(ctx, newUser); err != nil {
		c.logger.Error("Failed to save user", "error", err)
		return nil, err
	}

	c.logger.Info("User registered successfully", "user_id", newUser.ID, "email", email.String())

	// Map to response
	userResponse := mapper.ToUserResponse(newUser)

	return &dto.AuthResponse{
		User: userResponse,
	}, nil
}
