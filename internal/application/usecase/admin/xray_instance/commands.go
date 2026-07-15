package xray_instance

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/service"
	"github.com/suproxy/backend/internal/domain/audit"
	"github.com/suproxy/backend/internal/domain/xray"
	"github.com/suproxy/backend/internal/infrastructure/xray/runtime"
)

// StartInstanceCommand handles starting an Xray instance
type StartInstanceCommand struct {
	instanceRepo   xray.XrayInstanceRepository
	processManager runtime.Manager
	auditRepo      audit.Repository
}

func NewStartInstanceCommand(
	instanceRepo xray.XrayInstanceRepository,
	processManager runtime.Manager,
	auditRepo audit.Repository,
) *StartInstanceCommand {
	return &StartInstanceCommand{
		instanceRepo:   instanceRepo,
		processManager: processManager,
		auditRepo:      auditRepo,
	}
}

func (c *StartInstanceCommand) Execute(ctx context.Context, instanceID, adminID uuid.UUID, ip, userAgent string) error {
	// Find instance
	instance, err := c.instanceRepo.FindByID(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to find instance: %w", err)
	}

	// Start domain entity
	if err := instance.Start(); err != nil {
		return fmt.Errorf("failed to start instance entity: %w", err)
	}

	// Start process
	if err := c.processManager.Start(ctx, instanceID); err != nil {
		return fmt.Errorf("failed to start process: %w", err)
	}

	// Save instance
	if err := c.instanceRepo.Update(ctx, instance); err != nil {
		return fmt.Errorf("failed to update instance: %w", err)
	}

	// Create audit log
	auditLog := audit.NewLog(adminID, audit.ActionUpdate, "xray_instance", instanceID, ip, userAgent)
	auditLog.AddMetadata("event", "xray_instance_started")
	auditLog.AddMetadata("instance_id", instanceID.String())

	if err := c.auditRepo.Create(ctx, auditLog); err != nil {
		return fmt.Errorf("instance started but audit log failed: %w", err)
	}

	return nil
}

// StopInstanceCommand handles stopping an Xray instance
type StopInstanceCommand struct {
	instanceRepo   xray.XrayInstanceRepository
	processManager runtime.Manager
	auditRepo      audit.Repository
}

func NewStopInstanceCommand(
	instanceRepo xray.XrayInstanceRepository,
	processManager runtime.Manager,
	auditRepo audit.Repository,
) *StopInstanceCommand {
	return &StopInstanceCommand{
		instanceRepo:   instanceRepo,
		processManager: processManager,
		auditRepo:      auditRepo,
	}
}

func (c *StopInstanceCommand) Execute(ctx context.Context, instanceID, adminID uuid.UUID, ip, userAgent string) error {
	// Find instance
	instance, err := c.instanceRepo.FindByID(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to find instance: %w", err)
	}

	// Stop domain entity
	if err := instance.Stop(); err != nil {
		return fmt.Errorf("failed to stop instance entity: %w", err)
	}

	// Stop process
	if err := c.processManager.Stop(ctx, instanceID); err != nil {
		return fmt.Errorf("failed to stop process: %w", err)
	}

	// Save instance
	if err := c.instanceRepo.Update(ctx, instance); err != nil {
		return fmt.Errorf("failed to update instance: %w", err)
	}

	// Create audit log
	auditLog := audit.NewLog(adminID, audit.ActionUpdate, "xray_instance", instanceID, ip, userAgent)
	auditLog.AddMetadata("event", "xray_instance_stopped")
	auditLog.AddMetadata("instance_id", instanceID.String())

	if err := c.auditRepo.Create(ctx, auditLog); err != nil {
		return fmt.Errorf("instance stopped but audit log failed: %w", err)
	}

	return nil
}

// RestartInstanceCommand handles restarting an Xray instance
type RestartInstanceCommand struct {
	instanceRepo   xray.XrayInstanceRepository
	processManager runtime.Manager
	auditRepo      audit.Repository
}

func NewRestartInstanceCommand(
	instanceRepo xray.XrayInstanceRepository,
	processManager runtime.Manager,
	auditRepo audit.Repository,
) *RestartInstanceCommand {
	return &RestartInstanceCommand{
		instanceRepo:   instanceRepo,
		processManager: processManager,
		auditRepo:      auditRepo,
	}
}

func (c *RestartInstanceCommand) Execute(ctx context.Context, instanceID, adminID uuid.UUID, ip, userAgent string) error {
	// Find instance
	instance, err := c.instanceRepo.FindByID(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to find instance: %w", err)
	}

	// Restart domain entity
	if err := instance.Restart(); err != nil {
		return fmt.Errorf("failed to restart instance entity: %w", err)
	}

	// Restart process
	if err := c.processManager.Restart(ctx, instanceID); err != nil {
		return fmt.Errorf("failed to restart process: %w", err)
	}

	// Save instance
	if err := c.instanceRepo.Update(ctx, instance); err != nil {
		return fmt.Errorf("failed to update instance: %w", err)
	}

	// Create audit log
	auditLog := audit.NewLog(adminID, audit.ActionUpdate, "xray_instance", instanceID, ip, userAgent)
	auditLog.AddMetadata("event", "xray_instance_restarted")
	auditLog.AddMetadata("instance_id", instanceID.String())

	if err := c.auditRepo.Create(ctx, auditLog); err != nil {
		return fmt.Errorf("instance restarted but audit log failed: %w", err)
	}

	return nil
}

// ReloadInstanceCommand handles reloading Xray instance config
// REUSES XrayProvisioningService from Aşama 16
type ReloadInstanceCommand struct {
	provisioningService *service.XrayProvisioningService
	auditRepo           audit.Repository
}

func NewReloadInstanceCommand(
	provisioningService *service.XrayProvisioningService,
	auditRepo audit.Repository,
) *ReloadInstanceCommand {
	return &ReloadInstanceCommand{
		provisioningService: provisioningService,
		auditRepo:           auditRepo,
	}
}

func (c *ReloadInstanceCommand) Execute(ctx context.Context, instanceID, adminID uuid.UUID, ip, userAgent string) error {
	startTime := time.Now()

	// REUSE existing RegenerateAndReload from XrayProvisioningService
	if err := c.provisioningService.RegenerateAndReload(ctx, instanceID, adminID, ip, userAgent); err != nil {
		return fmt.Errorf("failed to reload instance: %w", err)
	}

	// Create audit log (additional to the ones created by provisioning service)
	auditLog := audit.NewLog(adminID, audit.ActionUpdate, "xray_instance", instanceID, ip, userAgent)
	auditLog.AddMetadata("event", "xray_instance_reload_admin_initiated")
	auditLog.AddMetadata("instance_id", instanceID.String())
	auditLog.AddMetadata("duration_ms", time.Since(startTime).Milliseconds())

	if err := c.auditRepo.Create(ctx, auditLog); err != nil {
		return fmt.Errorf("instance reloaded but audit log failed: %w", err)
	}

	return nil
}

// CheckInstanceHealthCommand handles health check for Xray instance
type CheckInstanceHealthCommand struct {
	processManager runtime.Manager
}

func NewCheckInstanceHealthCommand(processManager runtime.Manager) *CheckInstanceHealthCommand {
	return &CheckInstanceHealthCommand{
		processManager: processManager,
	}
}

// HealthCheckResult represents health check result
type HealthCheckResult struct {
	Healthy bool
	Status  string
	Message string
}

func (c *CheckInstanceHealthCommand) Execute(ctx context.Context, instanceID uuid.UUID) (*HealthCheckResult, error) {
	// Check if process is running
	isRunning, err := c.processManager.IsRunning(ctx, instanceID)
	if err != nil {
		return &HealthCheckResult{
			Healthy: false,
			Status:  "error",
			Message: fmt.Sprintf("failed to check instance status: %v", err),
		}, nil
	}

	if !isRunning {
		return &HealthCheckResult{
			Healthy: false,
			Status:  "stopped",
			Message: "instance is not running",
		}, nil
	}

	// Get process status
	status, err := c.processManager.Status(ctx, instanceID)
	if err != nil {
		return &HealthCheckResult{
			Healthy: true,
			Status:  "running",
			Message: "instance is running but status check failed",
		}, nil
	}

	// Get status string representation
	statusStr := "running"
	if status != nil {
		statusStr = "running" // ProcessStatus is a struct, use a default string
	}

	return &HealthCheckResult{
		Healthy: true,
		Status:  "healthy",
		Message: fmt.Sprintf("instance is healthy (status: %s)", statusStr),
	}, nil
}
