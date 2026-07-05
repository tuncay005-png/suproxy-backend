package xray

import (
	"time"

	"github.com/google/uuid"
)

// Inbound represents an Xray inbound configuration
type Inbound struct {
	ID              uuid.UUID
	XrayInstanceID  uuid.UUID
	Protocol        InboundProtocol
	Port            int
	Transport       TransportType
	Security        SecurityType
	Enabled         bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewInbound(instanceID uuid.UUID, protocol InboundProtocol, port int, transport TransportType, security SecurityType) (*Inbound, error) {
	if instanceID == uuid.Nil {
		return nil, ErrInvalidInstanceID
	}
	if !protocol.IsValid() {
		return nil, ErrInvalidProtocol
	}
	if port <= 0 || port > 65535 {
		return nil, ErrInvalidPort
	}
	if !transport.IsValid() {
		return nil, ErrInvalidTransport
	}
	if !security.IsValid() {
		return nil, ErrInvalidSecurity
	}

	return &Inbound{
		ID:             uuid.New(),
		XrayInstanceID: instanceID,
		Protocol:       protocol,
		Port:           port,
		Transport:      transport,
		Security:       security,
		Enabled:        true,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}, nil
}

func (i *Inbound) Enable() error {
	if i.Enabled {
		return ErrInboundAlreadyEnabled
	}
	i.Enabled = true
	i.UpdatedAt = time.Now().UTC()
	return nil
}

func (i *Inbound) Disable() error {
	if !i.Enabled {
		return ErrInboundAlreadyDisabled
	}
	i.Enabled = false
	i.UpdatedAt = time.Now().UTC()
	return nil
}

func (i *Inbound) ChangePort(newPort int) error {
	if newPort <= 0 || newPort > 65535 {
		return ErrInvalidPort
	}
	i.Port = newPort
	i.UpdatedAt = time.Now().UTC()
	return nil
}

func (i *Inbound) IsEnabled() bool {
	return i.Enabled
}

func (i *Inbound) UpdateTransport(transport TransportType) error {
	if !transport.IsValid() {
		return ErrInvalidTransport
	}
	i.Transport = transport
	i.UpdatedAt = time.Now().UTC()
	return nil
}

func (i *Inbound) UpdateSecurity(security SecurityType) error {
	if !security.IsValid() {
		return ErrInvalidSecurity
	}
	i.Security = security
	i.UpdatedAt = time.Now().UTC()
	return nil
}
