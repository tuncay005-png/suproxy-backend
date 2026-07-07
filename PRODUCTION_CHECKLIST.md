# Production Deployment Checklist

## Pre-Deployment

### Security Configuration
- [ ] `.env.production` dosyası oluşturuldu
- [ ] `JWT_SECRET` 64+ karakter random string ile değiştirildi
- [ ] `DB_PASSWORD` güçlü ve benzersiz bir şifre ile değiştirildi
- [ ] `GRAFANA_PASSWORD` güvenli bir şifre ile değiştirildi
- [ ] `XRAY_USE_MOCK=false` olarak ayarlandı
- [ ] `DB_SSLMODE=require` olarak ayarlandı
- [ ] `ENVIRONMENT=production` olarak ayarlandı
- [ ] Tüm default değerler (CHANGE_ME) değiştirildi

### Infrastructure
- [ ] Docker 20.10+ kurulu
- [ ] Docker Compose 2.0+ kurulu
- [ ] Minimum sistem gereksinimleri karşılanıyor (4GB RAM, 20GB disk)
- [ ] Network portları açık (8080, 9090, 3000)
- [ ] Disk space monitoring ayarlandı
- [ ] Backup dizini oluşturuldu ve yazılabilir

### SSL/TLS (Önerilir)
- [ ] Reverse proxy (nginx/traefik) kuruldu
- [ ] SSL sertifikası yapılandırıldı
- [ ] HTTPS yönlendirmesi aktif
- [ ] SSL sertifikası otomatik yenileme ayarlandı

### Database
- [ ] PostgreSQL bağlantı parametreleri doğru
- [ ] Database backup stratejisi belirlendi
- [ ] Connection pool limitleri production için ayarlandı
- [ ] Database disk alanı yeterli (minimum 20GB)

## Deployment

### Build & Deploy
- [ ] Docker image başarıyla build edildi
- [ ] `deploy.sh` veya `deploy.ps1` çalıştırıldı
- [ ] Tüm container'lar çalışıyor (docker ps)
- [ ] Health check'ler geçiyor

### Service Verification
- [ ] API health endpoint erişilebilir: `curl http://localhost:8080/health`
- [ ] Database bağlantısı başarılı
- [ ] Migrations tamamlandı
- [ ] Prometheus metrics endpoint erişilebilir: `http://localhost:8080/metrics`
- [ ] Prometheus UI erişilebilir: `http://localhost:9090`
- [ ] Grafana UI erişilebilir: `http://localhost:3000`

### API Tests
- [ ] User registration çalışıyor
- [ ] User login çalışıyor
- [ ] JWT token validation çalışıyor
- [ ] Admin endpoints erişilebilir (admin user ile)
- [ ] Authorization kontrolleri çalışıyor

## Post-Deployment

### Monitoring Setup
- [ ] Grafana datasource yapılandırması kontrol edildi
- [ ] Prometheus targets UP durumunda
- [ ] Metrics akışı kontrol edildi
- [ ] Alert rules yapılandırıldı (varsa)
- [ ] Dashboard'lar yüklendi

### Logging
- [ ] Log formatı JSON olarak ayarlandı
- [ ] Log aggregation sistemi yapılandırıldı (varsa)
- [ ] Log rotation ayarları yapıldı
- [ ] Log retention policy belirlendi

### Backup
- [ ] İlk manual backup alındı ve test edildi
- [ ] Backup scripti çalışıyor
- [ ] Backup rotation (30 gün) test edildi
- [ ] Restore prosedürü test edildi
- [ ] Backup monitoring kuruldu

### Security Hardening
- [ ] Firewall kuralları yapılandırıldı
- [ ] Gereksiz portlar kapatıldı
- [ ] Container'lar non-root user ile çalışıyor
- [ ] Read-only filesystem kontrol edildi
- [ ] Security scan yapıldı (gosec, trivy)
- [ ] Dependency vulnerabilities tarandı

### Performance
- [ ] Resource limits ayarlandı (CPU, memory)
- [ ] Database connection pool optimize edildi
- [ ] Response time'lar kabul edilebilir seviyede
- [ ] Load testing yapıldı (opsiyonel)

### Disaster Recovery
- [ ] Disaster recovery planı hazırlandı
- [ ] Backup restore prosedürü dokümante edildi
- [ ] Rollback stratejisi belirlendi
- [ ] Contact list oluşturuldu (on-call)

## Maintenance

### Daily
- [ ] Health check status kontrol
- [ ] Error log review
- [ ] Disk space monitoring
- [ ] Service availability monitoring

### Weekly
- [ ] Backup verification
- [ ] Performance metrics review
- [ ] Database vacuum analyze
- [ ] Security log review

### Monthly
- [ ] Database reindex
- [ ] Docker image cleanup
- [ ] Dependency updates check
- [ ] Security patch review
- [ ] Incident report review

### Quarterly
- [ ] Disaster recovery drill
- [ ] Capacity planning review
- [ ] Security audit
- [ ] Documentation update

## Monitoring Alerts (Recommended)

### Critical Alerts
- [ ] API down (health check fails)
- [ ] Database connection failed
- [ ] Disk space > 80%
- [ ] Memory usage > 90%
- [ ] CPU usage sustained > 90%
- [ ] High error rate (>5%)

### Warning Alerts
- [ ] Slow response time (>1s)
- [ ] Database connection pool exhaustion
- [ ] Disk space > 70%
- [ ] Memory usage > 80%
- [ ] Backup failure

## Compliance

### GDPR (if applicable)
- [ ] Data retention policy implemented
- [ ] User data export capability
- [ ] User data deletion capability
- [ ] Privacy policy updated
- [ ] Cookie consent (if web UI exists)

### Audit
- [ ] Audit logging enabled
- [ ] Sensitive operations logged
- [ ] Admin actions tracked
- [ ] Log retention policy defined
- [ ] Access logs preserved

## Documentation

- [ ] API documentation updated
- [ ] Deployment guide reviewed
- [ ] Runbook created
- [ ] Architecture diagram updated
- [ ] Contact information documented
- [ ] Known issues documented

## Rollback Plan

### Rollback Triggers
- [ ] Critical bugs identified
- [ ] Performance degradation
- [ ] Security vulnerability
- [ ] Database corruption
- [ ] Service unavailability

### Rollback Steps
1. [ ] Identify previous stable version
2. [ ] Stop current version: `docker-compose down`
3. [ ] Update VERSION in .env.production
4. [ ] Deploy previous version: `./scripts/deploy.sh`
5. [ ] Verify health checks
6. [ ] Restore database backup (if needed)
7. [ ] Communicate to stakeholders

## Sign-off

- [ ] Dev Team Lead approval
- [ ] DevOps approval
- [ ] Security Team approval (if applicable)
- [ ] Product Owner informed
- [ ] Stakeholders notified

---

**Deployment Date:** _____________

**Deployed By:** _____________

**Version:** _____________

**Notes:**
_____________________________________________________________
_____________________________________________________________
_____________________________________________________________
