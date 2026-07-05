package session

import "errors"

var (
	ErrSessionNotFound            = errors.New("session not found")
	ErrSessionAlreadyDisconnected = errors.New("session already disconnected")
	ErrInvalidUserID              = errors.New("invalid user id")
	ErrInvalidDeviceID            = errors.New("invalid device id")
	ErrInvalidNodeID              = errors.New("invalid node id")
	ErrInvalidSubscriptionID      = errors.New("invalid subscription id")
)
