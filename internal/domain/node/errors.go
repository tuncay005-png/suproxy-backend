package node

import "errors"

var (
	ErrNodeNotFound         = errors.New("node not found")
	ErrNodeAlreadyExists    = errors.New("node already exists")
	ErrNodeAlreadyActive    = errors.New("node already active")
	ErrNodeAlreadyInactive  = errors.New("node already inactive")
	ErrInvalidServerID      = errors.New("invalid server id")
	ErrInvalidProtocol      = errors.New("invalid protocol")
	ErrInvalidPort          = errors.New("invalid port")
	ErrInvalidConfiguration = errors.New("invalid configuration")
)
