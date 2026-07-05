package xray

import (
	"time"

	"github.com/google/uuid"
)

// Client represents an Xray client configuration (user)
type Client struct {
	ID        uuid.UUID
	InboundID uuid.UUID
	UserID    uuid.UUID
	UUID      string
	Flow      string
	Email     string
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewClient(inboundID, userID uuid.UUID, clientUUID, flow, email string) (*Client, error) {
	if inboundID == uuid.Nil {
		return nil, ErrInvalidInboundID
	}
	if userID == uuid.Nil {
		return nil, ErrInvalidUserID
	}
	if clientUUID == "" {
		return nil, ErrInvalidUUID
	}
	if email == "" {
		return nil, ErrInvalidEmail
	}

	return &Client{
		ID:        uuid.New(),
		InboundID: inboundID,
		UserID:    userID,
		UUID:      clientUUID,
		Flow:      flow,
		Email:     email,
		Enabled:   true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}

func (c *Client) Enable() error {
	if c.Enabled {
		return ErrClientAlreadyEnabled
	}
	c.Enabled = true
	c.UpdatedAt = time.Now().UTC()
	return nil
}

func (c *Client) Disable() error {
	if !c.Enabled {
		return ErrClientAlreadyDisabled
	}
	c.Enabled = false
	c.UpdatedAt = time.Now().UTC()
	return nil
}

func (c *Client) RegenerateUUID(newUUID string) error {
	if newUUID == "" {
		return ErrInvalidUUID
	}
	c.UUID = newUUID
	c.UpdatedAt = time.Now().UTC()
	return nil
}

func (c *Client) IsEnabled() bool {
	return c.Enabled
}

func (c *Client) UpdateFlow(flow string) {
	c.Flow = flow
	c.UpdatedAt = time.Now().UTC()
}

func (c *Client) UpdateEmail(email string) error {
	if email == "" {
		return ErrInvalidEmail
	}
	c.Email = email
	c.UpdatedAt = time.Now().UTC()
	return nil
}
