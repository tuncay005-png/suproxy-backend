# SuProxy Backend - Production Deployment Guide

## Overview

Bu döküman SuProxy Backend'in production ortamına deploy edilmesi için gereken tüm adımları içerir.

## Prerequisites

- Docker 20.10+
- Docker Compose 2.0+
- Minimum 4GB RAM
- Minimum 20GB disk space
- SSL sertifikası (production için önerilir)

## Quick Start (Development)

```bash
# .env dosyası oluştur
cp .env.example .env

# Servisleri başlat
docker-compose up -d

# Logları takip et
docker-compose logs -f api
```

API: http://localhost:8080

## Production Deployment

### 1. Environment Configuration

Production environment dosyasını oluştur:

```bash
cp .env.example .env.production
```

`.env.production` dosyasını düzenle ve aşağıdaki değerleri güvenli değerlerle değiştir:

**KRİTİK GÜVENLİK AYARLARI:**
- `DB_PASSWORD`: Güçlü database şifresi
- `JWT_SECRET`: 64+ karakter random string
- `GRAFANA_PASSWORD`: Güçlü Grafana şifresi
- `XRAY_USE_MOCK`: Production'da `false` olmalı
- `DB_SSLMODE`: Production'da `require` olmalı

**JWT Secret oluşturma:**
```bash
# Linux/Mac
openssl rand -base64 64

# Windows (PowerShell)
[Convert]::ToBase64String((1..64 | ForEach-Object { Get-Random -Minimum 0 -Maximum 256 }))
```

### 2. Docker Image Build

```bash
# Linux/Mac
./scripts/deploy.sh

# Windows (PowerShell)
.\scripts\deploy.ps1
```

veya manuel:

```bash
# Image build
docker build -t suproxy/backend:1.0.0 .

# Production servisleri başlat
docker-compose -f docker-compose.production.yml up -d
```

### 3. Health Check

```bash
# API health check
curl http://localhost:8080/health

# Container durumu
docker-compose -f docker-compose.production.yml ps

# Logları kontrol et
docker-compose -f docker-compose.production.yml logs api
```

### 4. Monitoring Setup

**Prometheus:** http://localhost:9090
**Grafana:** http://localhost:3000

Grafana default credentials:
- Username: admin
- Password: `.env.production` dosyasında tanımlı

## Container Architecture

### Services

1. **postgres**: PostgreSQL 15 database
   - Port: 5432
   - Volume: `postgres_data`
   - Health check: pg_isready
   - Resources: 2 CPU, 2GB RAM (production)

2. **api**: SuProxy Backend API
   - Port: 8080
   - Volumes: `xray_config`, `xray_logs`, `xray_backups`
   - Health check: `/health` endpoint
   - Resources: 4 CPU, 4GB RAM (production)
   - Security: non-root user, read-only filesystem

3. **prometheus**: Metrics collection
   - Port: 9090
   - Volume: `prometheus_data`
   - Retention: 30 days

4. **grafana**: Metrics visualization
   - Port: 3000
   - Volume: `grafana_data`

### Volumes

- `postgres_data`: Database persistent storage
- `xray_config`: Xray configuration files
- `xray_logs`: Xray logs
- `xray_backups`: Xray configuration backups
- `prometheus_data`: Prometheus time-series data
- `grafana_data`: Grafana dashboards and settings

### Network

- Network: `suproxy-network` (bridge)
- Subnet: 172.25.0.0/16 (production)

## Security Hardening

### Container Security

1. **Non-root User**: API container `suproxy:1000` kullanıcısı ile çalışır
2. **Read-only Filesystem**: API container read-only mode'da çalışır
3. **No New Privileges**: `security-opt: no-new-privileges:true`
4. **Resource Limits**: CPU ve memory limitleri tanımlı
5. **Minimal Image**: Alpine-based image (~30MB)

### Network Security

1. **Internal Network**: Tüm servisler izole bridge network'te
2. **Port Exposure**: Sadece gerekli portlar expose edilir
3. **SSL/TLS**: Production'da reverse proxy arkasında HTTPS kullan

### Database Security

1. **SSL Mode**: Production'da `require` veya `verify-full`
2. **Strong Password**: Minimum 16 karakter complex password
3. **Connection Pooling**: Max connection limitleri tanımlı
4. **Regular Backups**: Otomatik backup scriptleri

## Backup & Restore

### Database Backup

```bash
# Linux/Mac
./scripts/backup.sh

# Windows (PowerShell)
.\scripts\backup.ps1
```

Backup dosyaları `./backups` dizininde saklanır.
Otomatik olarak 30 günden eski backuplar silinir.

### Manual Backup

```bash
docker-compose -f docker-compose.production.yml exec postgres pg_dump \
    -U suproxy_prod \
    -d suproxy_prod \
    --format=plain \
    --no-owner \
    --no-acl > backup.sql
```

### Restore

```bash
docker-compose -f docker-compose.production.yml exec -T postgres psql \
    -U suproxy_prod \
    -d suproxy_prod < backup.sql
```

## Resource Requirements

### Minimum Requirements
- CPU: 2 cores
- RAM: 4GB
- Disk: 20GB
- Network: 100Mbps

### Recommended (Production)
- CPU: 4+ cores
- RAM: 8GB+
- Disk: 50GB+ SSD
- Network: 1Gbps

## Scaling

### Horizontal Scaling

API servisi horizontal scale edilebilir:

```yaml
# docker-compose.production.yml
services:
  api:
    deploy:
      replicas: 3
```

Load balancer (nginx, traefik) ile birlikte kullanın.

### Vertical Scaling

Resource limitleri artırın:

```yaml
deploy:
  resources:
    limits:
      cpus: '8'
      memory: 8G
```

## Monitoring & Observability

### Prometheus Metrics

- `/metrics` endpoint: API metrics
- HTTP request metrics: duration, count, errors
- Business metrics: active users, clients, xray instances
- Database metrics: connection pool, query performance
- Runtime metrics: Go runtime statistics

### Grafana Dashboards

Prometheus data source otomatik olarak provision edilir.
Custom dashboardlar `./grafana/provisioning/dashboards/` dizinine eklenebilir.

### Logging

API JSON format structured logging kullanır:
- Request ID tracking
- Correlation ID support
- Log levels: debug, info, warn, error
- Log aggregation için stdout/stderr kullanır

## Troubleshooting

### API çalışmıyor

```bash
# Container durumu
docker-compose -f docker-compose.production.yml ps

# API logs
docker-compose -f docker-compose.production.yml logs api

# Database connection test
docker-compose -f docker-compose.production.yml exec api wget -O- http://localhost:8080/health
```

### Database bağlantı hatası

```bash
# Database logs
docker-compose -f docker-compose.production.yml logs postgres

# Connection test
docker-compose -f docker-compose.production.yml exec postgres pg_isready -U suproxy_prod
```

### Yüksek memory kullanımı

```bash
# Resource kullanımı
docker stats

# Connection pool ayarları kontrol et
# .env.production içinde DB_MAX_OPEN_CONNS ve DB_MAX_IDLE_CONNS değerlerini düşür
```

### Disk doldu

```bash
# Disk kullanımı
docker system df

# Eski container/image cleanup
docker system prune -a

# Eski backupları sil
find ./backups -name "*.sql.gz" -mtime +30 -delete
```

## Graceful Shutdown

API graceful shutdown destekler:

```bash
# SIGTERM gönder (docker-compose down kullanır)
docker-compose -f docker-compose.production.yml down

# Veya container'a direct signal
docker kill --signal=SIGTERM suproxy-api-prod
```

Shutdown timeout: 30 saniye (configurable)

## Updates & Rollback

### Update

```bash
# Yeni version build et
docker build -t suproxy/backend:1.1.0 .

# .env.production içinde VERSION=1.1.0 yap

# Rolling update
docker-compose -f docker-compose.production.yml up -d
```

### Rollback

```bash
# Önceki version'a dön
# .env.production içinde VERSION=1.0.0 yap

# Restart
docker-compose -f docker-compose.production.yml up -d
```

## Production Checklist

- [ ] `.env.production` oluşturuldu ve güvenli değerler girildi
- [ ] JWT_SECRET 64+ karakter random string
- [ ] Database password güçlü ve unique
- [ ] XRAY_USE_MOCK=false
- [ ] DB_SSLMODE=require
- [ ] Grafana password değiştirildi
- [ ] Backup scripti test edildi
- [ ] Health check endpointleri çalışıyor
- [ ] Prometheus metrics erişilebilir
- [ ] Grafana dashboards yüklendi
- [ ] SSL/TLS certificate yapılandırıldı (reverse proxy)
- [ ] Firewall kuralları ayarlandı
- [ ] Monitoring alertleri yapılandırıldı
- [ ] Backup retention policy belirlendi
- [ ] Disaster recovery planı hazırlandı

## Support & Maintenance

### Log Rotation

Container logları Docker tarafından yönetilir. Production'da log driver yapılandırın:

```yaml
services:
  api:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

### Database Maintenance

```bash
# Vacuum analyze (haftalık)
docker-compose -f docker-compose.production.yml exec postgres \
    psql -U suproxy_prod -d suproxy_prod -c "VACUUM ANALYZE;"

# Reindex (aylık)
docker-compose -f docker-compose.production.yml exec postgres \
    psql -U suproxy_prod -d suproxy_prod -c "REINDEX DATABASE suproxy_prod;"
```

## License

Copyright © 2024 SuProxy. All rights reserved.
