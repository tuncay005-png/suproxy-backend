package runtime

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Manager handles Xray process lifecycle management
// Currently provides interface only - actual process management will be implemented later
type Manager interface {
	// Start starts the Xray process
	Start(ctx context.Context, instanceID uuid.UUID) error

	// Stop stops the Xray process
	Stop(ctx context.Context, instanceID uuid.UUID) error

	// Restart restarts the Xray process
	Restart(ctx context.Context, instanceID uuid.UUID) error

	// Reload reloads configuration without restarting (hot reload)
	Reload(ctx context.Context, instanceID uuid.UUID) error

	// Status returns the current process status
	Status(ctx context.Context, instanceID uuid.UUID) (*ProcessStatus, error)

	// IsRunning checks if the process is running
	IsRunning(ctx context.Context, instanceID uuid.UUID) (bool, error)

	// GetProcessID returns the OS process ID
	GetProcessID(ctx context.Context, instanceID uuid.UUID) (int, error)

	// GetLogs retrieves recent logs
	GetLogs(ctx context.Context, instanceID uuid.UUID, lines int) ([]string, error)

	// Kill forcefully terminates the process
	Kill(ctx context.Context, instanceID uuid.UUID) error
}

// ProcessStatus represents the current status of an Xray process
type ProcessStatus struct {
	InstanceID  uuid.UUID
	Running     bool
	ProcessID   int
	StartedAt   time.Time
	Uptime      time.Duration
	ConfigPath  string
	LogPath     string
	ErrorPath   string
	CPUUsage    float64
	MemoryUsage int64 // bytes
	ErrorReason string
}

// MockManager provides a mock implementation for testing
// TODO: Replace with real process manager using os/exec in production
type MockManager struct {
	processes map[uuid.UUID]*ProcessStatus
}

// NewMockManager creates a new mock manager for testing
func NewMockManager() Manager {
	return &MockManager{
		processes: make(map[uuid.UUID]*ProcessStatus),
	}
}

func (m *MockManager) Start(ctx context.Context, instanceID uuid.UUID) error {
	m.processes[instanceID] = &ProcessStatus{
		InstanceID: instanceID,
		Running:    true,
		ProcessID:  12345, // Mock PID
		StartedAt:  time.Now().UTC(),
		ConfigPath: "/etc/xray/" + instanceID.String() + ".json",
		LogPath:    "/var/log/xray/" + instanceID.String() + ".log",
		ErrorPath:  "/var/log/xray/" + instanceID.String() + ".error.log",
	}
	return nil
}

func (m *MockManager) Stop(ctx context.Context, instanceID uuid.UUID) error {
	if status, ok := m.processes[instanceID]; ok {
		status.Running = false
		status.ProcessID = 0
	}
	return nil
}

func (m *MockManager) Restart(ctx context.Context, instanceID uuid.UUID) error {
	if err := m.Stop(ctx, instanceID); err != nil {
		return err
	}
	return m.Start(ctx, instanceID)
}

func (m *MockManager) Reload(ctx context.Context, instanceID uuid.UUID) error {
	// Mock implementation - in real implementation, send SIGHUP or use API
	return nil
}

func (m *MockManager) Status(ctx context.Context, instanceID uuid.UUID) (*ProcessStatus, error) {
	if status, ok := m.processes[instanceID]; ok {
		// Update uptime
		if status.Running {
			status.Uptime = time.Since(status.StartedAt)
		}
		return status, nil
	}
	return &ProcessStatus{
		InstanceID: instanceID,
		Running:    false,
	}, nil
}

func (m *MockManager) IsRunning(ctx context.Context, instanceID uuid.UUID) (bool, error) {
	status, err := m.Status(ctx, instanceID)
	if err != nil {
		return false, err
	}
	return status.Running, nil
}

func (m *MockManager) GetProcessID(ctx context.Context, instanceID uuid.UUID) (int, error) {
	status, err := m.Status(ctx, instanceID)
	if err != nil {
		return 0, err
	}
	return status.ProcessID, nil
}

func (m *MockManager) GetLogs(ctx context.Context, instanceID uuid.UUID, lines int) ([]string, error) {
	// Mock implementation - in real implementation, tail log files
	return []string{
		"[Info] Xray started",
		"[Info] Configuration loaded",
		"[Info] Listening on port 443",
	}, nil
}

func (m *MockManager) Kill(ctx context.Context, instanceID uuid.UUID) error {
	return m.Stop(ctx, instanceID)
}
