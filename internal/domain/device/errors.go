package device

import "errors"

var (
	ErrDeviceNotFound        = errors.New("device not found")
	ErrDeviceAlreadyExists   = errors.New("device already exists")
	ErrDeviceAlreadyActive   = errors.New("device already active")
	ErrDeviceAlreadyInactive = errors.New("device already inactive")
	ErrInvalidUserID         = errors.New("invalid user id")
	ErrInvalidDeviceName     = errors.New("invalid device name")
	ErrInvalidDeviceType     = errors.New("invalid device type")
	ErrMaxDevicesReached     = errors.New("maximum number of devices reached")
)
