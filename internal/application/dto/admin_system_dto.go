package dto

import "time"

// SystemHealthResponse represents overall system health
type SystemHealthResponse struct {
	Status       string                       `json:"status"`        // ok, degraded, unhealthy
	Timestamp    time.Time                    `json:"timestamp"`
	Components   map[string]ComponentHealth   `json:"components"`
	Version      string                       `json:"version"`
	Uptime       int64                        `json:"uptime_seconds"`
}

// ComponentHealth represents individual component health
type ComponentHealth struct {
	Status  string `json:"status"`  // healthy, unhealthy, unknown
	Message string `json:"message,omitempty"`
	Latency int64  `json:"latency_ms,omitempty"` // Response time in milliseconds
}

// SystemStatsResponse represents system-wide statistics
type SystemStatsResponse struct {
	Users     UserStats     `json:"users"`
	Xray      XrayStats     `json:"xray"`
	Audit     AuditStats    `json:"audit"`
	Timestamp time.Time     `json:"timestamp"`
}

// UserStats represents user statistics
type UserStats struct {
	TotalUsers       int64 `json:"total_users"`
	ActiveUsers      int64 `json:"active_users"`
	SuspendedUsers   int64 `json:"suspended_users"`
	AdminUsers       int64 `json:"admin_users"`
	NewUsersToday    int64 `json:"new_users_today"`
	NewUsersThisWeek int64 `json:"new_users_this_week"`
}

// XrayStats represents Xray system statistics
type XrayStats struct {
	TotalInstances   int64 `json:"total_instances"`
	RunningInstances int64 `json:"running_instances"`
	StoppedInstances int64 `json:"stopped_instances"`
	TotalInbounds    int64 `json:"total_inbounds"`
	EnabledInbounds  int64 `json:"enabled_inbounds"`
	TotalClients     int64 `json:"total_clients"`
	EnabledClients   int64 `json:"enabled_clients"`
}

// AuditStats represents audit system statistics
type AuditStats struct {
	TotalLogs        int64 `json:"total_logs"`
	LogsToday        int64 `json:"logs_today"`
	LogsThisWeek     int64 `json:"logs_this_week"`
	UniqueUsersToday int64 `json:"unique_users_today"`
}

// VersionResponse represents application version information
type VersionResponse struct {
	Version     string    `json:"version"`
	BuildDate   string    `json:"build_date"`
	GitCommit   string    `json:"git_commit,omitempty"`
	GoVersion   string    `json:"go_version"`
	Environment string    `json:"environment"`
	Uptime      int64     `json:"uptime_seconds"`
	StartTime   time.Time `json:"start_time"`
}

// DatabaseStatusResponse represents database connection status
type DatabaseStatusResponse struct {
	Status       string `json:"status"`        // connected, disconnected
	Type         string `json:"type"`          // postgresql
	Latency      int64  `json:"latency_ms"`
	MaxOpenConns int    `json:"max_open_conns"`
	OpenConns    int    `json:"open_conns"`
	InUse        int    `json:"in_use"`
	Idle         int    `json:"idle"`
	WaitCount    int64  `json:"wait_count"`
	WaitDuration int64  `json:"wait_duration_ms"`
}

// XraySystemStatusResponse represents Xray system status
type XraySystemStatusResponse struct {
	Status           string                   `json:"status"` // ok, partial, down
	TotalInstances   int64                    `json:"total_instances"`
	RunningInstances int64                    `json:"running_instances"`
	HealthyInstances int64                    `json:"healthy_instances"`
	Instances        []XrayInstanceStatusInfo `json:"instances,omitempty"`
}

// XrayInstanceStatusInfo represents brief instance status
type XrayInstanceStatusInfo struct {
	ID       string `json:"id"`
	NodeID   string `json:"node_id"`
	Status   string `json:"status"`   // running, stopped
	Healthy  bool   `json:"healthy"`
	PID      int    `json:"pid,omitempty"`
	Uptime   int64  `json:"uptime_seconds,omitempty"`
}
