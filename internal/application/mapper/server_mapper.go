package mapper

import (
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/domain/server"
)

func ToServerResponse(s *server.Server) *dto.ServerResponse {
	if s == nil {
		return nil
	}

	return &dto.ServerResponse{
		ID:               s.ID,
		Name:             s.Name,
		Country:          s.Location.Country,
		City:             s.Location.City,
		CountryCode:      s.Location.CountryCode,
		IPAddress:        s.IPAddress,
		Domain:           s.Domain,
		Port:             s.Port,
		MaxUsers:         s.Capacity.MaxUsers,
		CurrentUsers:     s.Capacity.CurrentUsers,
		MaxBandwidth:     s.Capacity.MaxBandwidth,
		UsedBandwidth:    s.Capacity.UsedBandwidth,
		Status:           string(s.Status),
		UsagePercentage:  s.Capacity.UsagePercentage(),
		Tags:             s.Tags,
		CreatedAt:        s.CreatedAt,
		UpdatedAt:        s.UpdatedAt,
	}
}

func ToServerListResponse(servers []*server.Server, total int64, offset, limit int) *dto.ServerListResponse {
	responses := make([]*dto.ServerResponse, 0, len(servers))
	for _, s := range servers {
		responses = append(responses, ToServerResponse(s))
	}

	return &dto.ServerListResponse{
		Servers: responses,
		Total:   total,
		Offset:  offset,
		Limit:   limit,
	}
}
