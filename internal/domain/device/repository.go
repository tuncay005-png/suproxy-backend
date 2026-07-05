package device

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, device *Device) error
	FindByID(ctx context.Context, id uuid.UUID) (*Device, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*Device, error)
	FindByIdentifier(ctx context.Context, identifier DeviceIdentifier) (*Device, error)
	Update(ctx context.Context, device *Device) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	ExistsByIdentifier(ctx context.Context, identifier DeviceIdentifier) (bool, error)
}
