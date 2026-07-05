package xray

import "errors"

var (
	// XrayInstance errors
	ErrInstanceNotFound         = errors.New("xray instance not found")
	ErrInstanceAlreadyExists    = errors.New("xray instance already exists")
	ErrInstanceAlreadyRunning   = errors.New("xray instance already running")
	ErrInstanceAlreadyStopped   = errors.New("xray instance already stopped")
	ErrInvalidNodeID            = errors.New("invalid node id")
	ErrInvalidVersion           = errors.New("invalid version")
	ErrInvalidInstanceID        = errors.New("invalid instance id")

	// Inbound errors
	ErrInboundNotFound          = errors.New("inbound not found")
	ErrInboundAlreadyExists     = errors.New("inbound already exists")
	ErrInboundAlreadyEnabled    = errors.New("inbound already enabled")
	ErrInboundAlreadyDisabled   = errors.New("inbound already disabled")
	ErrInvalidInboundID         = errors.New("invalid inbound id")
	ErrInvalidProtocol          = errors.New("invalid protocol")
	ErrInvalidPort              = errors.New("invalid port")
	ErrInvalidTransport         = errors.New("invalid transport")
	ErrInvalidSecurity          = errors.New("invalid security")

	// Client errors
	ErrClientNotFound           = errors.New("client not found")
	ErrClientAlreadyExists      = errors.New("client already exists")
	ErrClientAlreadyEnabled     = errors.New("client already enabled")
	ErrClientAlreadyDisabled    = errors.New("client already disabled")
	ErrInvalidUserID            = errors.New("invalid user id")
	ErrInvalidUUID              = errors.New("invalid uuid")
	ErrInvalidEmail             = errors.New("invalid email")

	// Reality errors
	ErrRealityConfigNotFound    = errors.New("reality config not found")
	ErrRealityConfigExists      = errors.New("reality config already exists")
	ErrInvalidKeys              = errors.New("invalid keys")
	ErrInvalidServerName        = errors.New("invalid server name")
	ErrInvalidFingerprint       = errors.New("invalid fingerprint")
)
