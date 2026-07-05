package node

import (
	"time"

	"github.com/google/uuid"
)

type Node struct {
	ID            uuid.UUID
	ServerID      uuid.UUID
	Protocol      Protocol
	Configuration Configuration
	Port          int
	Status        Status
	Metrics       Metrics
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewNode(serverID uuid.UUID, protocol Protocol, config Configuration, port int) (*Node, error) {
	if serverID == uuid.Nil {
		return nil, ErrInvalidServerID
	}
	if port <= 0 || port > 65535 {
		return nil, ErrInvalidPort
	}

	return &Node{
		ID:            uuid.New(),
		ServerID:      serverID,
		Protocol:      protocol,
		Configuration: config,
		Port:          port,
		Status:        StatusActive,
		Metrics:       NewMetrics(),
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}, nil
}

func (n *Node) Activate() error {
	if n.Status == StatusActive {
		return ErrNodeAlreadyActive
	}
	n.Status = StatusActive
	n.UpdatedAt = time.Now().UTC()
	return nil
}

func (n *Node) Deactivate() error {
	if n.Status == StatusInactive {
		return ErrNodeAlreadyInactive
	}
	n.Status = StatusInactive
	n.UpdatedAt = time.Now().UTC()
	return nil
}

func (n *Node) UpdateConfiguration(config Configuration) {
	n.Configuration = config
	n.UpdatedAt = time.Now().UTC()
}

func (n *Node) UpdateMetrics(metrics Metrics) {
	n.Metrics = metrics
	n.UpdatedAt = time.Now().UTC()
}

func (n *Node) IsActive() bool {
	return n.Status == StatusActive
}
