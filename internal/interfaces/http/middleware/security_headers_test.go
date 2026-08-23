package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSecurityHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		environment string
		path        string
		wantHSTS    bool
		wantCSP     string
	}{
		{
			name:        "development environment - no HSTS",
			environment: "development",
			path:        "/api/v1/test",
			wantHSTS:    false,
			wantCSP:     "default-src 'none'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'",
		},
		{
			name:        "production environment - with HSTS",
			environment: "production",
			path:        "/api/v1/test",
			wantHSTS:    true,
			wantCSP:     "default-src 'none'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'",
		},
		{
			name:        "health endpoint - minimal CSP",
			environment: "production",
			path:        "/health",
			wantHSTS:    true,
			wantCSP:     "default-src 'none'; frame-ancestors 'none'",
		},
		{
			name:        "metrics endpoint - minimal CSP",
			environment: "production",
			path:        "/metrics",
			wantHSTS:    true,
			wantCSP:     "default-src 'none'; frame-ancestors 'none'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment
			_ = os.Setenv("SUPROXY_ENVIRONMENT", tt.environment)
			defer _ = os.Unsetenv("SUPROXY_ENVIRONMENT")

			// Create test router
			router := gin.New()
			router.Use(SecurityHeaders())
			router.GET(tt.path, func(c *gin.Context) {
				c.String(http.StatusOK, "test")
			})

			// Make request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.path, nil)
			router.ServeHTTP(w, req)

			// Check X-Frame-Options
			assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"), "X-Frame-Options should be DENY")

			// Check X-Content-Type-Options
			assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"), "X-Content-Type-Options should be nosniff")

			// Check X-XSS-Protection
			assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"), "X-XSS-Protection should be set")

			// Check Referrer-Policy
			assert.Equal(t, "strict-origin-when-cross-origin", w.Header().Get("Referrer-Policy"), "Referrer-Policy should be set")

			// Check Permissions-Policy
			assert.Contains(t, w.Header().Get("Permissions-Policy"), "geolocation=()", "Permissions-Policy should disable geolocation")

			// Check HSTS based on environment
			hstsHeader := w.Header().Get("Strict-Transport-Security")
			if tt.wantHSTS {
				assert.NotEmpty(t, hstsHeader, "HSTS should be set in production")
				assert.Contains(t, hstsHeader, "max-age=31536000", "HSTS should have 1 year max-age")
				assert.Contains(t, hstsHeader, "includeSubDomains", "HSTS should include subdomains")
			} else {
				assert.Empty(t, hstsHeader, "HSTS should not be set in development")
			}

			// Check CSP
			cspHeader := w.Header().Get("Content-Security-Policy")
			assert.Equal(t, tt.wantCSP, cspHeader, "CSP should match expected policy")
		})
	}
}

func TestSecurityHeaders_AllHeadersPresent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	_ = os.Setenv("SUPROXY_ENVIRONMENT", "production")
	defer _ = os.Unsetenv("SUPROXY_ENVIRONMENT")

	router := gin.New()
	router.Use(SecurityHeaders())
	router.GET("/api/v1/test", func(c *gin.Context) {
		c.String(http.StatusOK, "test")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/test", nil)
	router.ServeHTTP(w, req)

	// Verify all security headers are present
	expectedHeaders := map[string]bool{
		"X-Frame-Options":           true,
		"X-Content-Type-Options":    true,
		"X-XSS-Protection":          true,
		"Referrer-Policy":           true,
		"Permissions-Policy":        true,
		"Strict-Transport-Security": true,
		"Content-Security-Policy":   true,
	}

	for header := range expectedHeaders {
		value := w.Header().Get(header)
		assert.NotEmpty(t, value, "Header %s should be present", header)
	}
}

func TestSecurityHeaders_DoesNotBreakCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Security headers should not interfere with CORS headers
	router := gin.New()
	router.Use(SecurityHeaders())
	router.Use(func(c *gin.Context) {
		// Simulate CORS middleware setting headers
		c.Writer.Header().Set("Access-Control-Allow-Origin", "https://example.com")
		c.Next()
	})
	router.GET("/api/v1/test", func(c *gin.Context) {
		c.String(http.StatusOK, "test")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/test", nil)
	req.Header.Set("Origin", "https://example.com")
	router.ServeHTTP(w, req)

	// CORS header should still be present
	assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"), "CORS should work with security headers")

	// Security headers should also be present
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"), "Security headers should be present")
}

func TestIsAPIEndpoint(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"/api/v1/users", true},
		{"/api/auth/login", true},
		{"/api", false}, // Too short
		{"/health", false},
		{"/metrics", false},
		{"/", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := isAPIEndpoint(tt.path)
			assert.Equal(t, tt.want, got, "isAPIEndpoint(%s) = %v, want %v", tt.path, got, tt.want)
		})
	}
}

func TestIsMonitoringEndpoint(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"/health", true},
		{"/ready", true},
		{"/metrics", true},
		{"/api/v1/health", false},
		{"/", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := isMonitoringEndpoint(tt.path)
			assert.Equal(t, tt.want, got, "isMonitoringEndpoint(%s) = %v, want %v", tt.path, got, tt.want)
		})
	}
}

func TestSecureTransportDetectionMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		headers   map[string]string
		wantHTTPS bool
	}{
		{
			name:      "X-Forwarded-Proto https",
			headers:   map[string]string{"X-Forwarded-Proto": "https"},
			wantHTTPS: true,
		},
		{
			name:      "X-Forwarded-SSL on",
			headers:   map[string]string{"X-Forwarded-SSL": "on"},
			wantHTTPS: true,
		},
		{
			name:      "no secure headers",
			headers:   map[string]string{},
			wantHTTPS: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(SecureTransportDetectionMiddleware())
			router.GET("/test", func(c *gin.Context) {
				isHTTPS, exists := c.Get("is_https")
				if tt.wantHTTPS {
					assert.True(t, exists, "is_https should be set")
					assert.True(t, isHTTPS.(bool), "is_https should be true")
				} else {
					if exists {
						assert.False(t, isHTTPS.(bool), "is_https should be false")
					}
				}
				c.String(http.StatusOK, "test")
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}
			router.ServeHTTP(w, req)
		})
	}
}
