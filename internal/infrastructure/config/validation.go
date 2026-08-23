package config

import (
	"fmt"
	"strings"
)

// ValidateJWTSecret validates the JWT secret key based on environment.
// In production, it rejects insecure default values.
// In development, it allows defaults but logs a warning.
func ValidateJWTSecret(cfg *Config) error {
	secret := strings.TrimSpace(cfg.JWT.SecretKey)

	// Check for invalid secrets
	isInsecure := secret == "" ||
		secret == "change-me-in-production" ||
		strings.TrimSpace(secret) == ""

	if !isInsecure {
		return nil
	}

	// Production: reject insecure secrets
	if cfg.Environment == "production" {
		return fmt.Errorf(
			"SECURITY ERROR: SUPROXY_JWT_SECRET_KEY must be set to a secure value in production.\n"+
				"Current value: '%s'\n"+
				"Generate a secure secret using: openssl rand -base64 32",
			cfg.JWT.SecretKey,
		)
	}

	// Development: allow but don't log warning here (will be logged by bootstrap)
	return nil
}

// IsDefaultJWTSecret returns true if the JWT secret is using the default value.
func IsDefaultJWTSecret(cfg *Config) bool {
	secret := strings.TrimSpace(cfg.JWT.SecretKey)
	return secret == "" || secret == "change-me-in-production" || strings.TrimSpace(secret) == ""
}
