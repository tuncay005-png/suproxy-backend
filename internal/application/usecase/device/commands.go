package device

import (
	"context"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/application/dto"
)

type CreateDeviceCommand interface {
	Execute(ctx context.Context, userID uuid.UUID, req *dto.CreateDeviceRequest) (*dto.DeviceResponse, error)
}

type UpdateDeviceCommand interface {
	Execute(ctx context.Context, deviceID uuid.UUID, req *dto.UpdateDeviceRequest) (*dto.DeviceResponse, error)
}

type ActivateDeviceCommand interface {
	Execute(ctx context.Context, deviceID uuid.UUID) error
}

type DeactivateDeviceCommand interface {
	Execute(ctx context.Context, deviceID uuid.UUID) error
}

type DeleteDeviceCommand interface {
	Execute(ctx context.Context, deviceID uuid.UUID) error
}
