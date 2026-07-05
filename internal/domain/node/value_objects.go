package node

// Protocol represents the VPN protocol type
type Protocol string

const (
	ProtocolVMess      Protocol = "vmess"
	ProtocolVLESS      Protocol = "vless"
	ProtocolTrojan     Protocol = "trojan"
	ProtocolShadowsock Protocol = "shadowsocks"
)

func (p Protocol) IsValid() bool {
	switch p {
	case ProtocolVMess, ProtocolVLESS, ProtocolTrojan, ProtocolShadowsock:
		return true
	default:
		return false
	}
}

func (p Protocol) String() string {
	return string(p)
}

// HealthStatus represents the health status of a node
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusDegraded  HealthStatus = "degraded"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	HealthStatusUnknown   HealthStatus = "unknown"
)

func (h HealthStatus) IsValid() bool {
	switch h {
	case HealthStatusHealthy, HealthStatusDegraded, HealthStatusUnhealthy, HealthStatusUnknown:
		return true
	default:
		return false
	}
}

func (h HealthStatus) String() string {
	return string(h)
}
