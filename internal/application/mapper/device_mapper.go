package mapper

import (
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/domain/device"
)

func ToDeviceResponse(d *device.Device) *dto.DeviceResponse {
	if d == nil {
		return nil
	}

	return &dto.DeviceResponse{
		ID:            d.ID,
		UserID:        d.UserID,
		Name:          d.Name,
		DeviceType:    string(d.DeviceType),
		Identifier:    d.Identifier.String(),
		Status:        string(d.Status),
		LastSeenAt:    d.LastSeenAt,
		LastIPAddress: d.LastIPAddress,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}
}

func ToDeviceListResponse(devices []*device.Device, total int64) *dto.DeviceListResponse {
	responses := make([]*dto.DeviceResponse, 0, len(devices))
	for _, d := range devices {
		responses = append(responses, ToDeviceResponse(d))
	}

	return &dto.DeviceListResponse{
		Devices: responses,
		Total:   total,
	}
}
