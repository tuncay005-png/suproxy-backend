package bootstrap

import (
"os"
"strings"
"testing"
)

func TestInitialize_ProductionJWTSecretValidation(t *testing.T) {
// Save original environment
origEnv := os.Getenv("SUPROXY_ENVIRONMENT")
origSecret := os.Getenv("SUPROXY_JWT_SECRET_KEY")
defer func() {
os.Setenv("SUPROXY_ENVIRONMENT", origEnv)
os.Setenv("SUPROXY_JWT_SECRET_KEY", origSecret)
}()

tests := []struct {
name        string
environment string
secretKey   string
wantErr     bool
errContains []string
}{
{
name:        "production with default secret fails",
environment: "production",
secretKey:   "change-me-in-production",
wantErr:     true,
errContains: []string{"SECURITY ERROR", "SUPROXY_JWT_SECRET_KEY", "openssl rand -base64 32"},
},
{
name:        "production with empty secret fails",
environment: "production",
secretKey:   "",
wantErr:     true,
errContains: []string{"SECURITY ERROR", "SUPROXY_JWT_SECRET_KEY", "openssl rand -base64 32"},
},
{
name:        "development with default secret succeeds (but would log warning)",
environment: "development",
secretKey:   "change-me-in-production",
wantErr:     false, // Should not error, just warn in logs
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
os.Setenv("SUPROXY_ENVIRONMENT", tt.environment)
os.Setenv("SUPROXY_JWT_SECRET_KEY", tt.secretKey)

// Attempt to initialize
app, err := Initialize()

if tt.wantErr {
if err == nil {
t.Errorf("Initialize() expected error, got nil")
if app != nil {
app.Shutdown()
}
return
}

// Check error message contains expected strings
errMsg := err.Error()
for _, expected := range tt.errContains {
if !strings.Contains(errMsg, expected) {
t.Errorf("Initialize() error = %v, should contain %q", err, expected)
}
}
} else {
if err != nil {
// Development mode failures might be due to missing database, etc.
// Skip the full initialization check
if !strings.Contains(err.Error(), "SECURITY ERROR") {
t.Logf("Initialize() failed with non-security error (expected in test env): %v", err)
} else {
t.Errorf("Initialize() unexpected security error = %v", err)
}
}
if app != nil {
app.Shutdown()
}
}
})
}
}
