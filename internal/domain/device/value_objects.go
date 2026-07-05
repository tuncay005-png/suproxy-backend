package device

import "errors"

// DeviceIdentifier value object
type DeviceIdentifier struct {
	value string
}

func NewDeviceIdentifier(identifier string) (DeviceIdentifier, error) {
	if identifier == "" {
		return DeviceIdentifier{}, errors.New("device identifier cannot be empty")
	}
	return DeviceIdentifier{value: identifier}, nil
}

func (d DeviceIdentifier) String() string {
	return d.value
}

// DeviceType enum
type DeviceType string

const (
	DeviceTypeAndroid DeviceType = "android"
	DeviceTypeIOS     DeviceType = "ios"
	DeviceTypeWindows DeviceType = "windows"
	DeviceTypeMacOS   DeviceType = "macos"
	DeviceTypeLinux   DeviceType = "linux"
	DeviceTypeOther   DeviceType = "other"
)

func (dt DeviceType) IsValid() bool {
	switch dt {
	case DeviceTypeAndroid, DeviceTypeIOS, DeviceTypeWindows, DeviceTypeMacOS, DeviceTypeLinux, DeviceTypeOther:
		return true
	}
	return false
}

// Status enum
type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusInactive:
		return true
	}
	return false
}
