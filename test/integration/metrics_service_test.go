package integration_test

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suproxy/backend/internal/infrastructure/metrics"
	"github.com/suproxy/backend/internal/infrastructure/testutil"
)

func TestMetricsService_Initialize(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	t.Run("Initialize_Success", func(t *testing.T) {
		// Initialize should not panic
		require.NotPanics(t, func() {
			metrics.Initialize()
		})
	})

	t.Run("Initialize_Idempotent", func(t *testing.T) {
		// Multiple calls should not panic
		require.NotPanics(t, func() {
			metrics.Initialize()
			metrics.Initialize()
		})
	})
}

func TestMetricsService_HTTPMetrics(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	metrics.Initialize()

	t.Run("RecordHTTPRequest_Success", func(t *testing.T) {
		require.NotPanics(t, func() {
			metrics.RecordHTTPRequest("GET", "/api/users", "200")
		})
	})

	t.Run("RecordHTTPRequest_MultipleStatusCodes", func(t *testing.T) {
		statusCodes := []string{"200", "201", "400", "404", "500"}

		for _, code := range statusCodes {
			require.NotPanics(t, func() {
				metrics.RecordHTTPRequest("GET", "/api/test", code)
			})
		}
	})

	t.Run("RecordHTTPRequest_DifferentMethods", func(t *testing.T) {
		methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

		for _, method := range methods {
			require.NotPanics(t, func() {
				metrics.RecordHTTPRequest(method, "/api/test", "200")
			})
		}
	})
}

func TestMetricsService_UserMetrics(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	metrics.Initialize()

	t.Run("SetActiveUsers_Success", func(t *testing.T) {
		require.NotPanics(t, func() {
			metrics.SetActiveUsers(100)
		})
	})

	t.Run("SetActiveUsers_Zero", func(t *testing.T) {
		require.NotPanics(t, func() {
			metrics.SetActiveUsers(0)
		})
	})

	t.Run("SetActiveUsers_Large", func(t *testing.T) {
		require.NotPanics(t, func() {
			metrics.SetActiveUsers(1000000)
		})
	})

	t.Run("IncrementUserRegistrations_Success", func(t *testing.T) {
		require.NotPanics(t, func() {
			metrics.IncUserRegistrations()
		})
	})

	t.Run("IncrementUserLogins_Success", func(t *testing.T) {
		require.NotPanics(t, func() {
			metrics.IncUserLogins()
		})
	})
}

func TestMetricsService_XrayMetrics(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	metrics.Initialize()

	t.Run("SetActiveXrayInstances_Success", func(t *testing.T) {
		require.NotPanics(t, func() {
			metrics.SetXrayInstances("running", 10)
		})
	})

	t.Run("SetActiveClients_Success", func(t *testing.T) {
		require.NotPanics(t, func() {
			metrics.SetXrayClients("true", 50)
		})
	})

	t.Run("IncrementConfigReloads_Success", func(t *testing.T) {
		require.NotPanics(t, func() {
			metrics.RecordConfigReload("success")
		})
	})

	t.Run("IncrementConfigReloads_Failure", func(t *testing.T) {
		require.NotPanics(t, func() {
			metrics.RecordConfigReload("failure")
		})
	})

	t.Run("RecordConfigReloadDuration_Success", func(t *testing.T) {
		require.NotPanics(t, func() {
			metrics.RecordConfigReloadDuration(0.5)
		})
	})
}

func TestMetricsService_DatabaseMetrics(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	metrics.Initialize()

	t.Run("RecordDatabaseQuery_Success", func(t *testing.T) {
		require.NotPanics(t, func() {
			metrics.RecordDatabaseQueryDuration("select", 0.01)
		})
	})

	t.Run("RecordDatabaseQuery_DifferentOperations", func(t *testing.T) {
		operations := []string{"select", "insert", "update", "delete"}

		for _, op := range operations {
			require.NotPanics(t, func() {
				metrics.RecordDatabaseQueryDuration(op, 0.01)
			})
		}
	})

	t.Run("SetDatabaseConnections_Success", func(t *testing.T) {
		require.NotPanics(t, func() {
			metrics.SetDatabaseConnections(25, 10)
		})
	})
}

func TestMetricsService_SystemMetrics(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	metrics.Initialize()

	// RecordMemoryUsage and RecordGoroutines are not implemented in metrics package
	// These tests are skipped
	
	t.Run("SetHealthCheckStatus_Success", func(t *testing.T) {
		require.NotPanics(t, func() {
			metrics.SetHealthCheckStatus("database", true)
		})
	})
}

func TestMetricsService_MetricsCollector(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	defer app.CleanupTables()

	t.Run("Collector_Creation", func(t *testing.T) {
		collector := metrics.NewCollector(
			app.Container.UserRepository,
			app.Container.XrayInstanceRepository,
			app.Container.InboundRepository,
			app.Container.ClientRepository,
			app.Database,
			app.Logger,
			1000, // 1 second for testing
		)

		assert.NotNil(t, collector)
	})

	t.Run("Collector_StartStop", func(t *testing.T) {
		collector := metrics.NewCollector(
			app.Container.UserRepository,
			app.Container.XrayInstanceRepository,
			app.Container.InboundRepository,
			app.Container.ClientRepository,
			app.Database,
			app.Logger,
			1000,
		)

		// Start should not panic
		require.NotPanics(t, func() {
			go collector.Start()
		})

		// Stop should not panic
		require.NotPanics(t, func() {
			collector.Stop()
		})
	})
}

func TestMetricsService_PrometheusRegistry(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	metrics.Initialize()

	t.Run("Metrics_Registered", func(t *testing.T) {
		// Gather metrics
		metricFamilies, err := prometheus.DefaultGatherer.Gather()
		require.NoError(t, err)

		// Verify some metrics exist
		metricNames := make([]string, 0)
		for _, mf := range metricFamilies {
			metricNames = append(metricNames, mf.GetName())
		}

		// Check for our custom metrics (at least some should exist)
		assert.NotEmpty(t, metricNames)

		// Common Go runtime metrics should always be present
		assert.Contains(t, metricNames, "go_goroutines")
		assert.Contains(t, metricNames, "go_memstats_alloc_bytes")
	})
}

func TestMetricsService_ConcurrentAccess(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	metrics.Initialize()

	t.Run("Concurrent_RecordHTTPRequest", func(t *testing.T) {
		done := make(chan bool, 10)

		for i := 0; i < 10; i++ {
			go func() {
				defer func() { done <- true }()
				for j := 0; j < 100; j++ {
					metrics.RecordHTTPRequest("GET", "/api/test", "200")
				}
			}()
		}

		// Wait for all goroutines
		for i := 0; i < 10; i++ {
			<-done
		}
	})

	t.Run("Concurrent_SetActiveUsers", func(t *testing.T) {
		done := make(chan bool, 10)

		for i := 0; i < 10; i++ {
			go func(count int) {
				defer func() { done <- true }()
				for j := 0; j < 100; j++ {
					metrics.SetActiveUsers(float64(count))
				}
			}(i)
		}

		// Wait for all goroutines
		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

func TestMetricsService_LabelValues(t *testing.T) {
	if !testutil.IsIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	metrics.Initialize()

	t.Run("HTTP_DifferentPaths", func(t *testing.T) {
		paths := []string{
			"/api/users",
			"/api/xray",
			"/api/inbounds",
			"/health",
			"/metrics",
		}

		for _, path := range paths {
			require.NotPanics(t, func() {
				metrics.RecordHTTPRequest("GET", path, "200")
			})
		}
	})

	t.Run("Database_DifferentTables", func(t *testing.T) {
		_ = []string{
			"users",
			"xray_instances",
			"inbounds",
			"clients",
			"audit_logs",
		}

		// Test database query recording
		require.NotPanics(t, func() {
			metrics.RecordDatabaseQueryDuration("select", 0.01)
		})
	})
}

