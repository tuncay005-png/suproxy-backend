# Production HTTPS Setup Guide

This guide provides complete instructions for setting up HTTPS with SSL/TLS certificates for the SuProxy backend API.

## Overview

Two reverse proxy options are provided:
1. **Traefik** (Recommended) - Automatic SSL management, Docker-native
2. **Nginx** - Traditional, widely used, more manual control

Choose based on your preference and infrastructure.

---

## Prerequisites

✅ **Before Starting:**
- [ ] Domain name configured (e.g., `yourdomain.com`)
- [ ] DNS A record pointing to your server IP
  - `api.yourdomain.com` → `YOUR_SERVER_IP`
- [ ] Ports 80 and 443 open in firewall
- [ ] Docker and Docker Compose installed
- [ ] Server has internet access
- [ ] Email address for Let's Encrypt notifications

**DNS Verification:**
```bash
# Verify DNS is properly configured
nslookup api.yourdomain.com
# Should return your server IP
```

---

## Option 1: Traefik (Recommended)

### Why Traefik?
✅ Automatic SSL certificate management  
✅ Auto-renewal built-in  
✅ Docker-native integration  
✅ Dashboard for monitoring  
✅ Simpler configuration  

### Setup Steps

#### 1. Configure Environment

```bash
cd /opt/suproxy/backend
cp .env.production.example .env.production
nano .env.production
```

**Required settings:**
```bash
# Your domain (without https://)
DOMAIN=yourdomain.com

# Email for Let's Encrypt notifications
ACME_EMAIL=admin@yourdomain.com

# Generate secure passwords
JWT_SECRET=$(openssl rand -base64 32)
DB_PASSWORD=$(openssl rand -base64 32)
GRAFANA_PASSWORD=$(openssl rand -base64 16)

# Generate basic auth for Traefik dashboard
# htpasswd -nb admin your_password
TRAEFIK_DASHBOARD_AUTH='admin:$apr1$...'
MONITORING_AUTH='admin:$apr1$...'
```

#### 2. Initial Setup with Staging Certificates

**Why staging first?**  
Let's Encrypt has rate limits (5 certificates per domain per week). Testing with staging prevents hitting limits.

```bash
# Run setup script
chmod +x scripts/setup-https-traefik.sh
./scripts/setup-https-traefik.sh
```

The script will:
1. Verify configuration
2. Create necessary volumes
3. Start Traefik with staging certificates
4. Wait for certificate generation

#### 3. Verify Staging Setup

```bash
# Check Traefik logs
docker logs suproxy-traefik

# Test API (will show certificate warning - expected)
curl -k https://api.yourdomain.com/health

# Should return: {"status":"ok"}
```

**Expected:** Browser will show certificate warning (staging cert is not trusted).

#### 4. Enable Production Certificates

Once staging works:

```bash
# 1. Stop services
docker-compose -f docker-compose.production-https.yml down

# 2. Enable production certificates
# Edit traefik/traefik.yml and comment out staging line:
nano traefik/traefik.yml
# Comment: # caServer: https://acme-staging-v02.api.letsencrypt.org/directory

# 3. Remove staging certificates
docker volume rm traefik_letsencrypt
docker volume create traefik_letsencrypt

# 4. Start with production certificates
docker-compose -f docker-compose.production-https.yml up -d

# 5. Wait for certificate generation
sleep 30

# 6. Verify production certificate
curl https://api.yourdomain.com/health
# Should work without certificate warning
```

#### 5. Verify HTTPS

```bash
# Test HTTPS endpoint
curl https://api.yourdomain.com/health

# Check security headers
curl -I https://api.yourdomain.com/health

# Should include:
# Strict-Transport-Security: max-age=31536000
# X-Frame-Options: DENY
# X-Content-Type-Options: nosniff

# Test HTTP -> HTTPS redirect
curl -I http://api.yourdomain.com/health
# Should return: 301 Moved Permanently
# Location: https://api.yourdomain.com/health
```

#### 6. Access Traefik Dashboard

```bash
# Dashboard available at: https://traefik.yourdomain.com
# Login with credentials from TRAEFIK_DASHBOARD_AUTH

# Shows:
# - Active routes
# - SSL certificate status
# - Backend health
# - Request metrics
```

### Traefik Architecture

```
Internet
   │
   ├─ Port 80 (HTTP)
   │     │
   │     └─> Traefik ──> Redirect to HTTPS
   │
   └─ Port 443 (HTTPS + TLS)
         │
         └─> Traefik
               ├─> Let's Encrypt (auto SSL)
               │
               ├─ api.yourdomain.com
               │    └─> API Container :8080
               │
               ├─ grafana.yourdomain.com
               │    └─> Grafana :3000
               │
               └─ prometheus.yourdomain.com
                    └─> Prometheus :9090
```

---

## Option 2: Nginx

### Why Nginx?
✅ Traditional, widely known  
✅ Fine-grained control  
✅ Mature ecosystem  
✅ Extensive documentation  

### Setup Steps

#### 1. Configure Environment

Same as Traefik - configure `.env.production`

#### 2. Run Setup Script

```bash
chmod +x scripts/setup-https-nginx.sh
./scripts/setup-https-nginx.sh
```

The script will:
1. Update Nginx configuration with your domain
2. Start backend services
3. Obtain SSL certificate via Certbot
4. Start Nginx reverse proxy
5. Configure auto-renewal

#### 3. Verify Setup

```bash
# Test HTTPS
curl https://api.yourdomain.com/health

# Check certificate
openssl s_client -connect api.yourdomain.com:443 -servername api.yourdomain.com < /dev/null 2>/dev/null | openssl x509 -noout -dates

# Test HTTP redirect
curl -I http://api.yourdomain.com/health
```

#### 4. Certificate Renewal

Certificates auto-renew via certbot container (runs every 12 hours).

**Manual renewal (if needed):**
```bash
docker-compose -f docker-compose.production-nginx.yml exec certbot certbot renew
docker-compose -f docker-compose.production-nginx.yml exec nginx nginx -s reload
```

### Nginx Architecture

```
Internet
   │
   ├─ Port 80 (HTTP)
   │     │
   │     └─> Nginx ──> Redirect to HTTPS (except /.well-known)
   │
   └─ Port 443 (HTTPS + TLS)
         │
         └─> Nginx
               ├─> SSL Termination
               ├─> Security Headers
               ├─> Rate Limiting
               └─> Proxy to API :8080
```

---

## Security Features Implemented

### 1. HTTPS/TLS
- ✅ TLS 1.2 and 1.3 only
- ✅ Strong cipher suites
- ✅ Perfect Forward Secrecy
- ✅ OCSP Stapling (Nginx)

### 2. HTTP Strict Transport Security (HSTS)
```
Strict-Transport-Security: max-age=31536000; includeSubDomains; preload
```
Forces HTTPS for 1 year, includes subdomains.

### 3. Security Headers

| Header | Value | Purpose |
|--------|-------|---------|
| X-Frame-Options | DENY | Prevent clickjacking |
| X-Content-Type-Options | nosniff | Prevent MIME sniffing |
| X-XSS-Protection | 1; mode=block | XSS protection |
| Referrer-Policy | strict-origin-when-cross-origin | Limit referrer info |
| Permissions-Policy | geolocation=(), microphone=(), camera=() | Disable unnecessary APIs |
| Content-Security-Policy | default-src 'self'... | Restrict resource loading |

### 4. Cookie Security

Cookies are configured with:
```go
COOKIE_SECURE=true      // HTTPS only
COOKIE_SAMESITE=strict  // Strict same-site policy
```

### 5. Rate Limiting

**Traefik:**
- 100 requests/minute average
- 200 burst capacity

**Nginx:**
- Auth endpoints: 5 req/min per IP
- General API: 100 req/min per IP

---

## Monitoring HTTPS

### Certificate Expiration

**Traefik:**
- Auto-renews 30 days before expiration
- Check dashboard: https://traefik.yourdomain.com

**Nginx:**
```bash
# Check certificate expiration
docker-compose -f docker-compose.production-nginx.yml exec certbot certbot certificates

# Expected output shows expiration date
```

### SSL Labs Test

Test your SSL configuration:
```bash
# Visit: https://www.ssllabs.com/ssltest/
# Enter: api.yourdomain.com
# Expected Grade: A or A+
```

### Check HTTPS is Working

```bash
# Health check
curl https://api.yourdomain.com/health

# Login endpoint (should reject without credentials)
curl https://api.yourdomain.com/api/v1/auth/login

# Check security headers
curl -I https://api.yourdomain.com/health | grep -i "strict-transport\|x-frame\|x-content"
```

---

## Troubleshooting

### Issue: Certificate Not Generated

**Symptoms:**
- Browser shows "Connection not secure"
- 404 errors on HTTPS

**Solutions:**
```bash
# 1. Check DNS
nslookup api.yourdomain.com

# 2. Check Traefik/Certbot logs
docker logs suproxy-traefik  # Traefik
docker logs suproxy-certbot  # Nginx

# 3. Verify port 80 is accessible
curl http://api.yourdomain.com/.well-known/acme-challenge/test

# 4. Check firewall
sudo ufw status
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
```

### Issue: HTTP Not Redirecting to HTTPS

**Traefik:**
```bash
# Check redirect configuration
docker exec suproxy-traefik cat /etc/traefik/traefik.yml | grep -A5 "web:"

# Restart Traefik
docker-compose -f docker-compose.production-https.yml restart traefik
```

**Nginx:**
```bash
# Check nginx config
docker exec suproxy-nginx cat /etc/nginx/conf.d/default.conf | grep -A10 "listen 80"

# Test config
docker exec suproxy-nginx nginx -t

# Reload nginx
docker-compose -f docker-compose.production-nginx.yml exec nginx nginx -s reload
```

### Issue: Mixed Content Warnings

**Cause:** Frontend making HTTP requests to HTTPS API

**Solution:**
```javascript
// Update frontend API base URL
const API_BASE_URL = 'https://api.yourdomain.com'
```

### Issue: Rate Limit Errors

**Symptoms:** 
- Let's Encrypt error: "too many certificates"

**Solution:**
```bash
# Use staging environment first
# Edit traefik/traefik.yml:
caServer: https://acme-staging-v02.api.letsencrypt.org/directory

# Wait 1 week or use different (sub)domain
```

---

## Maintenance

### Certificate Renewal

**Automatic (both solutions):**
- Traefik: Checks every 24 hours, renews 30 days before expiry
- Nginx/Certbot: Runs every 12 hours, renews 30 days before expiry

**Manual Check:**
```bash
# Traefik - check dashboard or logs
docker logs suproxy-traefik | grep -i certificate

# Nginx - force renewal check
docker-compose -f docker-compose.production-nginx.yml exec certbot certbot renew --dry-run
```

### Updating Configuration

**Traefik:**
```bash
# 1. Edit configuration files
nano traefik/traefik.yml
nano traefik/dynamic/middleware.yml

# 2. Restart Traefik
docker-compose -f docker-compose.production-https.yml restart traefik
```

**Nginx:**
```bash
# 1. Edit nginx config
nano nginx/suproxy.conf

# 2. Test config
docker-compose -f docker-compose.production-nginx.yml exec nginx nginx -t

# 3. Reload nginx
docker-compose -f docker-compose.production-nginx.yml exec nginx nginx -s reload
```

### Log Rotation

Logs are stored in Docker volumes. Configure log rotation:

```bash
# Create /etc/docker/daemon.json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}

# Restart Docker
sudo systemctl restart docker
```

---

## Production Checklist

Before going live:

### DNS & Domain
- [ ] Domain DNS points to server IP
- [ ] Propagation verified (24-48 hours)
- [ ] Subdomain configured (api.yourdomain.com)

### SSL/TLS
- [ ] Certificate obtained successfully
- [ ] HTTPS accessible (no warnings)
- [ ] HTTP redirects to HTTPS
- [ ] SSL Labs grade A or A+
- [ ] Certificate auto-renewal configured

### Security
- [ ] Security headers present
- [ ] Cookie flags set (secure, httponly, samesite)
- [ ] Rate limiting active
- [ ] Firewall configured
- [ ] Sensitive ports closed (PostgreSQL, Prometheus)

### Functionality
- [ ] Health endpoint works: https://api.yourdomain.com/health
- [ ] API endpoints accessible
- [ ] Frontend can connect
- [ ] Authentication works
- [ ] Database connections stable

### Monitoring
- [ ] Traefik dashboard accessible (if using Traefik)
- [ ] Grafana dashboard accessible
- [ ] Certificate expiration monitored
- [ ] Logs accessible

---

## Comparison: Traefik vs Nginx

| Feature | Traefik | Nginx |
|---------|---------|-------|
| Setup Complexity | ⭐⭐ Easy | ⭐⭐⭐ Moderate |
| SSL Management | 🤖 Automatic | ✋ Semi-automatic |
| Docker Integration | ✅ Native | ⚠️ Manual |
| Configuration | YAML/Labels | Config files |
| Dashboard | ✅ Built-in | ❌ None |
| Performance | ⭐⭐⭐⭐ Good | ⭐⭐⭐⭐⭐ Excellent |
| Community | ⭐⭐⭐ Growing | ⭐⭐⭐⭐⭐ Huge |
| Learning Curve | Gentle | Steep |

**Recommendation:**
- **Traefik**: If you want automation and Docker-native experience
- **Nginx**: If you need maximum control and performance

---

## Support & Resources

### Let's Encrypt
- Rate Limits: https://letsencrypt.org/docs/rate-limits/
- Staging: https://letsencrypt.org/docs/staging-environment/
- Troubleshooting: https://community.letsencrypt.org/

### Traefik
- Docs: https://doc.traefik.io/traefik/
- ACME: https://doc.traefik.io/traefik/https/acme/

### Nginx
- Docs: https://nginx.org/en/docs/
- SSL: https://nginx.org/en/docs/http/configuring_https_servers.html

### Security Testing
- SSL Labs: https://www.ssllabs.com/ssltest/
- Security Headers: https://securityheaders.com/

---

## Next Steps

After HTTPS is configured:

1. **Update Frontend**
   ```javascript
   // Change API base URL to HTTPS
   const API_URL = 'https://api.yourdomain.com'
   ```

2. **Update Documentation**
   - Update API documentation with HTTPS URLs
   - Update deployment guides
   - Inform team members

3. **Configure Monitoring**
   - Set up alerts for certificate expiration
   - Monitor HTTPS traffic
   - Track SSL handshake metrics

4. **Test Load**
   - Run load tests over HTTPS
   - Verify performance is acceptable
   - Check SSL overhead is minimal

5. **Backup Configuration**
   ```bash
   # Backup certificates and config
   tar -czf https-backup-$(date +%Y%m%d).tar.gz \
     .env.production \
     traefik/ \
     nginx/ \
     certbot/conf/
   ```

---

## Emergency Rollback

If HTTPS causes issues:

```bash
# Revert to HTTP-only
docker-compose -f docker-compose.production.yml up -d

# Investigate and fix
# Then retry HTTPS setup
```

**Important:** Don't run HTTP and HTTPS simultaneously - choose one.

---

**Documentation Version:** 1.0.0  
**Last Updated:** 2024-08-21  
**Tested With:** Traefik 2.10, Nginx 1.25, Certbot 2.x
