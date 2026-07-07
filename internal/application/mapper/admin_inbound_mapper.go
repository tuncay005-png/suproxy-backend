package mapper

import (
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/domain/xray"
)

// ToAdminInboundResponse converts domain Inbound to admin response
func ToAdminInboundResponse(inbound *xray.Inbound) *dto.AdminInboundResponse {
	if inbound == nil {
		return nil
	}

	return &dto.AdminInboundResponse{
		ID:             inbound.ID,
		XrayInstanceID: inbound.XrayInstanceID,
		Protocol:       string(inbound.Protocol),
		Port:           inbound.Port,
		Transport:      string(inbound.Transport),
		Security:       string(inbound.Security),
		Enabled:        inbound.Enabled,
		CreatedAt:      inbound.CreatedAt,
		UpdatedAt:      inbound.UpdatedAt,
	}
}

// ToAdminInboundListResponse converts domain inbounds to admin list response
func ToAdminInboundListResponse(inbounds []*xray.Inbound, total int64, offset, limit int) *dto.AdminInboundListResponse {
	responses := make([]*dto.AdminInboundResponse, 0, len(inbounds))
	for _, inbound := range inbounds {
		responses = append(responses, ToAdminInboundResponse(inbound))
	}

	return &dto.AdminInboundListResponse{
		Inbounds: responses,
		Total:    total,
		Offset:   offset,
		Limit:    limit,
	}
}
