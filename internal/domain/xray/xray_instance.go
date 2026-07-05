package xray

import (
	"time"

	"github.com/google/uuid"
)

// XrayInstance represents a running Xray-core instance on a node
type XrayInstance struct {
	ID            uuid.UUID
	NodeID        uuid.UUID
	Version       string
	Status        InstanceStatus
	ConfigVersion int
	StartedAt     *time.Time
	StoppedAt     *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewXrayInstance(nodeID uuid.UUID, version string) (*XrayInstance, error) {
	if nodeID == uuid.Nil {
		return nil, ErrInvalidNodeID
	}
	if version == "" {
		return nil, ErrInvalidVersion
	}

	return &XrayInstance{
		ID:            uuid.New(),
		NodeID:        nodeID,
		Version:       version,
		Status:        StatusStopped,
		ConfigVersion: 1,
		StartedAt:     nil,
		StoppedAt:     nil,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}, nil
}

func (x *XrayInstance) Start() error {
	if x.Status == StatusRunning {
		return ErrInstanceAlreadyRunning
	}

	now := time.Now().UTC()
	x.Status = StatusRunning
	x.StartedAt = &now
	x.StoppedAt = nil
	x.UpdatedAt = now
	return nil
}

func (x *XrayInstance) Stop() error {
	if x.Status == StatusStopped {
		return ErrInstanceAlreadyStopped
	}

	now := time.Now().UTC()
	x.Status = StatusStopped
	x.StoppedAt = &now
	x.UpdatedAt = now
	return nil
}

func (x *XrayInstance) Restart() error {
	now := time.Now().UTC()
	x.Status = StatusRunning
	x.StartedAt = &now
	x.StoppedAt = nil
	x.ConfigVersion++
	x.UpdatedAt = now
	return nil
}

func (x *XrayInstance) IsRunning() bool {
	return x.Status == StatusRunning
}

func (x *XrayInstance) SetFailed(reason string) {
	x.Status = StatusFailed
	x.UpdatedAt = time.Now().UTC()
}

func (x *XrayInstance) UpdateVersion(version string) error {
	if version == "" {
		return ErrInvalidVersion
	}
	x.Version = version
	x.UpdatedAt = time.Now().UTC()
	return nil
}

func (x *XrayInstance) IncrementConfigVersion() {
	x.ConfigVersion++
	x.UpdatedAt = time.Now().UTC()
}

func (x *XrayInstance) GetUptime() time.Duration {
	if x.StartedAt == nil || x.Status != StatusRunning {
		return 0
	}
	return time.Since(*x.StartedAt)
}
