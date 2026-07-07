package mapper

import (
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/domain/xray"
)

// ToAdminXrayInstanceResponse converts domain XrayInstance to admin response
func ToAdminXrayInstanceResponse(instance *xray.XrayInstance) *dto.AdminXrayInstanceResponse {
	if instance == nil {
		return nil
	}

	uptime := int64(0)
	if instance.IsRunning() {
		uptime = int64(instance.GetUptime().Seconds())
	}

	return &dto.AdminXrayInstanceResponse{
		ID:            instance.ID,
		NodeID:        instance.NodeID,
		Version:       instance.Version,
		Status:        string(instance.Status),
		ConfigVersion: instance.ConfigVersion,
		StartedAt:     instance.StartedAt,
		StoppedAt:     instance.StoppedAt,
		Uptime:        uptime,
		CreatedAt:     instance.CreatedAt,
		UpdatedAt:     instance.UpdatedAt,
	}
}

// ToAdminXrayInstanceListResponse converts domain instances to admin list response
func ToAdminXrayInstanceListResponse(instances []*xray.XrayInstance, total int64, offset, limit int) *dto.AdminXrayInstanceListResponse {
	responses := make([]*dto.AdminXrayInstanceResponse, 0, len(instances))
	for _, instance := range instances {
		responses = append(responses, ToAdminXrayInstanceResponse(instance))
	}

	return &dto.AdminXrayInstanceListResponse{
		Instances: responses,
		Total:     total,
		Offset:    offset,
		Limit:     limit,
	}
}
