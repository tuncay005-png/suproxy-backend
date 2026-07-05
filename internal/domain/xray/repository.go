package xray

import (
	"context"

	"github.com/google/uuid"
)

// XrayInstanceRepository manages XrayInstance persistence
type XrayInstanceRepository interface {
	Create(ctx context.Context, instance *XrayInstance) error
	FindByID(ctx context.Context, id uuid.UUID) (*XrayInstance, error)
	FindByNodeID(ctx context.Context, nodeID uuid.UUID) (*XrayInstance, error)
	FindRunning(ctx context.Context) ([]*XrayInstance, error)
	Update(ctx context.Context, instance *XrayInstance) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*XrayInstance, error)
	Count(ctx context.Context) (int64, error)
}

// InboundRepository manages Inbound persistence
type InboundRepository interface {
	Create(ctx context.Context, inbound *Inbound) error
	FindByID(ctx context.Context, id uuid.UUID) (*Inbound, error)
	FindByInstanceID(ctx context.Context, instanceID uuid.UUID) ([]*Inbound, error)
	FindEnabledByInstanceID(ctx context.Context, instanceID uuid.UUID) ([]*Inbound, error)
	Update(ctx context.Context, inbound *Inbound) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*Inbound, error)
	Count(ctx context.Context) (int64, error)
}

// ClientRepository manages Client persistence
type ClientRepository interface {
	Create(ctx context.Context, client *Client) error
	FindByID(ctx context.Context, id uuid.UUID) (*Client, error)
	FindByInboundID(ctx context.Context, inboundID uuid.UUID) ([]*Client, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*Client, error)
	FindByUUID(ctx context.Context, clientUUID string) (*Client, error)
	FindEnabledByInboundID(ctx context.Context, inboundID uuid.UUID) ([]*Client, error)
	Update(ctx context.Context, client *Client) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*Client, error)
	Count(ctx context.Context) (int64, error)
}

// RealityConfigRepository manages RealityConfig persistence
type RealityConfigRepository interface {
	Create(ctx context.Context, config *RealityConfig) error
	FindByID(ctx context.Context, id uuid.UUID) (*RealityConfig, error)
	FindByInboundID(ctx context.Context, inboundID uuid.UUID) (*RealityConfig, error)
	Update(ctx context.Context, config *RealityConfig) error
	Delete(ctx context.Context, id uuid.UUID) error
}
