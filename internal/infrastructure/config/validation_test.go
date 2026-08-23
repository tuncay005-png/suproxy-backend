package config

import (
	"strings"
	"testing"
)

func TestValidateJWTSecret(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		secretKey   string
		wantErr     bool
		errContains string
	}{
		{
			name:        "production with default secret fails",
			environment: "production",
			secretKey:   "change-me-in-production",
			wantErr:     true,
			errContains: "SECURITY ERROR",
		},
		{
			name:        "production with empty secret fails",
			environment: "production",
			secretKey:   "",
			wantErr:     true,
			errContains: "SECURITY ERROR",
		},
		{
			name:        "production with whitespace-only secret fails",
			environment: "production",
			secretKey:   "   ",
			wantErr:     true,
			errContains: "SECURITY ERROR",
		},
		{
			name:        "production with valid secret succeeds",
			environment: "production",
			secretKey:   "aB3dF6hJ9kL2mN5pQ8rT1uV4wX7yZ0+c=",
			wantErr:     false,
		},
		{
			name:        "development with default secret succeeds",
			environment: "development",
			secretKey:   "change-me-in-production",
			wantErr:     false,
		},
		{
			name:        "development with valid secret succeeds",
			environment: "development",
			secretKey:   "aB3dF6hJ9kL2mN5pQ8rT1uV4wX7yZ0+c=",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Environment: tt.environment,
				JWT: JWTConfig{
					SecretKey: tt.secretKey,
				},
			}

			err := ValidateJWTSecret(cfg)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateJWTSecret() expected error, got nil")
					return
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ValidateJWTSecret() error = %v, should contain %q", err, tt.errContains)
				}
				if !strings.Contains(err.Error(), "SUPROXY_JWT_SECRET_KEY") {
					t.Errorf("ValidateJWTSecret() error should contain 'SUPROXY_JWT_SECRET_KEY', got %v", err)
				}
				if !strings.Contains(err.Error(), "openssl rand -base64 32") {
					t.Errorf("ValidateJWTSecret() error should contain 'openssl rand -base64 32', got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateJWTSecret() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestIsDefaultJWTSecret(t *testing.T) {
	tests := []struct {
		name      string
		secretKey string
		want      bool
	}{
		{
			name:      "default secret returns true",
			secretKey: "change-me-in-production",
			want:      true,
		},
		{
			name:      "empty secret returns true",
			secretKey: "",
			want:      true,
		},
		{
			name:      "whitespace-only secret returns true",
			secretKey: "   ",
			want:      true,
		},
		{
			name:      "valid secret returns false",
			secretKey: "aB3dF6hJ9kL2mN5pQ8rT1uV4wX7yZ0+c=",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				JWT: JWTConfig{
					SecretKey: tt.secretKey,
				},
			}

			got := IsDefaultJWTSecret(cfg)
			if got != tt.want {
				t.Errorf("IsDefaultJWTSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}
