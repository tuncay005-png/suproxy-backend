package mapper

import (
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/domain/node"
)

func ToNodeResponse(n *node.Node) *dto.NodeResponse {
	if n == nil {
		return nil
	}

	const GBToBytes = 1024 * 1024 * 1024
	
	bandwidthLimitGB := float64(n.BandwidthLimitBytes) / float64(GBToBytes)
	bandwidthUsedGB := float64(n.BandwidthUsedBytes) / float64(GBToBytes)
	
	var remainingBandwidthGB float64
	if n.HasUnlimitedBandwidth() {
		remainingBandwidthGB = -1 // Indicates unlimited
	} else {
		remainingBytes := n.RemainingBandwidth()
		remainingBandwidthGB = float64(remainingBytes) / float64(GBToBytes)
	}

	return &dto.NodeResponse{
		ID:                    n.ID,
		ServerID:              n.ServerID,
		Protocol:              string(n.Protocol),
		Port:                  n.Port,
		MaxUsers:              n.MaxUsers,
		CurrentUsers:          n.CurrentUsers,
		AvailableSlots:        n.AvailableSlots(),
		UserLoadPercentage:    n.UserLoadPercentage(),
		BandwidthLimitBytes:   n.BandwidthLimitBytes,
		BandwidthUsedBytes:    n.BandwidthUsedBytes,
		BandwidthLimitGB:      bandwidthLimitGB,
		BandwidthUsedGB:       bandwidthUsedGB,
		RemainingBandwidthGB:  remainingBandwidthGB,
		BandwidthUsagePercent: n.BandwidthUsagePercentage(),
		HasUnlimitedBandwidth: n.HasUnlimitedBandwidth(),
		CPUUsage:              n.CPUUsage,
		RAMUsage:              n.RAMUsage,
		LatencyMs:             n.LatencyMs,
		Version:               n.Version,
		HealthStatus:          string(n.HealthStatus),
		IsHealthy:             n.IsHealthy(),
		IsOverloaded:          n.IsOverloaded(),
		CanAcceptUser:         n.CanAcceptUser(),
		CreatedAt:             n.CreatedAt,
		UpdatedAt:             n.UpdatedAt,
	}
}
