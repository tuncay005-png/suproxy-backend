package metrics

import (
	"context"
	"time"

	"github.com/suproxy/backend/internal/domain/user"
	"github.com/suproxy/backend/internal/domain/xray"
	"github.com/suproxy/backend/internal/infrastructure/database"
	"github.com/suproxy/backend/internal/infrastructure/logger"
)

// Collector periodically collects business metrics
type Collector struct {
	userRepo     user.Repository
	instanceRepo xray.XrayInstanceRepository
	inboundRepo  xray.InboundRepository
	clientRepo   xray.ClientRepository
	db           *database.Database
	logger       *logger.Logger
	interval     time.Duration
	stopChan     chan struct{}
}

// NewCollector creates a new metrics collector
func NewCollector(
	userRepo user.Repository,
	instanceRepo xray.XrayInstanceRepository,
	inboundRepo xray.InboundRepository,
	clientRepo xray.ClientRepository,
	db *database.Database,
	log *logger.Logger,
	interval time.Duration,
) *Collector {
	return &Collector{
		userRepo:     userRepo,
		instanceRepo: instanceRepo,
		inboundRepo:  inboundRepo,
		clientRepo:   clientRepo,
		db:           db,
		logger:       log,
		interval:     interval,
		stopChan:     make(chan struct{}),
	}
}

// Start starts the metrics collector
func (c *Collector) Start() {
	c.logger.Info("Starting metrics collector", "interval", c.interval)
	
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	
	// Collect immediately on start
	c.collect()
	
	for {
		select {
		case <-ticker.C:
			c.collect()
		case <-c.stopChan:
			c.logger.Info("Metrics collector stopped")
			return
		}
	}
}

// Stop stops the metrics collector
func (c *Collector) Stop() {
	close(c.stopChan)
}

func (c *Collector) collect() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Collect user metrics
	c.collectUserMetrics(ctx)
	
	// Collect Xray metrics
	c.collectXrayMetrics(ctx)
	
	// Collect database metrics
	c.collectDatabaseMetrics()
}

func (c *Collector) collectUserMetrics(ctx context.Context) {
	users, _, err := c.userRepo.ListWithFilters(ctx, user.UserFilters{
		Offset: 0,
		Limit:  100000, // Large enough to get all
	})
	if err != nil {
		c.logger.Error("Failed to collect user metrics", "error", err)
		return
	}
	
	activeCount := 0
	for _, u := range users {
		if u.Status == user.StatusActive {
			activeCount++
		}
	}
	
	SetActiveUsers(float64(activeCount))
}

func (c *Collector) collectXrayMetrics(ctx context.Context) {
	// Collect instance metrics
	instances, _, err := c.instanceRepo.ListWithFilters(ctx, xray.XrayInstanceFilters{
		Offset: 0,
		Limit:  10000,
	})
	if err != nil {
		c.logger.Error("Failed to collect Xray instance metrics", "error", err)
		return
	}
	
	runningCount := 0
	stoppedCount := 0
	for _, inst := range instances {
		if inst.IsRunning() {
			runningCount++
		} else {
			stoppedCount++
		}
	}
	
	SetXrayInstances("running", float64(runningCount))
	SetXrayInstances("stopped", float64(stoppedCount))
	
	// Collect client metrics
	totalClients := 0
	enabledClients := 0
	totalInbounds := 0
	enabledInbounds := 0
	
	for _, inst := range instances {
		inbounds, err := c.inboundRepo.FindByInstanceID(ctx, inst.ID)
		if err != nil {
			continue
		}
		
		for _, inbound := range inbounds {
			totalInbounds++
			if inbound.Enabled {
				enabledInbounds++
			}
			
			clients, err := c.clientRepo.FindByInboundID(ctx, inbound.ID)
			if err != nil {
				continue
			}
			
			for _, client := range clients {
				totalClients++
				if client.Enabled {
					enabledClients++
				}
			}
		}
	}
	
	SetXrayClients("true", float64(enabledClients))
	SetXrayClients("false", float64(totalClients-enabledClients))
	SetXrayInbounds("true", float64(enabledInbounds))
	SetXrayInbounds("false", float64(totalInbounds-enabledInbounds))
}

func (c *Collector) collectDatabaseMetrics() {
	sqlDB, err := c.db.DB.DB()
	if err != nil {
		c.logger.Error("Failed to get database connection", "error", err)
		return
	}
	
	stats := sqlDB.Stats()
	SetDatabaseConnections(stats.InUse, stats.Idle)
	
	// Health check
	healthy := c.db.Ping() == nil
	SetHealthCheckStatus("database", healthy)
}
