package mapper

import (
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/domain/server"
)

func ToServerResponse(s *server.Server, nodeCount int) *dto.ServerResponse {
	if s == nil {
		return nil
	}

	return &dto.ServerResponse{
		ID:        s.ID,
		Name:      s.Name,
		Country:   s.Country,
		City:      s.City,
		Hostname:  s.Hostname,
		Provider:  s.Provider,
		IPv4:      s.IPv4,
		IPv6:      s.IPv6,
		Status:    string(s.Status),
		IsPublic:  s.IsPublic,
		IsOnline:  s.IsOnline(),
		NodeCount: nodeCount,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
