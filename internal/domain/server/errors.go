package server

import "errors"

var (
	ErrServerNotFound          = errors.New("server not found")
	ErrServerAlreadyExists     = errors.New("server already exists")
	ErrServerAlreadyActive     = errors.New("server already active")
	ErrServerAlreadyInactive   = errors.New("server already inactive")
	ErrInvalidServerName       = errors.New("invalid server name")
	ErrInvalidCountry          = errors.New("invalid country")
	ErrInvalidCity             = errors.New("invalid city")
	ErrInvalidHostname         = errors.New("invalid hostname")
	ErrInvalidProvider         = errors.New("invalid provider")
	ErrInvalidIPAddress        = errors.New("invalid IP address")
	ErrServerNotAvailable      = errors.New("server not available")
)
