package testutil

import (
	"os"
	"strconv"
)

// GetEnv retrieves environment variable or returns default value
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt retrieves integer environment variable or returns default value
func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetEnvBool retrieves boolean environment variable or returns default value
func GetEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// IsCI returns true if running in CI environment
func IsCI() bool {
	return GetEnvBool("CI", false)
}

// IsIntegrationTest returns true if integration tests should run
func IsIntegrationTest() bool {
	return GetEnvBool("INTEGRATION_TEST", false)
}
