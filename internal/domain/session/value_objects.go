package session

// ConnectionInfo value object
type ConnectionInfo struct {
	IPAddress string
	Port      int
	Protocol  string
	UserAgent string
}

func NewConnectionInfo(ipAddress string, port int, protocol, userAgent string) ConnectionInfo {
	return ConnectionInfo{
		IPAddress: ipAddress,
		Port:      port,
		Protocol:  protocol,
		UserAgent: userAgent,
	}
}

// Status enum
type Status string

const (
	StatusActive       Status = "active"
	StatusDisconnected Status = "disconnected"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusDisconnected:
		return true
	}
	return false
}
