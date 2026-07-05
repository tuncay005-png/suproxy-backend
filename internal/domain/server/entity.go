package server

import (
	"time"

	"github.com/google/uuid"
)

// Server represents a physical or virtual server that hosts VPN nodes
type Server struct {
	ID        uuid.UUID
	Name      string
	Country   string
	City      string
	Hostname  string
	Provider  string
	IPv4      string
	IPv6      string
	Status    Status
	IsPublic  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewServer(name, country, city, hostname, provider, ipv4 string) (*Server, error) {
	if name == "" {
		return nil, ErrInvalidServerName
	}
	if country == "" {
		return nil, ErrInvalidCountry
	}
	if hostname == "" {
		return nil, ErrInvalidHostname
	}
	if ipv4 == "" {
		return nil, ErrInvalidIPAddress
	}

	return &Server{
		ID:        uuid.New(),
		Name:      name,
		Country:   country,
		City:      city,
		Hostname:  hostname,
		Provider:  provider,
		IPv4:      ipv4,
		IPv6:      "",
		Status:    StatusOffline,
		IsPublic:  true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}

func (s *Server) Activate() error {
	if s.Status == StatusActive {
		return ErrServerAlreadyActive
	}
	s.Status = StatusActive
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *Server) Deactivate() error {
	if s.Status == StatusInactive {
		return ErrServerAlreadyInactive
	}
	s.Status = StatusInactive
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *Server) SetMaintenance() error {
	s.Status = StatusMaintenance
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *Server) IsOnline() bool {
	return s.Status == StatusActive
}

func (s *Server) IsAvailable() bool {
	return s.Status == StatusActive && s.IsPublic
}

func (s *Server) UpdateDetails(name, city, provider string) {
	if name != "" {
		s.Name = name
	}
	if city != "" {
		s.City = city
	}
	if provider != "" {
		s.Provider = provider
	}
	s.UpdatedAt = time.Now().UTC()
}

func (s *Server) UpdateIPv6(ipv6 string) {
	s.IPv6 = ipv6
	s.UpdatedAt = time.Now().UTC()
}

func (s *Server) MakePublic() {
	s.IsPublic = true
	s.UpdatedAt = time.Now().UTC()
}

func (s *Server) MakePrivate() {
	s.IsPublic = false
	s.UpdatedAt = time.Now().UTC()
}
