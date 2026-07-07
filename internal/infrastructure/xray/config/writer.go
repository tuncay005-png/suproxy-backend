package config

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Writer handles writing, reading, and managing Xray config files
// Currently provides interface only - actual file operations will be implemented later
type Writer interface {
	// Write writes configuration to file
	Write(ctx context.Context, instanceID uuid.UUID, config []byte) error

	// Read reads configuration from file
	Read(ctx context.Context, instanceID uuid.UUID) ([]byte, error)

	// Backup creates a backup of current configuration
	Backup(ctx context.Context, instanceID uuid.UUID) error

	// Restore restores configuration from backup
	Restore(ctx context.Context, instanceID uuid.UUID, backupTime time.Time) error

	// Delete deletes configuration file
	Delete(ctx context.Context, instanceID uuid.UUID) error

	// DeleteBackup deletes a specific backup
	DeleteBackup(ctx context.Context, instanceID uuid.UUID, timestamp int64) error

	// GetPath returns the configuration file path
	GetPath(instanceID uuid.UUID) string

	// ListBackups lists all available backups
	ListBackups(ctx context.Context, instanceID uuid.UUID) ([]BackupInfo, error)
}

// BackupInfo represents backup metadata
type BackupInfo struct {
	InstanceID uuid.UUID
	Timestamp  time.Time
	Path       string
	Size       int64
}

// MockWriter provides a mock implementation for testing
// TODO: Replace with real file writer in production
type MockWriter struct {
	configs map[uuid.UUID][]byte
	backups map[uuid.UUID][]BackupInfo
}

// NewMockWriter creates a new mock writer for testing
func NewMockWriter() Writer {
	return &MockWriter{
		configs: make(map[uuid.UUID][]byte),
		backups: make(map[uuid.UUID][]BackupInfo),
	}
}

func (w *MockWriter) Write(ctx context.Context, instanceID uuid.UUID, config []byte) error {
	w.configs[instanceID] = config
	return nil
}

func (w *MockWriter) Read(ctx context.Context, instanceID uuid.UUID) ([]byte, error) {
	if config, ok := w.configs[instanceID]; ok {
		return config, nil
	}
	return nil, ErrInvalidConfig
}

func (w *MockWriter) Backup(ctx context.Context, instanceID uuid.UUID) error {
	if config, ok := w.configs[instanceID]; ok {
		backup := BackupInfo{
			InstanceID: instanceID,
			Timestamp:  time.Now().UTC(),
			Path:       "/mock/backup/path",
			Size:       int64(len(config)),
		}
		w.backups[instanceID] = append(w.backups[instanceID], backup)
		return nil
	}
	return ErrInvalidConfig
}

func (w *MockWriter) Restore(ctx context.Context, instanceID uuid.UUID, backupTime time.Time) error {
	// Mock implementation - in real implementation, restore from backup file
	return nil
}

func (w *MockWriter) Delete(ctx context.Context, instanceID uuid.UUID) error {
	delete(w.configs, instanceID)
	return nil
}

func (w *MockWriter) DeleteBackup(ctx context.Context, instanceID uuid.UUID, timestamp int64) error {
	// Mock implementation - delete backup from list
	if backups, ok := w.backups[instanceID]; ok {
		filtered := []BackupInfo{}
		for _, backup := range backups {
			if backup.Timestamp.Unix() != timestamp {
				filtered = append(filtered, backup)
			}
		}
		w.backups[instanceID] = filtered
	}
	return nil
}

func (w *MockWriter) GetPath(instanceID uuid.UUID) string {
	return "/etc/xray/" + instanceID.String() + ".json"
}

func (w *MockWriter) ListBackups(ctx context.Context, instanceID uuid.UUID) ([]BackupInfo, error) {
	if backups, ok := w.backups[instanceID]; ok {
		return backups, nil
	}
	return []BackupInfo{}, nil
}
