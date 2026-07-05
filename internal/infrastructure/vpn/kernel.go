package vpn

import (
	"context"

	"github.com/google/uuid"
)

// Kernel represents a VPN core implementation (Xray, Sing-box, Hysteria, TUIC, etc.)
// This interface provides abstraction over different VPN kernels
type Kernel interface {
	// Name returns the kernel name (e.g., "xray", "singbox")
	Name() string

	// Version returns the current kernel version
	Version(ctx context.Context) (string, error)

	// GenerateConfig generates kernel-specific configuration
	GenerateConfig(ctx context.Context, instanceID uuid.UUID) ([]byte, error)

	// ValidateConfig validates the generated configuration
	ValidateConfig(ctx context.Context, config []byte) error

	// Start starts the kernel process
	Start(ctx context.Context, instanceID uuid.UUID) error

	// Stop stops the kernel process
	Stop(ctx context.Context, instanceID uuid.UUID) error

	// Restart restarts the kernel process
	Restart(ctx context.Context, instanceID uuid.UUID) error

	// Reload reloads configuration without restarting
	Reload(ctx context.Context, instanceID uuid.UUID) error

	// Status returns the current status of the kernel instance
	Status(ctx context.Context, instanceID uuid.UUID) (KernelStatus, error)

	// IsRunning checks if the kernel instance is running
	IsRunning(ctx context.Context, instanceID uuid.UUID) (bool, error)

	// GetProcessID returns the OS process ID if running
	GetProcessID(ctx context.Context, instanceID uuid.UUID) (int, error)
}

// KernelStatus represents the runtime status of a kernel instance
type KernelStatus struct {
	Running     bool
	ProcessID   int
	Uptime      int64 // seconds
	ConfigPath  string
	LogPath     string
	ErrorReason string
}

// KernelFactory creates kernel instances based on type
type KernelFactory interface {
	Create(kernelType string) (Kernel, error)
	SupportedKernels() []string
}
