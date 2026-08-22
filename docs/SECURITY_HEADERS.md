# Security Headers Implementation Guide

**Last Updated:** 2024-08-21  
**Version:** 1.0.0

---

## Overview

The SuProxy backend implements comprehensive security headers to protect against common web vulnerabilities. These headers are automatically applied to all API responses.

---

## Implemented Security Headers

### 1. X-Frame-Options: DENY

**Purpose:** Prevents clickjacking attacks

**Value:** `DENY`

**What it does:**
- Prevents the page from being displayed in frames, iframes, or objects
- Protects against clickjacking attacks where malicious sites try to trick users into clicking hidden elements

**Browser Support:** All modern browsers

**Example:**
```http
X-Frame-Options: DENY
```

---

### 2. X-Content-Type-Options: nosniff

**Purpose:** Prevents MIME-type sniffing

**Value:** `nosniff`

**What it does:**
- Forces browsers to respect the Content-Type header
- Prevents browsers from trying to "sniff" the content type
- Blocks execution of scripts if Content-Type is incorrect

**Browser Support:** All modern browsers

**Example:**
```http
X-Content-Type-Options: nosniff
```

---

### 3. X-XSS-Protection: 1; mode=block

**Purpose:** Legacy XSS protection (for older browsers)

**Value:** `1; mode=block`

**What it does:**
- Enables built-in XSS filter in older browsers
- Blocks the page if XSS attack is detected
- Modern browsers use CSP instead, but this provides backwards compatibility

**Browser Support:** Legacy IE, Safari, Chrome (deprecated in newer versions)

**Example:**
```http
X-XSS-Protection: 1; mode=block
```

---

### 4. Referrer-Policy: strict-origin-when-cross-origin

**Purpose:** Controls referrer information

**Value:** `strict-origin-when-cross-origin`

**What it does:**
- **Same-origin requests:** Sends full URL
- **Cross-origin HTTPS → HTTPS:** Sends only origin
- **HTTPS → HTTP:** Sends nothing (prevents leaking secure URLs)
- Protects user privacy and prevents information leakage

**Browser Support:** All modern browsers

**Example:**
```http
Referrer-Policy: strict-origin-when-cross-origin
```

**Behavior:**
```
Request from: https://api.suproxy.com/api/v1/users/123
To same origin: https://api.suproxy.com/api/v1/profile
Referrer sent: https://api.suproxy.com/api/v1/users/123 (full URL)

Request from: https://api.suproxy.com/api/v1/users/123
To different origin: https://external.com/endpoint
Referrer sent: https://api.suproxy.com (origin only)

Request from: https://api.suproxy.com/api/v1/users/123
To HTTP site: http://external.com/endpoint
Referrer sent: (nothing - protects from leaking HTTPS URLs)
```

---

### 5. Permissions-Policy

**Purpose:** Controls browser features and APIs

**Value:** `geolocation=(), microphone=(), camera=(), payment=(), usb=(), magnetometer=(), gyroscope=(), accelerometer=()`

**What it does:**
- Explicitly disables unnecessary browser features
- Prevents malicious scripts from accessing device APIs
- Improves privacy and security

**Features Disabled:**
- Geolocation API
- Microphone access
- Camera access
- Payment API
- USB device access
- Motion sensors (magnetometer, gyroscope, accelerometer)

**Browser Support:** Modern browsers (Chrome, Edge, Firefox)

**Example:**
```http
Permissions-Policy: geolocation=(), microphone=(), camera=()
```

---

### 6. Strict-Transport-Security (HSTS)

**Purpose:** Forces HTTPS connections

**Value (Production only):** `max-age=31536000; includeSubDomains`

**What it does:**
- Forces browsers to use HTTPS for all future requests
- Prevents downgrade attacks
- Applied for 1 year (31536000 seconds)
- Includes all subdomains

**Environment Behavior:**
- **Production:** HSTS enabled
- **Development:** HSTS disabled (allows HTTP)

**Browser Support:** All modern browsers

**Example:**
```http
Strict-Transport-Security: max-age=31536000; includeSubDomains
```

**Important Notes:**
- Only set when HTTPS is actually used
- Once set, browser enforces HTTPS even if user types http://
- Subdomains must also support HTTPS

---

### 7. Content-Security-Policy (CSP)

**Purpose:** Prevents XSS and data injection attacks

**Values (path-dependent):**

#### For API Endpoints (`/api/*`)
```
default-src 'none'; 
frame-ancestors 'none'; 
base-uri 'self'; 
form-action 'self'
```

**What it does:**
- Blocks all content loading by default
- Prevents page from being framed
- Restricts base URL changes
- Only allows form submissions to same origin

#### For Health/Metrics Endpoints (`/health`, `/ready`, `/metrics`)
```
default-src 'none'; 
frame-ancestors 'none'
```

**What it does:**
- Minimal policy for monitoring endpoints
- Prevents framing and content loading

#### For Other Endpoints
```
default-src 'self'; 
frame-ancestors 'none'; 
base-uri 'self'; 
form-action 'self'
```

**What it does:**
- Allows content from same origin
- Prevents framing
- Restricts base and form actions

**Browser Support:** All modern browsers

---

## Implementation Details

### Middleware Order

Security headers are applied early in the middleware chain:

```go
r.engine.Use(middleware.RequestIDMiddleware())                    // 1. Request ID
r.engine.Use(middleware.MetricsMiddleware())                      // 2. Metrics
r.engine.Use(middleware.SecureTransportDetectionMiddleware())     // 3. HTTPS detection
r.engine.Use(middleware.SecurityHeaders())                        // 4. Security headers
r.engine.Use(middleware.CORS())                                   // 5. CORS (after security)
r.engine.Use(middleware.ErrorHandler(r.logger))                   // 6. Error handling
r.engine.Use(middleware.RequestLogger(r.logger))                  // 7. Request logging
```

**Why this order?**
1. Request ID first for tracing
2. Metrics to measure everything
3. Detect HTTPS before setting HSTS
4. Security headers applied early
5. CORS after security (CORS adds headers, doesn't override)
6. Error handling and logging last

### HTTPS Detection

The middleware detects HTTPS through reverse proxy headers:

```go
// Checks these headers:
X-Forwarded-Proto: https
X-Forwarded-SSL: on
// Or direct TLS connection
```

This ensures HSTS is only set when actually using HTTPS.

---

## Testing Security Headers

### Manual Testing

```bash
# Test with curl
curl -I https://api.yourdomain.com/health

# Expected headers:
# X-Frame-Options: DENY
# X-Content-Type-Options: nosniff
# X-XSS-Protection: 1; mode=block
# Referrer-Policy: strict-origin-when-cross-origin
# Permissions-Policy: geolocation=(), microphone=(), camera=()...
# Strict-Transport-Security: max-age=31536000; includeSubDomains
# Content-Security-Policy: default-src 'none'; frame-ancestors 'none'
```

### Automated Testing

```bash
# Run middleware tests
go test ./internal/interfaces/http/middleware/security_headers_test.go \
        ./internal/interfaces/http/middleware/security_headers.go -v
```

### Online Security Scanners

1. **Security Headers** (https://securityheaders.com/)
   ```
   Visit: https://securityheaders.com/
   Enter: https://api.yourdomain.com
   Expected Grade: A or A+
   ```

2. **Mozilla Observatory** (https://observatory.mozilla.org/)
   ```
   Comprehensive security analysis
   Checks headers, SSL, and more
   ```

---

## Frontend Compatibility

### No Breaking Changes

The security headers are designed to not break frontend functionality:

1. **CORS Still Works**
   - Security headers applied before CORS
   - CORS headers added afterwards
   - No conflicts or overrides

2. **CSP is API-Friendly**
   - API endpoints have strict CSP that allows JSON responses
   - No script execution needed for API
   - Frontend runs in different origin (admin.yourdomain.com)

3. **HSTS Doesn't Affect Development**
   - Only enabled in production
   - Development can still use HTTP

### Frontend Requirements

The frontend should:
1. ✅ Use HTTPS in production
2. ✅ Not embed API in iframes (already blocked by X-Frame-Options)
3. ✅ Handle CORS properly (already implemented)
4. ✅ Use proper Content-Type headers (already doing)

**No changes needed in existing frontend code.**

---

## Configuration

### Environment Variables

```bash
# Production (HSTS enabled)
SUPROXY_ENVIRONMENT=production

# Development (HSTS disabled)
SUPROXY_ENVIRONMENT=development
```

### Behind Reverse Proxy

When behind Nginx or Traefik, ensure these headers are forwarded:

**Nginx:**
```nginx
proxy_set_header X-Forwarded-Proto $scheme;
proxy_set_header X-Forwarded-SSL on;
```

**Traefik:**
```yaml
# Automatically handles X-Forwarded-Proto
# No configuration needed
```

---

## Security Benefits

### Protection Against:

1. **Clickjacking** ✓
   - X-Frame-Options prevents framing attacks

2. **XSS (Cross-Site Scripting)** ✓
   - CSP prevents inline script execution
   - X-XSS-Protection provides legacy browser protection

3. **MIME Confusion** ✓
   - X-Content-Type-Options prevents type sniffing

4. **Information Leakage** ✓
   - Referrer-Policy controls what information is shared

5. **Man-in-the-Middle** ✓
   - HSTS forces HTTPS connections

6. **Unauthorized Device Access** ✓
   - Permissions-Policy blocks device APIs

---

## Compliance

### Standards Compliance

- ✅ **OWASP Top 10:** Addresses A2 (Cryptographic Failures), A5 (Security Misconfiguration)
- ✅ **NIST:** Follows secure coding guidelines
- ✅ **PCI DSS:** Meets secure transmission requirements
- ✅ **GDPR:** Privacy-enhancing (Referrer-Policy)

### Industry Best Practices

- ✅ **Mozilla Web Security Guidelines:** Level 2 compliance
- ✅ **OWASP Secure Headers Project:** All recommended headers
- ✅ **CWE:** Mitigates multiple Common Weakness Enumerations

---

## Troubleshooting

### Headers Not Appearing

**Problem:** Security headers not in response

**Solutions:**
```bash
# 1. Verify middleware is registered
grep "SecurityHeaders" internal/interfaces/http/router/router.go

# 2. Check order in middleware chain
# SecurityHeaders should be before CORS

# 3. Rebuild application
go build -o api cmd/api/main.go

# 4. Test locally
curl -I http://localhost:8080/health
```

### HSTS Not Working in Development

**Problem:** HSTS header not appearing in development

**Solution:**
```bash
# This is expected - HSTS only enabled in production
# Set environment to production for testing:
SUPROXY_ENVIRONMENT=production ./api
```

### CSP Blocking Frontend

**Problem:** CSP blocks frontend requests

**Solution:**
```
This shouldn't happen as API and frontend are separate origins.
CSP only applies to content loaded BY the API endpoint itself.
Frontend makes requests TO the API (not affected by API's CSP).

If issues occur:
1. Check if frontend is trying to embed API in iframe (not recommended)
2. Verify CORS is properly configured
3. Check browser console for actual CSP errors
```

---

## Monitoring

### Log Security Header Violations

While the backend sets headers, browsers enforce them. To monitor violations:

1. **Browser DevTools**
   ```
   Open: Chrome DevTools → Console
   Look for: CSP violation reports
   ```

2. **Reverse Proxy Logs**
   ```bash
   # Check for unusual patterns
   grep "X-Frame-Options" /var/log/nginx/access.log
   ```

3. **Security Scanners**
   ```bash
   # Regular scans
   curl -s https://api.yourdomain.com | grep -i "X-Frame"
   ```

---

## Maintenance

### Regular Checks

**Weekly:**
- Run automated tests
- Check security scanner scores

**Monthly:**
- Review CSP policies
- Update if frontend requirements change

**Quarterly:**
- Review all security headers
- Check for new security standards
- Update documentation

---

## References

- [OWASP Secure Headers Project](https://owasp.org/www-project-secure-headers/)
- [MDN Web Security](https://developer.mozilla.org/en-US/docs/Web/Security)
- [Content Security Policy](https://content-security-policy.com/)
- [Security Headers Scanner](https://securityheaders.com/)

---

**Document Version:** 1.0.0  
**Last Reviewed:** 2024-08-21  
**Next Review:** 2024-09-21
