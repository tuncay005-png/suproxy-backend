package xray

import (
	"time"

	"github.com/google/uuid"
)

// RealityConfig represents REALITY protocol configuration for an inbound
type RealityConfig struct {
	ID          uuid.UUID
	InboundID   uuid.UUID
	PrivateKey  string
	PublicKey   string
	ShortID     string
	ServerName  string
	Fingerprint string
	SpiderX     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewRealityConfig(inboundID uuid.UUID, privateKey, publicKey, shortID, serverName, fingerprint, spiderX string) (*RealityConfig, error) {
	if inboundID == uuid.Nil {
		return nil, ErrInvalidInboundID
	}
	if privateKey == "" || publicKey == "" {
		return nil, ErrInvalidKeys
	}
	if serverName == "" {
		return nil, ErrInvalidServerName
	}
	if fingerprint == "" {
		return nil, ErrInvalidFingerprint
	}

	return &RealityConfig{
		ID:          uuid.New(),
		InboundID:   inboundID,
		PrivateKey:  privateKey,
		PublicKey:   publicKey,
		ShortID:     shortID,
		ServerName:  serverName,
		Fingerprint: fingerprint,
		SpiderX:     spiderX,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}, nil
}

func (r *RealityConfig) RegenerateKeys(privateKey, publicKey string) error {
	if privateKey == "" || publicKey == "" {
		return ErrInvalidKeys
	}
	r.PrivateKey = privateKey
	r.PublicKey = publicKey
	r.UpdatedAt = time.Now().UTC()
	return nil
}

func (r *RealityConfig) ChangeFingerprint(fingerprint string) error {
	if fingerprint == "" {
		return ErrInvalidFingerprint
	}
	r.Fingerprint = fingerprint
	r.UpdatedAt = time.Now().UTC()
	return nil
}

func (r *RealityConfig) UpdateServerName(serverName string) error {
	if serverName == "" {
		return ErrInvalidServerName
	}
	r.ServerName = serverName
	r.UpdatedAt = time.Now().UTC()
	return nil
}

func (r *RealityConfig) UpdateShortID(shortID string) {
	r.ShortID = shortID
	r.UpdatedAt = time.Now().UTC()
}

func (r *RealityConfig) UpdateSpiderX(spiderX string) {
	r.SpiderX = spiderX
	r.UpdatedAt = time.Now().UTC()
}
