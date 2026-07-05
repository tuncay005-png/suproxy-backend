package device

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	Name          string
	DeviceType    DeviceType
	Identifier    DeviceIdentifier
	Status        Status
	LastSeenAt    *time.Time
	LastIPAddress string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewDevice(userID uuid.UUID, name string, deviceType DeviceType, identifier DeviceIdentifier) (*Device, error) {
	if userID == uuid.Nil {
		return nil, ErrInvalidUserID
	}
	if name == "" {
		return nil, ErrInvalidDeviceName
	}

	return &Device{
		ID:         uuid.New(),
		UserID:     userID,
		Name:       name,
		DeviceType: deviceType,
		Identifier: identifier,
		Status:     StatusActive,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}, nil
}

func (d *Device) Activate() error {
	if d.Status == StatusActive {
		return ErrDeviceAlreadyActive
	}
	d.Status = StatusActive
	d.UpdatedAt = time.Now().UTC()
	return nil
}

func (d *Device) Deactivate() error {
	if d.Status == StatusInactive {
		return ErrDeviceAlreadyInactive
	}
	d.Status = StatusInactive
	d.UpdatedAt = time.Now().UTC()
	return nil
}

func (d *Device) UpdateLastSeen(ipAddress string) {
	now := time.Now().UTC()
	d.LastSeenAt = &now
	d.LastIPAddress = ipAddress
	d.UpdatedAt = now
}

func (d *Device) Rename(newName string) error {
	if newName == "" {
		return ErrInvalidDeviceName
	}
	d.Name = newName
	d.UpdatedAt = time.Now().UTC()
	return nil
}

func (d *Device) IsActive() bool {
	return d.Status == StatusActive
}
