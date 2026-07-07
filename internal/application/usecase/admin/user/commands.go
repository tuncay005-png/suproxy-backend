package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/audit"
	"github.com/suproxy/backend/internal/domain/user"
	"github.com/suproxy/backend/internal/interfaces/http/middleware"
)

// UpdateUserStatusCommand handles updating user status
type UpdateUserStatusCommand struct {
	userRepo  user.Repository
	auditRepo audit.Repository
}

func NewUpdateUserStatusCommand(userRepo user.Repository, auditRepo audit.Repository) *UpdateUserStatusCommand {
	return &UpdateUserStatusCommand{
		userRepo:  userRepo,
		auditRepo: auditRepo,
	}
}

func (c *UpdateUserStatusCommand) Execute(ctx context.Context, userID uuid.UUID, newStatus string, adminID uuid.UUID, ip, userAgent string) error {
	// Find user
	u, err := c.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	oldStatus := u.Status

	// Change status based on request
	switch newStatus {
	case "active":
		if err := u.Activate(); err != nil {
			return fmt.Errorf("failed to activate user: %w", err)
		}
	case "inactive":
		if err := u.Deactivate(); err != nil {
			return fmt.Errorf("failed to deactivate user: %w", err)
		}
	case "suspended":
		if err := u.Suspend(); err != nil {
			return fmt.Errorf("failed to suspend user: %w", err)
		}
	default:
		return fmt.Errorf("invalid status: %s", newStatus)
	}

	// Save user
	if err := c.userRepo.Update(ctx, u); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Create audit log
	auditLog := audit.NewLog(adminID, audit.ActionUpdate, "user", userID, ip, userAgent)
	auditLog.AddMetadata("event", "user_status_changed")
	auditLog.AddMetadata("old_status", string(oldStatus))
	auditLog.AddMetadata("new_status", newStatus)
	auditLog.AddMetadata("target_user_email", u.Email.String())

	if err := c.auditRepo.Create(ctx, auditLog); err != nil {
		// Log audit creation failure but don't fail the operation
		return fmt.Errorf("user status updated but audit log failed: %w", err)
	}

	return nil
}

// UpdateUserRoleCommand handles updating user role
type UpdateUserRoleCommand struct {
	userRepo  user.Repository
	auditRepo audit.Repository
}

func NewUpdateUserRoleCommand(userRepo user.Repository, auditRepo audit.Repository) *UpdateUserRoleCommand {
	return &UpdateUserRoleCommand{
		userRepo:  userRepo,
		auditRepo: auditRepo,
	}
}

func (c *UpdateUserRoleCommand) Execute(ctx context.Context, userID uuid.UUID, newRole string, adminID uuid.UUID, ip, userAgent string) error {
	// Prevent self-demotion
	if adminID == userID && newRole != "admin" {
		return middleware.ErrAdminSelfDemotion
	}

	// Find user
	u, err := c.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	oldRole := u.Role

	// Update role
	u.Role = user.Role(newRole)
	u.UpdatedAt = u.UpdatedAt // trigger UpdatedAt update

	// Save user
	if err := c.userRepo.Update(ctx, u); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Create audit log
	auditLog := audit.NewLog(adminID, audit.ActionUpdate, "user", userID, ip, userAgent)
	auditLog.AddMetadata("event", "user_role_changed")
	auditLog.AddMetadata("old_role", string(oldRole))
	auditLog.AddMetadata("new_role", newRole)
	auditLog.AddMetadata("target_user_email", u.Email.String())

	if err := c.auditRepo.Create(ctx, auditLog); err != nil {
		// Log audit creation failure but don't fail the operation
		return fmt.Errorf("user role updated but audit log failed: %w", err)
	}

	return nil
}
