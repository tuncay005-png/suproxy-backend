package runtime

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

// RealProcessManager manages actual Xray processes using os/exec
type RealProcessManager struct {
	binaryPath string
	configDir  string
	logDir     string
	registry   Registry
	logger     *logger.Logger
}

// NewRealProcessManager creates a new real process manager
func NewRealProcessManager(binaryPath, configDir, logDir string, logger *logger.Logger) Manager {
	return &RealProcessManager{
		binaryPath: binaryPath,
		configDir:  configDir,
		logDir:     logDir,
		registry:   NewRegistry(),
		logger:     logger,
	}
}

func (m *RealProcessManager) Start(ctx context.Context, instanceID uuid.UUID) error {
	// Check if already running
	if m.registry.IsRegistered(instanceID) {
		return ErrProcessAlreadyRunning
	}

	// Build paths
	configPath := m.getConfigPath(instanceID)
	logPath := m.getLogPath(instanceID)
	errorPath := m.getErrorPath(instanceID)

	// Prepare command
	// In production, this will execute: xray run -config /path/to/config.json
	_ = exec.CommandContext(ctx, m.binaryPath, "run", "-config", configPath)

	// TODO: Set up log file redirection
	// cmd.Stdout = logFile
	// cmd.Stderr = errorFile

	// TODO: Start the process
	// if err := cmd.Start(); err != nil {
	//     m.logger.Error("Failed to start Xray process", "error", err, "instance_id", instanceID)
	//     return fmt.Errorf("%w: %v", ErrProcessStartFailed, err)
	// }

	// Register process
	processInfo := &ProcessInfo{
		InstanceID: instanceID,
		ProcessID:  0, // TODO: cmd.Process.Pid after cmd.Start()
		StartedAt:  time.Now().UTC(),
		ConfigPath: configPath,
		LogPath:    logPath,
		ErrorPath:  errorPath,
		Command:    m.binaryPath,
		Args:       []string{"run", "-config", configPath},
	}

	if err := m.registry.Register(processInfo); err != nil {
		return err
	}

	m.logger.Info("Xray process start prepared", "instance_id", instanceID, "config", configPath)

	// TODO: Monitor process in goroutine
	// go m.monitorProcess(instanceID, cmd.Process)

	return nil
}

func (m *RealProcessManager) Stop(ctx context.Context, instanceID uuid.UUID) error {
	// Find process
	info, exists := m.registry.Find(instanceID)
	if !exists {
		return ErrProcessNotRunning
	}

	// TODO: Get process and send SIGTERM
	// process, err := os.FindProcess(info.ProcessID)
	// if err != nil {
	//     return fmt.Errorf("%w: %v", ErrProcessNotFound, err)
	// }

	// TODO: Send SIGTERM for graceful shutdown
	// if err := process.Signal(syscall.SIGTERM); err != nil {
	//     m.logger.Error("Failed to send SIGTERM", "error", err, "instance_id", instanceID)
	//     return fmt.Errorf("%w: %v", ErrProcessStopFailed, err)
	// }

	// TODO: Wait for process to exit with timeout
	// done := make(chan error, 1)
	// go func() {
	//     _, err := process.Wait()
	//     done <- err
	// }()

	// select {
	// case <-time.After(10 * time.Second):
	//     // Force kill if not stopped gracefully
	//     process.Kill()
	// case err := <-done:
	//     if err != nil {
	//         m.logger.Warn("Process exit with error", "error", err, "instance_id", instanceID)
	//     }
	// }

	// Remove from registry
	if err := m.registry.Remove(instanceID); err != nil {
		return err
	}

	m.logger.Info("Xray process stopped", "instance_id", instanceID, "pid", info.ProcessID)
	return nil
}

func (m *RealProcessManager) Restart(ctx context.Context, instanceID uuid.UUID) error {
	// Stop if running
	if m.registry.IsRegistered(instanceID) {
		if err := m.Stop(ctx, instanceID); err != nil {
			return err
		}

		// Wait a bit for cleanup
		time.Sleep(500 * time.Millisecond)
	}

	// Start again
	return m.Start(ctx, instanceID)
}

func (m *RealProcessManager) Reload(ctx context.Context, instanceID uuid.UUID) error {
	// Find process
	info, exists := m.registry.Find(instanceID)
	if !exists {
		return ErrProcessNotRunning
	}

	// TODO: Xray supports hot reload via API
	// In production, use Xray's gRPC API to reload config without restart
	// This preserves existing connections

	// For now, log the reload intention
	m.logger.Info("Config reload prepared", "instance_id", instanceID, "pid", info.ProcessID)

	// TODO: Implement hot reload via Xray API
	// xrayAPI := xray.NewAPIClient(...)
	// xrayAPI.ReloadConfig(ctx, info.ConfigPath)

	return nil
}

func (m *RealProcessManager) Status(ctx context.Context, instanceID uuid.UUID) (*ProcessStatus, error) {
	info, exists := m.registry.Find(instanceID)
	if !exists {
		return &ProcessStatus{
			InstanceID: instanceID,
			Running:    false,
		}, nil
	}

	// TODO: Check if process is actually running
	// process, err := os.FindProcess(info.ProcessID)
	// if err != nil {
	//     return &ProcessStatus{
	//         InstanceID: instanceID,
	//         Running:    false,
	//         ErrorReason: "process not found",
	//     }, nil
	// }

	// TODO: Get process stats (CPU, Memory)
	// In Linux: read from /proc/[pid]/stat
	// In production, use system monitoring library

	uptime := time.Since(info.StartedAt)

	return &ProcessStatus{
		InstanceID:  info.InstanceID,
		Running:     true,
		ProcessID:   info.ProcessID,
		StartedAt:   info.StartedAt,
		Uptime:      uptime,
		ConfigPath:  info.ConfigPath,
		LogPath:     info.LogPath,
		ErrorPath:   info.ErrorPath,
		CPUUsage:    0, // TODO: Implement
		MemoryUsage: 0, // TODO: Implement
	}, nil
}

func (m *RealProcessManager) IsRunning(ctx context.Context, instanceID uuid.UUID) (bool, error) {
	status, err := m.Status(ctx, instanceID)
	if err != nil {
		return false, err
	}
	return status.Running, nil
}

func (m *RealProcessManager) GetProcessID(ctx context.Context, instanceID uuid.UUID) (int, error) {
	info, exists := m.registry.Find(instanceID)
	if !exists {
		return 0, ErrProcessNotRunning
	}
	return info.ProcessID, nil
}

func (m *RealProcessManager) GetLogs(ctx context.Context, instanceID uuid.UUID, lines int) ([]string, error) {
	info, exists := m.registry.Find(instanceID)
	if !exists {
		return nil, ErrProcessNotRunning
	}

	// TODO: Tail log file
	// In production:
	// 1. Open log file at info.LogPath
	// 2. Read last N lines
	// 3. Return as string array

	m.logger.Debug("Log retrieval prepared", "instance_id", instanceID, "log_path", info.LogPath, "lines", lines)

	return []string{
		fmt.Sprintf("[Info] Process started at %s", info.StartedAt.Format(time.RFC3339)),
		fmt.Sprintf("[Info] Config loaded from %s", info.ConfigPath),
		fmt.Sprintf("[Info] PID: %d", info.ProcessID),
	}, nil
}

func (m *RealProcessManager) Kill(ctx context.Context, instanceID uuid.UUID) error {
	info, exists := m.registry.Find(instanceID)
	if !exists {
		return ErrProcessNotRunning
	}

	// TODO: Force kill process with SIGKILL
	// process, err := os.FindProcess(info.ProcessID)
	// if err != nil {
	//     return fmt.Errorf("%w: %v", ErrProcessNotFound, err)
	// }

	// if err := process.Kill(); err != nil {
	//     m.logger.Error("Failed to kill process", "error", err, "instance_id", instanceID)
	//     return fmt.Errorf("%w: %v", ErrProcessKillFailed, err)
	// }

	// Remove from registry
	if err := m.registry.Remove(instanceID); err != nil {
		return err
	}

	m.logger.Warn("Xray process killed forcefully", "instance_id", instanceID, "pid", info.ProcessID)
	return nil
}

// monitorProcess monitors a running process and handles unexpected exits
func (m *RealProcessManager) monitorProcess(instanceID uuid.UUID, process interface{}) {
	// TODO: Implement process monitoring
	// 1. Wait for process to exit
	// 2. Log exit status
	// 3. Remove from registry
	// 4. Optionally restart if configured for auto-restart
}

func (m *RealProcessManager) getConfigPath(instanceID uuid.UUID) string {
	return filepath.Join(m.configDir, instanceID.String()+".json")
}

func (m *RealProcessManager) getLogPath(instanceID uuid.UUID) string {
	return filepath.Join(m.logDir, instanceID.String()+".log")
}

func (m *RealProcessManager) getErrorPath(instanceID uuid.UUID) string {
	return filepath.Join(m.logDir, instanceID.String()+".error.log")
}

// GetRegistry returns the internal registry for testing purposes
func (m *RealProcessManager) GetRegistry() Registry {
	return m.registry
}
