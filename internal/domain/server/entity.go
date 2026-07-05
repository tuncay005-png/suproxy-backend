package server

import (
	"time"

	"github.com/google/uuid"
)

type Server struct {
	ID          uuid.UUID
	Name        string
	Location    Location
	IPAddress   string
	Domain      string
	Port        int
	Capacity    Capacity
	Status      Status
	Tags        []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewServer(name string, location Location, ipAddress string, domain string, port int, capacity Capacity) (*Server, error) {
	if name == "" {
		return nil, ErrInvalidServerName
	}
	if ipAddress == "" {
		return nil, ErrInvalidIPAddress
	}
	if port <= 0 || port > 65535 {
		return nil, ErrInvalidPort
	}

	return &Server{
		ID:        uuid.New(),
		Name:      name,
		Location:  location,
		IPAddress: ipAddress,
		Domain:    domain,
		Port:      port,
		Capacity:  capacity,
		Status:    StatusActive,
		Tags:      []string{},
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
	if s.Status == StatusMaintenance {
		return ErrServerAlreadyInMaintenance
	}
	s.Status = StatusMaintenance
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *Server) UpdateCapacity(capacity Capacity) {
	s.Capacity = capacity
	s.UpdatedAt = time.Now().UTC()
}

func (s *Server) AddTag(tag string) {
	for _, t := range s.Tags {
		if t == tag {
			return
		}
	}
	s.Tags = append(s.Tags, tag)
	s.UpdatedAt = time.Now().UTC()
}

func (s *Server) RemoveTag(tag string) {
	newTags := []string{}
	for _, t := range s.Tags {
		if t != tag {
			newTags = append(newTags, t)
		}
	}
	s.Tags = newTags
	s.UpdatedAt = time.Now().UTC()
}

func (s *Server) IsActive() bool {
	return s.Status == StatusActive
}

func (s *Server) IsAvailable() bool {
	return s.Status == StatusActive && !s.Capacity.IsFull()
}
