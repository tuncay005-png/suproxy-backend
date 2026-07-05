package device

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
)

type GetDeviceQuery interface {
	Execute(ctx context.Context, deviceID uuid.UUID) (*dto.DeviceResponse, error)
}

type ListDevicesByUserQuery interface {
	Execute(ctx context.Context, userID uuid.UUID) (*dto.DeviceListResponse, error)
}
