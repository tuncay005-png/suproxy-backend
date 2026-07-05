package binary

import (
	"context"
	"errors"
	"time"
)

var (
	ErrBinaryNotFound     = errors.New("xray binary not found")
	ErrInvalidVersion     = errors.New("invalid version")
	ErrDownloadFailed     = errors.New("download failed")
	ErrVerificationFailed = errors.New("verification failed")
)

// Manager handles Xray binary management
// Currently provides interface only - actual binary operations will be implemented later
type Manager interface {
	// Detect detects Xray binary in system
	Detect(ctx context.Context) (string, error)

	// Validate validates the binary is executable and correct
	Validate(ctx context.Context, binaryPath string) error

	// CurrentVersion returns the currently installed version
	CurrentVersion(ctx context.Context) (string, error)

	// LatestVersion returns the latest available version
	LatestVersion(ctx context.Context) (string, error)

	// Download downloads a specific version
	Download(ctx context.Context, version string) error

	// Upgrade upgrades to a specific version
	Upgrade(ctx context.Context, version string) error

	// GetPath returns the binary path
	GetPath() string

	// IsInstalled checks if Xray is installed
	IsInstalled(ctx context.Context) bool
}

// BinaryInfo represents Xray binary information
type BinaryInfo struct {
	Path       string
	Version    string
	Size       int64
	ModifiedAt time.Time
	Executable bool
	Verified   bool
}

// MockManager provides a mock implementation for testing
// TODO: Replace with real binary manager using os/exec and HTTP downloads in production
type MockManager struct {
	binaryPath     string
	currentVersion string
	installed      bool
}

// NewMockManager creates a new mock binary manager for testing
func NewMockManager() Manager {
	return &MockManager{
		binaryPath:     "/usr/local/bin/xray",
		currentVersion: "1.8.7",
		installed:      true,
	}
}

func (m *MockManager) Detect(ctx context.Context) (string, error) {
	if m.installed {
		return m.binaryPath, nil
	}
	return "", ErrBinaryNotFound
}

func (m *MockManager) Validate(ctx context.Context, binaryPath string) error {
	// Mock implementation - in real implementation, check file permissions and execute --version
	if !m.installed {
		return ErrBinaryNotFound
	}
	return nil
}

func (m *MockManager) CurrentVersion(ctx context.Context) (string, error) {
	if !m.installed {
		return "", ErrBinaryNotFound
	}
	// Mock implementation - in real implementation, execute: xray version
	return m.currentVersion, nil
}

func (m *MockManager) LatestVersion(ctx context.Context) (string, error) {
	// Mock implementation - in real implementation, fetch from GitHub API
	return "1.8.8", nil
}

func (m *MockManager) Download(ctx context.Context, version string) error {
	// Mock implementation - in real implementation:
	// 1. Construct download URL based on version and OS/arch
	// 2. Download from GitHub releases
	// 3. Verify checksum
	// 4. Extract archive
	// 5. Set executable permissions
	m.currentVersion = version
	return nil
}

func (m *MockManager) Upgrade(ctx context.Context, version string) error {
	// Mock implementation - in real implementation:
	// 1. Download new version
	// 2. Stop running instances
	// 3. Backup old binary
	// 4. Replace binary
	// 5. Restart instances
	if err := m.Download(ctx, version); err != nil {
		return err
	}
	m.currentVersion = version
	return nil
}

func (m *MockManager) GetPath() string {
	return m.binaryPath
}

func (m *MockManager) IsInstalled(ctx context.Context) bool {
	return m.installed
}

// DownloadOptions represents download configuration
type DownloadOptions struct {
	Version        string
	TargetPath     string
	VerifyChecksum bool
	OS             string
	Arch           string
}

// UpgradeOptions represents upgrade configuration
type UpgradeOptions struct {
	Version        string
	BackupOld      bool
	RestartRunning bool
	Force          bool
}
