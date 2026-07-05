package node

import "time"

// Protocol enum
type Protocol string

const (
	ProtocolVMess       Protocol = "vmess"
	ProtocolVLESS       Protocol = "vless"
	ProtocolTrojan      Protocol = "trojan"
	ProtocolShadowsocks Protocol = "shadowsocks"
)

func (p Protocol) IsValid() bool {
	switch p {
	case ProtocolVMess, ProtocolVLESS, ProtocolTrojan, ProtocolShadowsocks:
		return true
	}
	return false
}

// Configuration value object
type Configuration struct {
	UUID      string
	AlterID   int
	Security  string
	Network   string
	TLS       bool
	SNI       string
	Path      string
	Host      string
}

func NewConfiguration(uuid, security, network string, tls bool) Configuration {
	return Configuration{
		UUID:     uuid,
		Security: security,
		Network:  network,
		TLS:      tls,
	}
}

// Metrics value object
type Metrics struct {
	TotalConnections   int64
	ActiveConnections  int64
	TotalBandwidthIn   int64
	TotalBandwidthOut  int64
	LastCheckedAt      *time.Time
}

func NewMetrics() Metrics {
	return Metrics{
		TotalConnections:  0,
		ActiveConnections: 0,
		TotalBandwidthIn:  0,
		TotalBandwidthOut: 0,
	}
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
