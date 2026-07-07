package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	once sync.Once
	
	// HTTP Metrics
	httpRequestsTotal *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
	httpRequestsInFlight prometheus.Gauge
	httpErrorsTotal *prometheus.CounterVec
	
	// Business Metrics - Users
	activeUsersTotal prometheus.Gauge
	userRegistrationsTotal prometheus.Counter
	userLoginsTotal prometheus.Counter
	userLoginFailuresTotal prometheus.Counter
	
	// Business Metrics - Xray
	xrayInstancesTotal *prometheus.GaugeVec
	xrayClientsTotal *prometheus.GaugeVec
	xrayInboundsTotal *prometheus.GaugeVec
	
	// Business Metrics - Provisioning
	provisioningOperationsTotal *prometheus.CounterVec
	provisioningDuration *prometheus.HistogramVec
	provisioningErrorsTotal *prometheus.CounterVec
	configReloadTotal *prometheus.CounterVec
	configReloadDuration prometheus.Histogram
	
	// Database Metrics
	databaseConnectionsInUse prometheus.Gauge
	databaseConnectionsIdle prometheus.Gauge
	databaseConnectionsWaitCount prometheus.Counter
	databaseQueryDuration *prometheus.HistogramVec
	
	// System Metrics
	healthCheckStatus *prometheus.GaugeVec
	auditLogsTotal prometheus.Counter
)

// Initialize initializes all Prometheus metrics
// This should be called once at application startup
func Initialize() {
	once.Do(func() {
		// HTTP Metrics
		httpRequestsTotal = promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "suproxy_http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		)
		
		httpRequestDuration = promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "suproxy_http_request_duration_seconds",
				Help: "HTTP request duration in seconds",
				Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			},
			[]string{"method", "path"},
		)
		
		httpRequestsInFlight = promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "suproxy_http_requests_in_flight",
				Help: "Current number of HTTP requests being processed",
			},
		)
		
		httpErrorsTotal = promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "suproxy_http_errors_total",
				Help: "Total number of HTTP errors",
			},
			[]string{"method", "path", "status"},
		)
		
		// Business Metrics - Users
		activeUsersTotal = promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "suproxy_active_users_total",
				Help: "Total number of active users",
			},
		)
		
		userRegistrationsTotal = promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "suproxy_user_registrations_total",
				Help: "Total number of user registrations",
			},
		)
		
		userLoginsTotal = promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "suproxy_user_logins_total",
				Help: "Total number of successful user logins",
			},
		)
		
		userLoginFailuresTotal = promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "suproxy_user_login_failures_total",
				Help: "Total number of failed user login attempts",
			},
		)
		
		// Business Metrics - Xray
		xrayInstancesTotal = promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "suproxy_xray_instances_total",
				Help: "Total number of Xray instances by status",
			},
			[]string{"status"},
		)
		
		xrayClientsTotal = promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "suproxy_xray_clients_total",
				Help: "Total number of Xray clients",
			},
			[]string{"enabled"},
		)
		
		xrayInboundsTotal = promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "suproxy_xray_inbounds_total",
				Help: "Total number of Xray inbounds",
			},
			[]string{"enabled"},
		)
		
		// Business Metrics - Provisioning
		provisioningOperationsTotal = promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "suproxy_provisioning_operations_total",
				Help: "Total number of provisioning operations",
			},
			[]string{"operation", "status"},
		)
		
		provisioningDuration = promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "suproxy_provisioning_duration_seconds",
				Help: "Provisioning operation duration in seconds",
				Buckets: []float64{.1, .25, .5, 1, 2.5, 5, 10, 30, 60},
			},
			[]string{"operation"},
		)
		
		provisioningErrorsTotal = promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "suproxy_provisioning_errors_total",
				Help: "Total number of provisioning errors",
			},
			[]string{"operation", "error_type"},
		)
		
		configReloadTotal = promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "suproxy_config_reload_total",
				Help: "Total number of configuration reloads",
			},
			[]string{"status"},
		)
		
		configReloadDuration = promauto.NewHistogram(
			prometheus.HistogramOpts{
				Name: "suproxy_config_reload_duration_seconds",
				Help: "Configuration reload duration in seconds",
				Buckets: []float64{.1, .25, .5, 1, 2.5, 5, 10, 30},
			},
		)
		
		// Database Metrics
		databaseConnectionsInUse = promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "suproxy_database_connections_in_use",
				Help: "Number of database connections currently in use",
			},
		)
		
		databaseConnectionsIdle = promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "suproxy_database_connections_idle",
				Help: "Number of idle database connections",
			},
		)
		
		databaseConnectionsWaitCount = promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "suproxy_database_connections_wait_count_total",
				Help: "Total number of times a connection had to wait",
			},
		)
		
		databaseQueryDuration = promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "suproxy_database_query_duration_seconds",
				Help: "Database query duration in seconds",
				Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
			},
			[]string{"operation"},
		)
		
		// System Metrics
		healthCheckStatus = promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "suproxy_health_check_status",
				Help: "Health check status (1 = healthy, 0 = unhealthy)",
			},
			[]string{"component"},
		)
		
		auditLogsTotal = promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "suproxy_audit_logs_total",
				Help: "Total number of audit log entries",
			},
		)
	})
}

// RecordHTTPRequest records an HTTP request metric
func RecordHTTPRequest(method, path, status string) {
	if httpRequestsTotal != nil {
		httpRequestsTotal.WithLabelValues(method, path, status).Inc()
	}
}

// RecordHTTPDuration records HTTP request duration
func RecordHTTPDuration(method, path string, duration float64) {
	if httpRequestDuration != nil {
		httpRequestDuration.WithLabelValues(method, path).Observe(duration)
	}
}

// IncHTTPRequestsInFlight increments in-flight requests
func IncHTTPRequestsInFlight() {
	if httpRequestsInFlight != nil {
		httpRequestsInFlight.Inc()
	}
}

// DecHTTPRequestsInFlight decrements in-flight requests
func DecHTTPRequestsInFlight() {
	if httpRequestsInFlight != nil {
		httpRequestsInFlight.Dec()
	}
}

// RecordHTTPError records an HTTP error
func RecordHTTPError(method, path, status string) {
	if httpErrorsTotal != nil {
		httpErrorsTotal.WithLabelValues(method, path, status).Inc()
	}
}

// SetActiveUsers sets the active users gauge
func SetActiveUsers(count float64) {
	if activeUsersTotal != nil {
		activeUsersTotal.Set(count)
	}
}

// IncUserRegistrations increments user registrations counter
func IncUserRegistrations() {
	if userRegistrationsTotal != nil {
		userRegistrationsTotal.Inc()
	}
}

// IncUserLogins increments successful login counter
func IncUserLogins() {
	if userLoginsTotal != nil {
		userLoginsTotal.Inc()
	}
}

// IncUserLoginFailures increments failed login counter
func IncUserLoginFailures() {
	if userLoginFailuresTotal != nil {
		userLoginFailuresTotal.Inc()
	}
}

// SetXrayInstances sets Xray instances gauge by status
func SetXrayInstances(status string, count float64) {
	if xrayInstancesTotal != nil {
		xrayInstancesTotal.WithLabelValues(status).Set(count)
	}
}

// SetXrayClients sets Xray clients gauge
func SetXrayClients(enabled string, count float64) {
	if xrayClientsTotal != nil {
		xrayClientsTotal.WithLabelValues(enabled).Set(count)
	}
}

// SetXrayInbounds sets Xray inbounds gauge
func SetXrayInbounds(enabled string, count float64) {
	if xrayInboundsTotal != nil {
		xrayInboundsTotal.WithLabelValues(enabled).Set(count)
	}
}

// RecordProvisioningOperation records a provisioning operation
func RecordProvisioningOperation(operation, status string) {
	if provisioningOperationsTotal != nil {
		provisioningOperationsTotal.WithLabelValues(operation, status).Inc()
	}
}

// RecordProvisioningDuration records provisioning duration
func RecordProvisioningDuration(operation string, duration float64) {
	if provisioningDuration != nil {
		provisioningDuration.WithLabelValues(operation).Observe(duration)
	}
}

// RecordProvisioningError records a provisioning error
func RecordProvisioningError(operation, errorType string) {
	if provisioningErrorsTotal != nil {
		provisioningErrorsTotal.WithLabelValues(operation, errorType).Inc()
	}
}

// RecordConfigReload records a config reload operation
func RecordConfigReload(status string) {
	if configReloadTotal != nil {
		configReloadTotal.WithLabelValues(status).Inc()
	}
}

// RecordConfigReloadDuration records config reload duration
func RecordConfigReloadDuration(duration float64) {
	if configReloadDuration != nil {
		configReloadDuration.Observe(duration)
	}
}

// SetDatabaseConnections sets database connection metrics
func SetDatabaseConnections(inUse, idle int) {
	if databaseConnectionsInUse != nil {
		databaseConnectionsInUse.Set(float64(inUse))
	}
	if databaseConnectionsIdle != nil {
		databaseConnectionsIdle.Set(float64(idle))
	}
}

// IncDatabaseWaitCount increments database wait count
func IncDatabaseWaitCount() {
	if databaseConnectionsWaitCount != nil {
		databaseConnectionsWaitCount.Inc()
	}
}

// RecordDatabaseQueryDuration records database query duration
func RecordDatabaseQueryDuration(operation string, duration float64) {
	if databaseQueryDuration != nil {
		databaseQueryDuration.WithLabelValues(operation).Observe(duration)
	}
}

// SetHealthCheckStatus sets health check status
func SetHealthCheckStatus(component string, healthy bool) {
	if healthCheckStatus != nil {
		status := 0.0
		if healthy {
			status = 1.0
		}
		healthCheckStatus.WithLabelValues(component).Set(status)
	}
}

// IncAuditLogs increments audit logs counter
func IncAuditLogs() {
	if auditLogsTotal != nil {
		auditLogsTotal.Inc()
	}
}
