package testutil

import (
	"github.com/suproxy/backend/internal/infrastructure/config"
)

// TestConfig returns a configuration suitable for integration tests
func TestConfig() *config.Config {
	return &config.Config{
		Environment: "test",
		Server: config.ServerConfig{
			Address:         ":0", // Random available port
			ReadTimeout:     10,
			WriteTimeout:    10,
			ShutdownTimeout: 5,
		},
		Database: config.DatabaseConfig{
			Host:            GetEnv("TEST_DB_HOST", "localhost"),
			Port:            GetEnvInt("TEST_DB_PORT", 5432),
			User:            GetEnv("TEST_DB_USER", "suproxy_test"),
			Password:        GetEnv("TEST_DB_PASSWORD", "suproxy_test"),
			DBName:          GetEnv("TEST_DB_NAME", "suproxy_test"),
			SSLMode:         "disable",
			MaxOpenConns:    10,
			MaxIdleConns:    5,
			ConnMaxLifetime: 5,
			ConnMaxIdleTime: 10,
		},
		Log: config.LogConfig{
			Level:  "error", // Reduce noise in tests
			Format: "json",
		},
		JWT: config.JWTConfig{
			SecretKey:          "test-secret-key-for-integration-tests-only",
			AccessTokenExpiry:  15,
			RefreshTokenExpiry: 168,
			Issuer:             "suproxy-test",
		},
		Xray: config.XrayConfig{
			UseMock:    true, // Always use mock in tests
			BinaryPath: "/usr/local/bin/xray",
			ConfigDir:  "/tmp/xray-test/config",
			LogDir:     "/tmp/xray-test/logs",
			BackupDir:  "/tmp/xray-test/backups",
			InstallDir: "/tmp/xray-test/install",
		},
		Metrics: config.MetricsConfig{
			Enabled:            false, // Disable metrics in tests
			CollectionInterval: 30,
		},
	}
}

// MinimalConfig returns minimal configuration for unit tests
func MinimalConfig() *config.Config {
	cfg := TestConfig()
	cfg.Log.Level = "fatal" // Suppress all logs in unit tests
	return cfg
}
