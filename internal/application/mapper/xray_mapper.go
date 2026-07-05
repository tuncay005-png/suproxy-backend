package mapper

import (
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/domain/xray"
)

func ToXrayInstanceResponse(x *xray.XrayInstance) *dto.XrayInstanceResponse {
	if x == nil {
		return nil
	}

	return &dto.XrayInstanceResponse{
		ID:            x.ID,
		NodeID:        x.NodeID,
		Version:       x.Version,
		Status:        string(x.Status),
		ConfigVersion: x.ConfigVersion,
		StartedAt:     x.StartedAt,
		StoppedAt:     x.StoppedAt,
		Uptime:        int64(x.GetUptime().Seconds()),
		IsRunning:     x.IsRunning(),
		CreatedAt:     x.CreatedAt,
		UpdatedAt:     x.UpdatedAt,
	}
}

func ToInboundResponse(i *xray.Inbound) *dto.InboundResponse {
	if i == nil {
		return nil
	}

	return &dto.InboundResponse{
		ID:             i.ID,
		XrayInstanceID: i.XrayInstanceID,
		Protocol:       string(i.Protocol),
		Port:           i.Port,
		Transport:      string(i.Transport),
		Security:       string(i.Security),
		Enabled:        i.Enabled,
		CreatedAt:      i.CreatedAt,
		UpdatedAt:      i.UpdatedAt,
	}
}

func ToClientResponse(c *xray.Client) *dto.ClientResponse {
	if c == nil {
		return nil
	}

	return &dto.ClientResponse{
		ID:        c.ID,
		InboundID: c.InboundID,
		UserID:    c.UserID,
		UUID:      c.UUID,
		Flow:      c.Flow,
		Email:     c.Email,
		Enabled:   c.Enabled,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func ToRealityConfigResponse(r *xray.RealityConfig) *dto.RealityConfigResponse {
	if r == nil {
		return nil
	}

	return &dto.RealityConfigResponse{
		ID:          r.ID,
		InboundID:   r.InboundID,
		PrivateKey:  r.PrivateKey,
		PublicKey:   r.PublicKey,
		ShortID:     r.ShortID,
		ServerName:  r.ServerName,
		Fingerprint: r.Fingerprint,
		SpiderX:     r.SpiderX,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}
