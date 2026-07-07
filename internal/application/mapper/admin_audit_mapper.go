package mapper

import (
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/domain/audit"
)

// ToAdminAuditResponse converts audit log to response DTO
func ToAdminAuditResponse(log *audit.Log) dto.AdminAuditResponse {
	return dto.AdminAuditResponse{
		ID:         log.ID,
		UserID:     log.UserID,
		Action:     string(log.Action),
		EntityType: log.EntityType,
		EntityID:   log.EntityID,
		IPAddress:  log.IPAddress,
		UserAgent:  log.UserAgent,
		Metadata:   log.Metadata,
		CreatedAt:  log.CreatedAt,
	}
}

// ToAdminAuditListResponse converts audit log list to paginated response
func ToAdminAuditListResponse(logs []*audit.Log, total int64, offset, limit int) dto.AdminAuditListResponse {
	data := make([]dto.AdminAuditResponse, 0, len(logs))
	for _, log := range logs {
		data = append(data, ToAdminAuditResponse(log))
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return dto.AdminAuditListResponse{
		Data:       data,
		Total:      total,
		Offset:     offset,
		Limit:      limit,
		TotalPages: totalPages,
	}
}
