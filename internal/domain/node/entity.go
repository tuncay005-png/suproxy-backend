package node

import (
	"time"

	"github.com/google/uuid"
)

const (
	UnlimitedBandwidth int64 = 0
	BytesInGB          int64 = 1024 * 1024 * 1024
)

// Node represents a VPN node running on a server
type Node struct {
	ID                  uuid.UUID
	ServerID            uuid.UUID
	Protocol            Protocol
	Port                int
	MaxUsers            int
	CurrentUsers        int
	BandwidthLimitBytes int64
	BandwidthUsedBytes  int64
	CPUUsage            float64
	RAMUsage            float64
	LatencyMs           int
	Version             string
	HealthStatus        HealthStatus
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func NewNode(serverID uuid.UUID, protocol Protocol, port, maxUsers int, bandwidthLimitGB int64) (*Node, error) {
	if serverID == uuid.Nil {
		return nil, ErrInvalidServerID
	}
	if !protocol.IsValid() {
		return nil, ErrInvalidProtocol
	}
	if port <= 0 || port > 65535 {
		return nil, ErrInvalidPort
	}
	if maxUsers <= 0 {
		return nil, ErrInvalidMaxUsers
	}
	if bandwidthLimitGB < 0 {
		return nil, ErrInvalidBandwidthLimit
	}

	bandwidthLimitBytes := bandwidthLimitGB * BytesInGB
	if bandwidthLimitGB == 0 {
		bandwidthLimitBytes = UnlimitedBandwidth
	}

	return &Node{
		ID:                  uuid.New(),
		ServerID:            serverID,
		Protocol:            protocol,
		Port:                port,
		MaxUsers:            maxUsers,
		CurrentUsers:        0,
		BandwidthLimitBytes: bandwidthLimitBytes,
		BandwidthUsedBytes:  0,
		CPUUsage:            0,
		RAMUsage:            0,
		LatencyMs:           0,
		Version:             "",
		HealthStatus:        HealthStatusUnknown,
		CreatedAt:           time.Now().UTC(),
		UpdatedAt:           time.Now().UTC(),
	}, nil
}

func (n *Node) CanAcceptUser() bool {
	if !n.IsHealthy() {
		return false
	}
	if n.IsOverloaded() {
		return false
	}
	return n.AvailableSlots() > 0
}

func (n *Node) IsHealthy() bool {
	return n.HealthStatus == HealthStatusHealthy
}

func (n *Node) IsOverloaded() bool {
	// CPU usage > 90% or RAM usage > 90%
	if n.CPUUsage > 90.0 || n.RAMUsage > 90.0 {
		return true
	}
	
	// Bandwidth limit exceeded (if not unlimited)
	if n.BandwidthLimitBytes > 0 && n.BandwidthUsedBytes >= n.BandwidthLimitBytes {
		return true
	}
	
	return false
}

func (n *Node) AvailableSlots() int {
	available := n.MaxUsers - n.CurrentUsers
	if available < 0 {
		return 0
	}
	return available
}

func (n *Node) UpdateMetrics(cpuUsage, ramUsage float64, latencyMs int) error {
	if cpuUsage < 0 || cpuUsage > 100 {
		return ErrInvalidCPUUsage
	}
	if ramUsage < 0 || ramUsage > 100 {
		return ErrInvalidRAMUsage
	}
	if latencyMs < 0 {
		return ErrInvalidLatency
	}

	n.CPUUsage = cpuUsage
	n.RAMUsage = ramUsage
	n.LatencyMs = latencyMs
	n.UpdatedAt = time.Now().UTC()

	// Auto-update health status based on metrics
	n.updateHealthStatus()
	
	return nil
}

func (n *Node) updateHealthStatus() {
	if n.CPUUsage > 90 || n.RAMUsage > 90 || n.LatencyMs > 1000 {
		n.HealthStatus = HealthStatusUnhealthy
	} else if n.CPUUsage > 70 || n.RAMUsage > 70 || n.LatencyMs > 500 {
		n.HealthStatus = HealthStatusDegraded
	} else {
		n.HealthStatus = HealthStatusHealthy
	}
}

func (n *Node) IncrementUsers() error {
	if n.CurrentUsers >= n.MaxUsers {
		return ErrNodeFull
	}
	n.CurrentUsers++
	n.UpdatedAt = time.Now().UTC()
	return nil
}

func (n *Node) DecrementUsers() error {
	if n.CurrentUsers <= 0 {
		return ErrInvalidUserCount
	}
	n.CurrentUsers--
	n.UpdatedAt = time.Now().UTC()
	return nil
}

func (n *Node) AddBandwidthUsage(bytes int64) error {
	if bytes < 0 {
		return ErrInvalidBandwidthUsage
	}
	n.BandwidthUsedBytes += bytes
	n.UpdatedAt = time.Now().UTC()
	return nil
}

func (n *Node) ResetBandwidthUsage() {
	n.BandwidthUsedBytes = 0
	n.UpdatedAt = time.Now().UTC()
}

func (n *Node) UpdateVersion(version string) {
	n.Version = version
	n.UpdatedAt = time.Now().UTC()
}

func (n *Node) UpdateMaxUsers(maxUsers int) error {
	if maxUsers <= 0 {
		return ErrInvalidMaxUsers
	}
	n.MaxUsers = maxUsers
	n.UpdatedAt = time.Now().UTC()
	return nil
}

func (n *Node) UpdateBandwidthLimit(limitGB int64) error {
	if limitGB < 0 {
		return ErrInvalidBandwidthLimit
	}
	
	if limitGB == 0 {
		n.BandwidthLimitBytes = UnlimitedBandwidth
	} else {
		n.BandwidthLimitBytes = limitGB * BytesInGB
	}
	
	n.UpdatedAt = time.Now().UTC()
	return nil
}

func (n *Node) HasUnlimitedBandwidth() bool {
	return n.BandwidthLimitBytes == UnlimitedBandwidth
}

func (n *Node) RemainingBandwidth() int64 {
	if n.HasUnlimitedBandwidth() {
		return 0 // 0 indicates unlimited
	}
	
	remaining := n.BandwidthLimitBytes - n.BandwidthUsedBytes
	if remaining < 0 {
		return 0
	}
	return remaining
}

func (n *Node) BandwidthUsagePercentage() float64 {
	if n.HasUnlimitedBandwidth() {
		return 0
	}
	if n.BandwidthLimitBytes == 0 {
		return 0
	}
	return float64(n.BandwidthUsedBytes) / float64(n.BandwidthLimitBytes) * 100
}

func (n *Node) UserLoadPercentage() float64 {
	if n.MaxUsers == 0 {
		return 0
	}
	return float64(n.CurrentUsers) / float64(n.MaxUsers) * 100
}
