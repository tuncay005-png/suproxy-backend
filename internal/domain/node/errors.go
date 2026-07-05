package node

import "errors"

var (
	ErrNodeNotFound          = errors.New("node not found")
	ErrNodeAlreadyExists     = errors.New("node already exists")
	ErrInvalidServerID       = errors.New("invalid server id")
	ErrInvalidProtocol       = errors.New("invalid protocol")
	ErrInvalidPort           = errors.New("invalid port")
	ErrInvalidMaxUsers       = errors.New("invalid max users")
	ErrInvalidBandwidthLimit = errors.New("invalid bandwidth limit")
	ErrInvalidBandwidthUsage = errors.New("invalid bandwidth usage")
	ErrInvalidCPUUsage       = errors.New("invalid CPU usage")
	ErrInvalidRAMUsage       = errors.New("invalid RAM usage")
	ErrInvalidLatency        = errors.New("invalid latency")
	ErrInvalidUserCount      = errors.New("invalid user count")
	ErrNodeFull              = errors.New("node is full")
	ErrNodeNotHealthy        = errors.New("node is not healthy")
	ErrNodeOverloaded        = errors.New("node is overloaded")
)
