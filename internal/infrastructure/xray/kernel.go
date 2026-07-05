package xray

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/infrastructure/vpn"
	"github.com/suproxy/backend/internal/infrastructure/xray/binary"
	"github.com/suproxy/backend/internal/infrastructure/xray/config"
	"github.com/suproxy/backend/internal/infrastructure/xray/runtime"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

// Kernel implements the VPN Kernel interface for Xray-core
type Kernel struct {
	configGenerator config.Generator
	configValidator config.Validator
	configWriter    config.Writer
	runtimeManager  runtime.Manager
	binaryManager   binary.Manager
	logger          *logger.Logger
}

// NewKernel creates a new Xray kernel instance
func NewKernel(
	configGenerator config.Generator,
	configValidator config.Validator,
	configWriter config.Writer,
	runtimeManager runtime.Manager,
	binaryManager binary.Manager,
	logger *logger.Logger,
) vpn.Kernel {
	return &Kernel{
		configGenerator: configGenerator,
		configValidator: configValidator,
		configWriter:    configWriter,
		runtimeManager:  runtimeManager,
		binaryManager:   binaryManager,
		logger:          logger,
	}
}

func (k *Kernel) Name() string {
	return "xray"
}

func (k *Kernel) Version(ctx context.Context) (string, error) {
	version, err := k.binaryManager.CurrentVersion(ctx)
	if err != nil {
		k.logger.Error("Failed to get Xray version", "error", err)
		return "", err
	}
	return version, nil
}

func (k *Kernel) GenerateConfig(ctx context.Context, instanceID uuid.UUID) ([]byte, error) {
	// Generate config from domain entities
	configJSON, err := k.configGenerator.GenerateJSON(ctx, instanceID)
	if err != nil {
		k.logger.Error("Failed to generate Xray config", "error", err, "instance_id", instanceID)
		return nil, err
	}

	// Validate generated config
	if err := k.configValidator.ValidateJSON(configJSON); err != nil {
		k.logger.Error("Generated config validation failed", "error", err, "instance_id", instanceID)
		return nil, err
	}

	k.logger.Info("Xray config generated successfully", "instance_id", instanceID)
	return configJSON, nil
}

func (k *Kernel) ValidateConfig(ctx context.Context, configData []byte) error {
	return k.configValidator.ValidateJSON(configData)
}

func (k *Kernel) Start(ctx context.Context, instanceID uuid.UUID) error {
	// Check if already running
	running, err := k.runtimeManager.IsRunning(ctx, instanceID)
	if err != nil {
		return err
	}
	if running {
		return fmt.Errorf("instance %s is already running", instanceID)
	}

	// Generate and write config
	configJSON, err := k.GenerateConfig(ctx, instanceID)
	if err != nil {
		return err
	}

	if err := k.configWriter.Write(ctx, instanceID, configJSON); err != nil {
		k.logger.Error("Failed to write Xray config", "error", err, "instance_id", instanceID)
		return err
	}

	// Start the process
	if err := k.runtimeManager.Start(ctx, instanceID); err != nil {
		k.logger.Error("Failed to start Xray process", "error", err, "instance_id", instanceID)
		return err
	}

	k.logger.Info("Xray instance started successfully", "instance_id", instanceID)
	return nil
}

func (k *Kernel) Stop(ctx context.Context, instanceID uuid.UUID) error {
	// Check if running
	running, err := k.runtimeManager.IsRunning(ctx, instanceID)
	if err != nil {
		return err
	}
	if !running {
		return fmt.Errorf("instance %s is not running", instanceID)
	}

	// Stop the process
	if err := k.runtimeManager.Stop(ctx, instanceID); err != nil {
		k.logger.Error("Failed to stop Xray process", "error", err, "instance_id", instanceID)
		return err
	}

	k.logger.Info("Xray instance stopped successfully", "instance_id", instanceID)
	return nil
}

func (k *Kernel) Restart(ctx context.Context, instanceID uuid.UUID) error {
	// Backup current config
	if err := k.configWriter.Backup(ctx, instanceID); err != nil {
		k.logger.Warn("Failed to backup config before restart", "error", err, "instance_id", instanceID)
	}

	// Generate new config
	configJSON, err := k.GenerateConfig(ctx, instanceID)
	if err != nil {
		return err
	}

	// Write new config
	if err := k.configWriter.Write(ctx, instanceID, configJSON); err != nil {
		k.logger.Error("Failed to write Xray config", "error", err, "instance_id", instanceID)
		return err
	}

	// Restart the process
	if err := k.runtimeManager.Restart(ctx, instanceID); err != nil {
		k.logger.Error("Failed to restart Xray process", "error", err, "instance_id", instanceID)
		return err
	}

	k.logger.Info("Xray instance restarted successfully", "instance_id", instanceID)
	return nil
}

func (k *Kernel) Reload(ctx context.Context, instanceID uuid.UUID) error {
	// Generate and write new config
	configJSON, err := k.GenerateConfig(ctx, instanceID)
	if err != nil {
		return err
	}

	if err := k.configWriter.Write(ctx, instanceID, configJSON); err != nil {
		k.logger.Error("Failed to write Xray config", "error", err, "instance_id", instanceID)
		return err
	}

	// Reload without restarting (hot reload)
	if err := k.runtimeManager.Reload(ctx, instanceID); err != nil {
		k.logger.Error("Failed to reload Xray config", "error", err, "instance_id", instanceID)
		return err
	}

	k.logger.Info("Xray config reloaded successfully", "instance_id", instanceID)
	return nil
}

func (k *Kernel) Status(ctx context.Context, instanceID uuid.UUID) (vpn.KernelStatus, error) {
	processStatus, err := k.runtimeManager.Status(ctx, instanceID)
	if err != nil {
		return vpn.KernelStatus{}, err
	}

	return vpn.KernelStatus{
		Running:     processStatus.Running,
		ProcessID:   processStatus.ProcessID,
		Uptime:      int64(processStatus.Uptime.Seconds()),
		ConfigPath:  processStatus.ConfigPath,
		LogPath:     processStatus.LogPath,
		ErrorReason: processStatus.ErrorReason,
	}, nil
}

func (k *Kernel) IsRunning(ctx context.Context, instanceID uuid.UUID) (bool, error) {
	return k.runtimeManager.IsRunning(ctx, instanceID)
}

func (k *Kernel) GetProcessID(ctx context.Context, instanceID uuid.UUID) (int, error) {
	return k.runtimeManager.GetProcessID(ctx, instanceID)
}
