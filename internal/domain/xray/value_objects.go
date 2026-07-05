package xray

// InstanceStatus represents the status of an Xray instance
type InstanceStatus string

const (
	StatusRunning InstanceStatus = "running"
	StatusStopped InstanceStatus = "stopped"
	StatusFailed  InstanceStatus = "failed"
)

func (s InstanceStatus) IsValid() bool {
	switch s {
	case StatusRunning, StatusStopped, StatusFailed:
		return true
	default:
		return false
	}
}

func (s InstanceStatus) String() string {
	return string(s)
}

// InboundProtocol represents the protocol type for inbound
type InboundProtocol string

const (
	ProtocolVMess      InboundProtocol = "vmess"
	ProtocolVLESS      InboundProtocol = "vless"
	ProtocolTrojan     InboundProtocol = "trojan"
	ProtocolShadowsock InboundProtocol = "shadowsocks"
)

func (p InboundProtocol) IsValid() bool {
	switch p {
	case ProtocolVMess, ProtocolVLESS, ProtocolTrojan, ProtocolShadowsock:
		return true
	default:
		return false
	}
}

func (p InboundProtocol) String() string {
	return string(p)
}

// TransportType represents the transport protocol
type TransportType string

const (
	TransportTCP        TransportType = "tcp"
	TransportWebSocket  TransportType = "ws"
	TransportHTTP       TransportType = "http"
	TransportQUIC       TransportType = "quic"
	TransportGRPC       TransportType = "grpc"
	TransportHTTPUpgrade TransportType = "httpupgrade"
)

func (t TransportType) IsValid() bool {
	switch t {
	case TransportTCP, TransportWebSocket, TransportHTTP, TransportQUIC, TransportGRPC, TransportHTTPUpgrade:
		return true
	default:
		return false
	}
}

func (t TransportType) String() string {
	return string(t)
}

// SecurityType represents the security/encryption type
type SecurityType string

const (
	SecurityNone    SecurityType = "none"
	SecurityTLS     SecurityType = "tls"
	SecurityREALITY SecurityType = "reality"
)

func (s SecurityType) IsValid() bool {
	switch s {
	case SecurityNone, SecurityTLS, SecurityREALITY:
		return true
	default:
		return false
	}
}

func (s SecurityType) String() string {
	return string(s)
}
