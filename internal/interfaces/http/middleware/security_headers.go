package middleware

import (
"os"

"github.com/gin-gonic/gin"
)

// SecurityHeaders adds comprehensive security headers to all responses
func SecurityHeaders() gin.HandlerFunc {
// Determine environment
environment := os.Getenv("SUPROXY_ENVIRONMENT")
if environment == "" {
environment = "development"
}

// Check if we're behind HTTPS (for HSTS)
// In production, this should be true when behind reverse proxy with HTTPS
isHTTPS := environment == "production"

return func(c *gin.Context) {
// X-Frame-Options: Prevents clickjacking attacks
// DENY: The page cannot be displayed in a frame
c.Writer.Header().Set("X-Frame-Options", "DENY")

// X-Content-Type-Options: Prevents MIME-type sniffing
// nosniff: Browser should not try to guess content type
c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

// X-XSS-Protection: Legacy XSS protection (older browsers)
// 1; mode=block: Enable XSS filter and block page if attack detected
c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")

// Referrer-Policy: Controls how much referrer information is sent
// strict-origin-when-cross-origin: Send full URL for same-origin,
// only origin for HTTPS cross-origin, nothing for HTTP destinations
c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

// Permissions-Policy: Control browser features and APIs
// Disable unnecessary features like geolocation, camera, microphone
c.Writer.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=(), payment=(), usb=(), magnetometer=(), gyroscope=(), accelerometer=()")

// HTTP Strict Transport Security (HSTS)
// Only set in production when HTTPS is being used
// Forces browsers to use HTTPS for all future requests
if isHTTPS {
// max-age=31536000: 1 year
// includeSubDomains: Apply to all subdomains
c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
}

// Content-Security-Policy (CSP)
// This is carefully crafted to work with the frontend
// Adjust if frontend requirements change
csp := buildContentSecurityPolicy(c)
if csp != "" {
c.Writer.Header().Set("Content-Security-Policy", csp)
}

c.Next()
}
}

// buildContentSecurityPolicy creates a CSP header appropriate for the endpoint
func buildContentSecurityPolicy(c *gin.Context) string {
path := c.Request.URL.Path

// For API endpoints, use a strict CSP
// This allows the API to work but prevents any scripts from running
if isAPIEndpoint(path) {
return "default-src 'none'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'"
}

// For health/metrics endpoints, minimal CSP
if isMonitoringEndpoint(path) {
return "default-src 'none'; frame-ancestors 'none'"
}

// Default CSP for other endpoints
return "default-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'"
}

// isAPIEndpoint checks if the path is an API endpoint
func isAPIEndpoint(path string) bool {
return len(path) > 4 && path[:5] == "/api/"
}

// isMonitoringEndpoint checks if the path is a monitoring endpoint
func isMonitoringEndpoint(path string) bool {
return path == "/health" || path == "/ready" || path == "/metrics"
}

// SecureTransportDetectionMiddleware detects if request came through HTTPS
// This should be used when behind a reverse proxy
func SecureTransportDetectionMiddleware() gin.HandlerFunc {
return func(c *gin.Context) {
// Check various headers that reverse proxies set
proto := c.Request.Header.Get("X-Forwarded-Proto")
if proto == "https" {
c.Set("is_https", true)
}

// Alternative headers
if c.Request.Header.Get("X-Forwarded-SSL") == "on" {
c.Set("is_https", true)
}

if c.Request.TLS != nil {
c.Set("is_https", true)
}

c.Next()
}
}
