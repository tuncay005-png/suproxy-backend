package mapper

import (
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/domain/xray"
)

// ToAdminClientResponse converts domain Client to admin response
func ToAdminClientResponse(client *xray.Client) *dto.AdminClientResponse {
	if client == nil {
		return nil
	}

	return &dto.AdminClientResponse{
		ID:        client.ID,
		InboundID: client.InboundID,
		UserID:    client.UserID,
		UUID:      client.UUID,
		Flow:      client.Flow,
		Email:     client.Email,
		Enabled:   client.Enabled,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	}
}

// ToAdminClientListResponse converts domain clients to admin list response
func ToAdminClientListResponse(clients []*xray.Client, total int64, offset, limit int) *dto.AdminClientListResponse {
	responses := make([]*dto.AdminClientResponse, 0, len(clients))
	for _, client := range clients {
		responses = append(responses, ToAdminClientResponse(client))
	}

	return &dto.AdminClientListResponse{
		Clients: responses,
		Total:   total,
		Offset:  offset,
		Limit:   limit,
	}
}
