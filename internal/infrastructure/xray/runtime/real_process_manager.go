package runtime

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
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
	commands   map[uuid.UUID]*exec.Cmd // Track running commands
}

// NewRealProcessManager creates a new real process manager
func NewRealProcessManager(binaryPath, configDir, logDir string, logger *logger.Logger) Manager {
	return &RealProcessManager{
		binaryPath: binaryPath,
		configDir:  configDir,
		logDir:     logDir,
		registry:   NewRegistry(),
		logger:     logger,
		commands:   make(map[uuid.UUID]*exec.Cmd),
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

	// Ensure log directory exists
	if err := os.MkdirAll(m.logDir, 0755); err != nil {
		m.logger.Error("Failed to create log directory", "error", err)
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open log files
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		m.logger.Error("Failed to open log file", "error", err, "path", logPath)
		return fmt.Errorf("failed to open log file: %w", err)
	}

	errorFile, err := os.OpenFile(errorPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logFile.Close()
		m.logger.Error("Failed to open error log file", "error", err, "path", errorPath)
		return fmt.Errorf("failed to open error log file: %w", err)
	}

	// Prepare command
	cmd := exec.CommandContext(ctx, m.binaryPath, "run", "-config", configPath)
	cmd.Stdout = logFile
	cmd.Stderr = errorFile

	// Set process group for proper signal handling
	m.setProcessGroup(cmd)

	// Start the process
	if err := cmd.Start(); err != nil {
		logFile.Close()
		errorFile.Close()
		m.logger.Error("Failed to start Xray process", "error", err, "instance_id", instanceID)
		return fmt.Errorf("%w: %v", ErrProcessStartFailed, err)
	}

	// Store command for later reference
	m.commands[instanceID] = cmd

	// Register process
	processInfo := &ProcessInfo{
		InstanceID: instanceID,
		ProcessID:  cmd.Process.Pid,
		StartedAt:  time.Now().UTC(),
		ConfigPath: configPath,
		LogPath:    logPath,
		ErrorPath:  errorPath,
		Command:    m.binaryPath,
		Args:       []string{"run", "-config", configPath},
	}

	if err := m.registry.Register(processInfo); err != nil {
		// Kill the process if registration fails
		if killErr := cmd.Process.Kill(); killErr != nil {
			m.logger.Error("Failed to kill process after registration failure", "error", killErr, "pid", cmd.Process.Pid)
		}
		logFile.Close()
		errorFile.Close()
		delete(m.commands, instanceID)
		return err
	}

	m.logger.Info("Xray process started successfully",
		"instance_id", instanceID,
		"pid", cmd.Process.Pid,
		"config", configPath)

	// Monitor process in background
	go m.monitorProcess(instanceID, cmd, logFile, errorFile)

	return nil
}

func (m *RealProcessManager) Stop(ctx context.Context, instanceID uuid.UUID) error {
	// Find process
	info, exists := m.registry.Find(instanceID)
	if !exists {
		return ErrProcessNotRunning
	}

	// Get process
	process, err := os.FindProcess(info.ProcessID)
	if err != nil {
		m.logger.Error("Failed to find process", "error", err, "instance_id", instanceID)
		return fmt.Errorf("%w: %v", ErrProcessNotFound, err)
	}

	// Send SIGTERM for graceful shutdown
	m.logger.Info("Sending SIGTERM to process", "instance_id", instanceID, "pid", info.ProcessID)
	if err := m.sendSignal(process, syscall.SIGTERM); err != nil {
		m.logger.Error("Failed to send SIGTERM", "error", err, "instance_id", instanceID)
		return fmt.Errorf("%w: %v", ErrProcessStopFailed, err)
	}

	// Wait for process to exit with timeout
	done := make(chan error, 1)
	go func() {
		_, err := process.Wait()
		done <- err
	}()

	select {
	case <-time.After(10 * time.Second):
		// Force kill if not stopped gracefully
		m.logger.Warn("Process did not stop gracefully, forcing kill", "instance_id", instanceID)
		if err := process.Kill(); err != nil {
			m.logger.Error("Failed to kill process", "error", err, "instance_id", instanceID)
		}
	case err := <-done:
		if err != nil {
			m.logger.Warn("Process exited with error", "error", err, "instance_id", instanceID)
		}
	}

	// Cleanup
	delete(m.commands, instanceID)

	// Remove from registry
	if err := m.registry.Remove(instanceID); err != nil {
		return err
	}

	m.logger.Info("Xray process stopped successfully", "instance_id", instanceID, "pid", info.ProcessID)
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

	// Get process
	process, err := os.FindProcess(info.ProcessID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrProcessNotFound, err)
	}

	// Send SIGHUP for config reload
	// Xray supports SIGHUP for hot reload without dropping connections
	m.logger.Info("Sending SIGHUP for config reload", "instance_id", instanceID, "pid", info.ProcessID)
	if err := m.sendSignal(process, syscall.SIGHUP); err != nil {
		m.logger.Error("Failed to send SIGHUP", "error", err, "instance_id", instanceID)
		return fmt.Errorf("failed to reload config: %w", err)
	}

	m.logger.Info("Config reload signal sent successfully", "instance_id", instanceID, "pid", info.ProcessID)
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

	// Check if process is actually running
	process, err := os.FindProcess(info.ProcessID)
	if err != nil {
		return &ProcessStatus{
			InstanceID:  instanceID,
			Running:     false,
			ErrorReason: "process not found",
		}, nil
	}

	// Verify process is still alive by sending signal 0
	if err := m.sendSignal(process, syscall.Signal(0)); err != nil {
		return &ProcessStatus{
			InstanceID:  instanceID,
			Running:     false,
			ErrorReason: "process not responding",
		}, nil
	}

	uptime := time.Since(info.StartedAt)

	// Get process stats (basic implementation)
	cpuUsage, memoryUsage := m.getProcessStats(info.ProcessID)

	return &ProcessStatus{
		InstanceID:  info.InstanceID,
		Running:     true,
		ProcessID:   info.ProcessID,
		StartedAt:   info.StartedAt,
		Uptime:      uptime,
		ConfigPath:  info.ConfigPath,
		LogPath:     info.LogPath,
		ErrorPath:   info.ErrorPath,
		CPUUsage:    cpuUsage,
		MemoryUsage: int64(memoryUsage),
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

	// Read last N lines from log file
	logLines, err := m.tailFile(info.LogPath, lines)
	if err != nil {
		m.logger.Error("Failed to read log file", "error", err, "path", info.LogPath)
		return nil, fmt.Errorf("failed to read log file: %w", err)
	}

	return logLines, nil
}

func (m *RealProcessManager) Kill(ctx context.Context, instanceID uuid.UUID) error {
	info, exists := m.registry.Find(instanceID)
	if !exists {
		return ErrProcessNotRunning
	}

	// Force kill process with SIGKILL
	process, err := os.FindProcess(info.ProcessID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrProcessNotFound, err)
	}

	m.logger.Warn("Force killing process", "instance_id", instanceID, "pid", info.ProcessID)
	if err := process.Kill(); err != nil {
		m.logger.Error("Failed to kill process", "error", err, "instance_id", instanceID)
		// Note: Process might already be dead, so we continue with cleanup
		// errcheck: acknowledged - Kill failure is logged but not critical for cleanup flow
		_ = err
	}

	// Cleanup
	delete(m.commands, instanceID)

	// Remove from registry
	if err := m.registry.Remove(instanceID); err != nil {
		return err
	}

	m.logger.Warn("Xray process killed forcefully", "instance_id", instanceID, "pid", info.ProcessID)
	return nil
}

// monitorProcess monitors a running process and handles unexpected exits
func (m *RealProcessManager) monitorProcess(instanceID uuid.UUID, cmd *exec.Cmd, logFile, errorFile *os.File) {
	defer logFile.Close()
	defer errorFile.Close()

	// Wait for process to exit
	err := cmd.Wait()

	// Log exit status
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			m.logger.Error("Process exited with error",
				"instance_id", instanceID,
				"pid", cmd.Process.Pid,
				"exit_code", exitErr.ExitCode(),
				"error", exitErr)
		} else {
			m.logger.Error("Process wait failed",
				"instance_id", instanceID,
				"pid", cmd.Process.Pid,
				"error", err)
		}
	} else {
		m.logger.Info("Process exited normally",
			"instance_id", instanceID,
			"pid", cmd.Process.Pid)
	}

	// Remove from registry
	if err := m.registry.Remove(instanceID); err != nil {
		m.logger.Error("Failed to remove process from registry",
			"instance_id", instanceID,
			"error", err)
	}

	// Cleanup command reference
	delete(m.commands, instanceID)

	// TODO: Implement auto-restart logic if configured
	// if m.autoRestart {
	//     m.logger.Info("Auto-restarting crashed process", "instance_id", instanceID)
	//     go m.Start(context.Background(), instanceID)
	// }
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

// sendSignal sends a signal to a process (cross-platform compatible)
func (m *RealProcessManager) sendSignal(process *os.Process, sig syscall.Signal) error {
	if runtime.GOOS == "windows" {
		// Windows doesn't support POSIX signals
		// For graceful shutdown, we need to send a different signal or use process.Kill()
		if sig == syscall.SIGTERM || sig == syscall.SIGHUP {
			return process.Signal(os.Interrupt)
		}
		return process.Kill()
	}
	return process.Signal(sig)
}

// setProcessGroup sets process group for proper signal handling
func (m *RealProcessManager) setProcessGroup(cmd *exec.Cmd) {
	if runtime.GOOS != "windows" {
		// On Unix systems, create a new process group
		// Note: Setpgid is Unix-specific, handled via build tags if needed
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	// Windows doesn't need special process group handling
}

// tailFile reads the last N lines from a file
func (m *RealProcessManager) tailFile(filePath string, lines int) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	defer file.Close()

	// Read all lines (for simplicity)
	// In production, use more efficient tail algorithm for large files
	var allLines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		allLines = append(allLines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Return last N lines
	if len(allLines) <= lines {
		return allLines, nil
	}
	return allLines[len(allLines)-lines:], nil
}

// getProcessStats retrieves CPU and memory usage of a process
func (m *RealProcessManager) getProcessStats(pid int) (cpuUsage, memoryUsage float64) {
	// Basic implementation - returns 0 for now
	// In production, use platform-specific methods:
	// - Linux: Read /proc/[pid]/stat and /proc/[pid]/status
	// - Windows: Use performance counters or WMI
	// - Or use a library like gopsutil for cross-platform support

	// TODO: Implement with gopsutil or platform-specific code
	// import "github.com/shirou/gopsutil/v3/process"
	// p, err := process.NewProcess(int32(pid))
	// if err != nil {
	//     return 0, 0
	// }
	// cpuPercent, _ := p.CPUPercent()
	// memInfo, _ := p.MemoryInfo()
	// return cpuPercent, float64(memInfo.RSS) / 1024 / 1024 // MB

	return 0.0, 0.0
}
