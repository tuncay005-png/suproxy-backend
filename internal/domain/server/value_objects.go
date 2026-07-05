package server

// Status represents the server status
type Status string

const (
	StatusActive      Status = "active"
	StatusInactive    Status = "inactive"
	StatusOffline     Status = "offline"
	StatusMaintenance Status = "maintenance"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusInactive, StatusOffline, StatusMaintenance:
		return true
	default:
		return false
	}
}

func (s Status) String() string {
	return string(s)
}
