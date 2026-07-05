package server

import "errors"

var (
	ErrServerNotFound              = errors.New("server not found")
	ErrServerAlreadyExists         = errors.New("server already exists")
	ErrServerAlreadyActive         = errors.New("server already active")
	ErrServerAlreadyInactive       = errors.New("server already inactive")
	ErrServerAlreadyInMaintenance  = errors.New("server already in maintenance")
	ErrInvalidServerName           = errors.New("invalid server name")
	ErrInvalidIPAddress            = errors.New("invalid ip address")
	ErrInvalidPort                 = errors.New("invalid port")
	ErrServerCapacityFull          = errors.New("server capacity full")
	ErrServerNotAvailable          = errors.New("server not available")
)
