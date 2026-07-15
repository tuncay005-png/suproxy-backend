package system

import (
	"context"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/audit"
	"github.com/suproxy/backend/internal/domain/user"
	"github.com/suproxy/backend/internal/domain/xray"
	"github.com/suproxy/backend/internal/infrastructure/database"
	xrayRuntime "github.com/suproxy/backend/internal/infrastructure/xray/runtime"
)

var (
	// Application metadata (set at build time)
	Version     = "dev"
	BuildDate   = "unknown"
	GitCommit   = ""
	Environment = "development"
	StartTime   = time.Now()
)

// GetSystemHealthQuery handles system-wide health check
type GetSystemHealthQuery struct {
	db             *database.Database
	instanceRepo   xray.XrayInstanceRepository
	processManager xrayRuntime.Manager
}

func NewGetSystemHealthQuery(
	db *database.Database,
	instanceRepo xray.XrayInstanceRepository,
	processManager xrayRuntime.Manager,
) *GetSystemHealthQuery {
	return &GetSystemHealthQuery{
		db:             db,
		instanceRepo:   instanceRepo,
		processManager: processManager,
	}
}

type SystemHealth struct {
	Status     string
	Components map[string]ComponentHealth
	Uptime     int64
}

type ComponentHealth struct {
	Status  string
	Message string
	Latency int64
}

func (q *GetSystemHealthQuery) Execute(ctx context.Context) (*SystemHealth, error) {
	components := make(map[string]ComponentHealth)

	// Check database
	dbStart := time.Now()
	dbStatus := "healthy"
	dbMessage := "connected"
	if err := q.db.Ping(); err != nil {
		dbStatus = "unhealthy"
		dbMessage = err.Error()
	}
	components["database"] = ComponentHealth{
		Status:  dbStatus,
		Message: dbMessage,
		Latency: time.Since(dbStart).Milliseconds(),
	}

	// Check Xray instances
	xrayStart := time.Now()
	instances, _, err := q.instanceRepo.ListWithFilters(ctx, xray.XrayInstanceFilters{
		Offset: 0,
		Limit:  1000, // Get all instances for health check
	})
	xrayStatus := "healthy"
	xrayMessage := "all instances operational"
	if err != nil {
		xrayStatus = "unhealthy"
		xrayMessage = "failed to fetch instances"
	} else {
		runningCount := 0
		for _, inst := range instances {
			if inst.IsRunning() {
				runningCount++
			}
		}
		if runningCount == 0 && len(instances) > 0 {
			xrayStatus = "unhealthy"
			xrayMessage = "no running instances"
		} else if runningCount < len(instances) {
			xrayStatus = "degraded"
			xrayMessage = "some instances not running"
		}
	}
	components["xray"] = ComponentHealth{
		Status:  xrayStatus,
		Message: xrayMessage,
		Latency: time.Since(xrayStart).Milliseconds(),
	}

	// Overall status
	overallStatus := "ok"
	if dbStatus == "unhealthy" {
		overallStatus = "unhealthy"
	} else if xrayStatus == "unhealthy" {
		overallStatus = "degraded"
	} else if xrayStatus == "degraded" {
		overallStatus = "degraded"
	}

	return &SystemHealth{
		Status:     overallStatus,
		Components: components,
		Uptime:     int64(time.Since(StartTime).Seconds()),
	}, nil
}

// GetSystemStatsQuery handles system-wide statistics
type GetSystemStatsQuery struct {
	userRepo     user.Repository
	instanceRepo xray.XrayInstanceRepository
	inboundRepo  xray.InboundRepository
	clientRepo   xray.ClientRepository
	auditRepo    audit.Repository
}

func NewGetSystemStatsQuery(
	userRepo user.Repository,
	instanceRepo xray.XrayInstanceRepository,
	inboundRepo xray.InboundRepository,
	clientRepo xray.ClientRepository,
	auditRepo audit.Repository,
) *GetSystemStatsQuery {
	return &GetSystemStatsQuery{
		userRepo:     userRepo,
		instanceRepo: instanceRepo,
		inboundRepo:  inboundRepo,
		clientRepo:   clientRepo,
		auditRepo:    auditRepo,
	}
}

type SystemStats struct {
	Users UserStats
	Xray  XrayStats
	Audit AuditStatsData
}

type UserStats struct {
	TotalUsers       int64
	ActiveUsers      int64
	SuspendedUsers   int64
	AdminUsers       int64
	NewUsersToday    int64
	NewUsersThisWeek int64
}

type XrayStats struct {
	TotalInstances   int64
	RunningInstances int64
	StoppedInstances int64
	TotalInbounds    int64
	EnabledInbounds  int64
	TotalClients     int64
	EnabledClients   int64
}

type AuditStatsData struct {
	TotalLogs        int64
	LogsToday        int64
	LogsThisWeek     int64
	UniqueUsersToday int64
}

func (q *GetSystemStatsQuery) Execute(ctx context.Context) (*SystemStats, error) {
	stats := &SystemStats{}

	// User statistics
	allUsers, _, err := q.userRepo.ListWithFilters(ctx, user.UserFilters{
		Offset: 0,
		Limit:  100000, // Large enough to get all users
	})
	if err == nil {
		stats.Users.TotalUsers = int64(len(allUsers))

		today := time.Now().Truncate(24 * time.Hour)
		weekAgo := today.AddDate(0, 0, -7)

		for _, u := range allUsers {
			if u.Status == user.StatusActive {
				stats.Users.ActiveUsers++
			} else if u.Status == user.StatusSuspended {
				stats.Users.SuspendedUsers++
			}

			if u.Role == user.RoleAdmin {
				stats.Users.AdminUsers++
			}

			if u.CreatedAt.After(today) {
				stats.Users.NewUsersToday++
			}
			if u.CreatedAt.After(weekAgo) {
				stats.Users.NewUsersThisWeek++
			}
		}
	}

	// Xray statistics
	instances, _, err := q.instanceRepo.ListWithFilters(ctx, xray.XrayInstanceFilters{
		Offset: 0,
		Limit:  10000,
	})
	if err == nil {
		stats.Xray.TotalInstances = int64(len(instances))

		for _, inst := range instances {
			if inst.IsRunning() {
				stats.Xray.RunningInstances++
			} else {
				stats.Xray.StoppedInstances++
			}
		}
	}

	// Count inbounds and clients across all instances
	for _, inst := range instances {
		inbounds, err := q.inboundRepo.FindByInstanceID(ctx, inst.ID)
		if err == nil {
			stats.Xray.TotalInbounds += int64(len(inbounds))

			for _, inbound := range inbounds {
				if inbound.Enabled {
					stats.Xray.EnabledInbounds++
				}

				// Count clients for this inbound
				clients, err := q.clientRepo.FindByInboundID(ctx, inbound.ID)
				if err == nil {
					stats.Xray.TotalClients += int64(len(clients))

					for _, client := range clients {
						if client.Enabled {
							stats.Xray.EnabledClients++
						}
					}
				}
			}
		}
	}

	// Audit statistics
	totalLogs, err := q.auditRepo.Count(ctx)
	if err == nil {
		stats.Audit.TotalLogs = totalLogs
	}

	// Logs today and this week
	today := time.Now().Truncate(24 * time.Hour)
	weekAgo := today.AddDate(0, 0, -7)

	filters := audit.AuditFilters{
		DateFrom: &today,
		Limit:    10000, // Large enough to count
	}
	todayLogs, total, err := q.auditRepo.ListWithFilters(ctx, filters)
	if err == nil {
		stats.Audit.LogsToday = total

		// Count unique users today
		uniqueUsers := make(map[uuid.UUID]bool)
		for _, log := range todayLogs {
			uniqueUsers[log.UserID] = true
		}
		stats.Audit.UniqueUsersToday = int64(len(uniqueUsers))
	}

	weekFilters := audit.AuditFilters{
		DateFrom: &weekAgo,
		Limit:    100000,
	}
	_, weekTotal, err := q.auditRepo.ListWithFilters(ctx, weekFilters)
	if err == nil {
		stats.Audit.LogsThisWeek = weekTotal
	}

	return stats, nil
}

// GetVersionQuery handles version information
type GetVersionQuery struct{}

func NewGetVersionQuery() *GetVersionQuery {
	return &GetVersionQuery{}
}

type VersionInfo struct {
	Version     string
	BuildDate   string
	GitCommit   string
	GoVersion   string
	Environment string
	Uptime      int64
	StartTime   time.Time
}

func (q *GetVersionQuery) Execute(ctx context.Context) *VersionInfo {
	return &VersionInfo{
		Version:     Version,
		BuildDate:   BuildDate,
		GitCommit:   GitCommit,
		GoVersion:   runtime.Version(),
		Environment: Environment,
		Uptime:      int64(time.Since(StartTime).Seconds()),
		StartTime:   StartTime,
	}
}

// GetDatabaseStatusQuery handles database status check
type GetDatabaseStatusQuery struct {
	db *database.Database
}

func NewGetDatabaseStatusQuery(db *database.Database) *GetDatabaseStatusQuery {
	return &GetDatabaseStatusQuery{
		db: db,
	}
}

type DatabaseStatus struct {
	Status       string
	Type         string
	Latency      int64
	MaxOpenConns int
	OpenConns    int
	InUse        int
	Idle         int
	WaitCount    int64
	WaitDuration int64
}

func (q *GetDatabaseStatusQuery) Execute(ctx context.Context) (*DatabaseStatus, error) {
	start := time.Now()

	status := "connected"
	if err := q.db.Ping(); err != nil {
		status = "disconnected"
	}

	latency := time.Since(start).Milliseconds()

	// Get connection pool stats
	sqlDB, err := q.db.DB.DB()
	if err != nil {
		return nil, err
	}

	dbStats := sqlDB.Stats()

	return &DatabaseStatus{
		Status:       status,
		Type:         "postgresql",
		Latency:      latency,
		MaxOpenConns: dbStats.MaxOpenConnections,
		OpenConns:    dbStats.OpenConnections,
		InUse:        dbStats.InUse,
		Idle:         dbStats.Idle,
		WaitCount:    dbStats.WaitCount,
		WaitDuration: dbStats.WaitDuration.Milliseconds(),
	}, nil
}

// GetXraySystemStatusQuery handles Xray system status
type GetXraySystemStatusQuery struct {
	instanceRepo xray.XrayInstanceRepository
}

func NewGetXraySystemStatusQuery(instanceRepo xray.XrayInstanceRepository) *GetXraySystemStatusQuery {
	return &GetXraySystemStatusQuery{
		instanceRepo: instanceRepo,
	}
}

type XraySystemStatus struct {
	Status           string
	TotalInstances   int64
	RunningInstances int64
	HealthyInstances int64
	Instances        []XrayInstanceStatus
}

type XrayInstanceStatus struct {
	ID      uuid.UUID
	NodeID  uuid.UUID
	Status  string
	Healthy bool
	PID     int
	Uptime  int64
}

func (q *GetXraySystemStatusQuery) Execute(ctx context.Context) (*XraySystemStatus, error) {
	instances, _, err := q.instanceRepo.ListWithFilters(ctx, xray.XrayInstanceFilters{
		Offset: 0,
		Limit:  10000,
	})
	if err != nil {
		return nil, err
	}

	status := &XraySystemStatus{
		TotalInstances: int64(len(instances)),
		Instances:      make([]XrayInstanceStatus, 0, len(instances)),
	}

	for _, inst := range instances {
		instanceStatus := XrayInstanceStatus{
			ID:      inst.ID,
			NodeID:  inst.NodeID,
			Status:  string(inst.Status),
			Healthy: inst.Status == xray.StatusRunning, // Running means healthy
			PID:     0,                                 // PID not tracked in domain
		}

		if inst.IsRunning() {
			status.RunningInstances++
			instanceStatus.Uptime = int64(inst.GetUptime().Seconds())
		}

		if inst.Status == xray.StatusRunning {
			status.HealthyInstances++
		}

		status.Instances = append(status.Instances, instanceStatus)
	}

	// Determine overall status
	if status.TotalInstances == 0 {
		status.Status = "ok" // No instances configured yet
	} else if status.RunningInstances == 0 {
		status.Status = "down"
	} else if status.RunningInstances < status.TotalInstances {
		status.Status = "partial"
	} else {
		status.Status = "ok"
	}

	return status, nil
}
