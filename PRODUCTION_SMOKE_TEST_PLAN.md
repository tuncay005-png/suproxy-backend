# 🔥 Production Smoke Test Plan - Suproxy Backend

**Test Tarihi:** -  
**Test Eden:** -  
**Environment:** Production  
**Base URL:** `http://localhost:8080` (veya production URL)

---

## ✅ Test Checklist

### 1️⃣ Server Startup
**Amaç:** Backend'in düzgün başladığından emin olmak

#### Test Adımları:
```bash
# Docker ile başlatma
docker-compose up -d

# Logları izleme
docker-compose logs -f api

# Container durumu kontrolü
docker-compose ps
```

#### Beklenen Sonuç:
- ✅ Container `healthy` durumda
- ✅ Log'da "Server started on :8080" mesajı var
- ✅ Database bağlantısı başarılı
- ❌ Hata/panic yok

**Durum:** [ ] Başarılı [ ] Başarısız

**Notlar:**
```

```

---

### 2️⃣ Health Endpoint
**Amaç:** Health check endpoint'lerinin çalıştığını doğrulamak

#### Test 1: Public Health Check
```bash
curl -X GET http://localhost:8080/health
```

**Beklenen Response (200 OK):**
```json
{
  "status": "healthy",
  "timestamp": "2026-07-16T12:00:00Z"
}
```

#### Test 2: Readiness Check
```bash
curl -X GET http://localhost:8080/ready
```

**Beklenen Response (200 OK):**
```json
{
  "status": "ready",
  "checks": {
    "database": "ok",
    "redis": "ok"
  }
}
```

**Durum:** [ ] Başarılı [ ] Başarısız

**Notlar:**
```

```

---

### 3️⃣ Auth - Register Flow
**Amaç:** Kullanıcı kayıt işleminin çalıştığını doğrulamak

#### API Endpoint:
- **Method:** `POST`
- **URL:** `/api/v1/auth/register`
- **Authentication:** ❌ Gerekli değil (Public)

#### Test Komutu:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "SecurePass123!",
    "first_name": "Test",
    "last_name": "User"
  }'
```

#### Beklenen Response (201 Created):
```json
{
  "user": {
    "id": "uuid-here",
    "email": "testuser@example.com",
    "first_name": "Test",
    "last_name": "User",
    "status": "active",
    "role": "user",
    "created_at": "2026-07-16T12:00:00Z",
    "updated_at": "2026-07-16T12:00:00Z"
  },
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc..."
}
```

**Durum:** [ ] Başarılı [ ] Başarısız

**Access Token:** `____________________`  
**Refresh Token:** `____________________`

**Notlar:**
```

```

---

### 4️⃣ Auth - Login Flow
**Amaç:** Kullanıcı login işleminin çalıştığını doğrulamak

#### API Endpoint:
- **Method:** `POST`
- **URL:** `/api/v1/auth/login`
- **Authentication:** ❌ Gerekli değil (Public)

#### Test Komutu:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "SecurePass123!",
    "device_name": "Test Device",
    "platform": "curl"
  }'
```

#### Beklenen Response (200 OK):
```json
{
  "user": {
    "id": "uuid-here",
    "email": "testuser@example.com",
    "first_name": "Test",
    "last_name": "User",
    "status": "active",
    "role": "user"
  },
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc..."
}
```

**Durum:** [ ] Başarılı [ ] Başarısız

**Notlar:**
```

```

---

### 5️⃣ JWT Access-Refresh Flow
**Amaç:** Token yenileme mekanizmasının çalıştığını doğrulamak

#### Test 1: Get Current User (Access Token)
```bash
export ACCESS_TOKEN="your-access-token-here"

curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

**Beklenen Response (200 OK):**
```json
{
  "id": "uuid-here",
  "email": "testuser@example.com",
  "first_name": "Test",
  "last_name": "User",
  "role": "user",
  "status": "active"
}
```

#### Test 2: Refresh Token
```bash
export REFRESH_TOKEN="your-refresh-token-here"

curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }"
```

**Beklenen Response (200 OK):**
```json
{
  "access_token": "new-access-token",
  "refresh_token": "new-refresh-token"
}
```

**Durum:** [ ] Başarılı [ ] Başarısı

**Notlar:**
```

```

---

### 6️⃣ Admin User Oluşturma
**Amaç:** Admin kullanıcısı oluşturup admin yetkilerini test etmek

#### Senaryo A: Database'den Manuel Oluşturma
```sql
-- PostgreSQL'e bağlan
psql -U suproxy_user -d suproxy_db

-- Admin kullanıcısını manuel oluştur
UPDATE users 
SET role = 'admin' 
WHERE email = 'testuser@example.com';

-- Kontrol et
SELECT id, email, role FROM users WHERE email = 'testuser@example.com';
```

#### Senaryo B: Environment Variable ile İlk Admin
```bash
# .env dosyasına ekle
ADMIN_EMAIL=admin@example.com
ADMIN_PASSWORD=AdminSecure123!

# Container'ı yeniden başlat
docker-compose restart api
```

#### Test: Admin Health Check
```bash
export ADMIN_TOKEN="admin-access-token"

curl -X GET http://localhost:8080/api/v1/admin/health \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

**Beklenen Response (200 OK):**
```json
{
  "status": "healthy",
  "role": "admin"
}
```

**Durum:** [ ] Başarılı [ ] Başarısız

**Admin Token:** `____________________`

**Notlar:**
```

```

---

### 7️⃣ Client Oluşturma
**Amaç:** Xray client oluşturma işlemini test etmek

#### Ön Koşul:
- ✅ Admin token hazır
- ✅ Xray instance ID var
- ✅ Inbound ID var
- ✅ User ID var

#### API Endpoint:
- **Method:** `POST`
- **URL:** `/api/v1/admin/xray/clients`
- **Authentication:** ✅ Gerekli (Admin JWT)

#### Request Body:
```json
{
  "inbound_id": "uuid-inbound-id",
  "user_id": "uuid-user-id",
  "email": "client@example.com",
  "flow": "xtls-rprx-vision"
}
```

#### Test Komutu:
```bash
export ADMIN_TOKEN="your-admin-token"
export INBOUND_ID="your-inbound-id"
export USER_ID="your-user-id"

curl -X POST http://localhost:8080/api/v1/admin/xray/clients \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"inbound_id\": \"$INBOUND_ID\",
    \"user_id\": \"$USER_ID\",
    \"email\": \"client@example.com\",
    \"flow\": \"xtls-rprx-vision\"
  }"
```

#### Beklenen Response (201 Created):
```json
{
  "id": "client-uuid",
  "inbound_id": "inbound-uuid",
  "user_id": "user-uuid",
  "uuid": "generated-uuid-for-xray",
  "flow": "xtls-rprx-vision",
  "email": "client@example.com",
  "enabled": true,
  "created_at": "2026-07-16T12:00:00Z",
  "updated_at": "2026-07-16T12:00:00Z"
}
```

**Durum:** [ ] Başarılı [ ] Başarısız

**Client ID:** `____________________`  
**Client UUID:** `____________________`

**Notlar:**
```

```

---

### 8️⃣ Xray Inbound Oluşturma
**Amaç:** Xray inbound yapılandırmasını test etmek

#### Ön Koşul:
- ✅ Admin token hazır
- ✅ Xray instance ID var

#### API Endpoint:
- **Method:** `POST`
- **URL:** `/api/v1/admin/xray/inbounds`
- **Authentication:** ✅ Gerekli (Admin JWT)

#### Request Body:
```json
{
  "xray_instance_id": "uuid-instance-id",
  "protocol": "vless",
  "port": 8443,
  "transport": "tcp",
  "security": "tls"
}
```

#### Test Komutu:
```bash
export ADMIN_TOKEN="your-admin-token"
export INSTANCE_ID="your-instance-id"

curl -X POST http://localhost:8080/api/v1/admin/xray/inbounds \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"xray_instance_id\": \"$INSTANCE_ID\",
    \"protocol\": \"vless\",
    \"port\": 8443,
    \"transport\": \"tcp\",
    \"security\": \"tls\"
  }"
```

#### Beklenen Response (201 Created):
```json
{
  "id": "inbound-uuid",
  "xray_instance_id": "instance-uuid",
  "protocol": "vless",
  "port": 8443,
  "transport": "tcp",
  "security": "tls",
  "enabled": true,
  "created_at": "2026-07-16T12:00:00Z",
  "updated_at": "2026-07-16T12:00:00Z"
}
```

**Durum:** [ ] Başarılı [ ] Başarısız

**Inbound ID:** `____________________`

**Notlar:**
```

```

---

### 9️⃣ VLESS Config Üretimi
**Amaç:** Client için VLESS configuration link'inin üretildiğini doğrulamak

#### Test: Client Detaylarını Getir
```bash
export ADMIN_TOKEN="your-admin-token"
export CLIENT_ID="your-client-id"

curl -X GET http://localhost:8080/api/v1/admin/xray/clients/$CLIENT_ID \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

#### Beklenen Response (200 OK):
```json
{
  "id": "client-uuid",
  "inbound_id": "inbound-uuid",
  "user_id": "user-uuid",
  "uuid": "vless-uuid-here",
  "flow": "xtls-rprx-vision",
  "email": "client@example.com",
  "enabled": true,
  "created_at": "2026-07-16T12:00:00Z",
  "updated_at": "2026-07-16T12:00:00Z"
}
```

#### Manual VLESS Link Oluşturma:
```
vless://{client-uuid}@{server-ip}:{port}?type=tcp&security=tls&flow=xtls-rprx-vision#{client-email}
```

**Örnek:**
```
vless://550e8400-e29b-41d4-a716-446655440000@example.com:8443?type=tcp&security=tls&flow=xtls-rprx-vision#client@example.com
```

**Durum:** [ ] Başarılı [ ] Başarısız

**VLESS Link:** `____________________`

**Notlar:**
```

```

---

### 🔟 Client Silme
**Amaç:** Client silme işleminin düzgün çalıştığını doğrulamak

#### API Endpoint:
- **Method:** `DELETE`
- **URL:** `/api/v1/admin/xray/clients/:id`
- **Authentication:** ✅ Gerekli (Admin JWT)

#### Test Komutu:
```bash
export ADMIN_TOKEN="your-admin-token"
export CLIENT_ID="your-client-id"

curl -X DELETE http://localhost:8080/api/v1/admin/xray/clients/$CLIENT_ID \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

#### Beklenen Response (200 OK):
```json
{
  "success": true,
  "message": "Client deleted successfully"
}
```

#### Doğrulama: Client'ın silindiğini kontrol et
```bash
curl -X GET http://localhost:8080/api/v1/admin/xray/clients/$CLIENT_ID \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

**Beklenen Response (404 Not Found):**
```json
{
  "error": "client not found"
}
```

**Durum:** [ ] Başarılı [ ] Başarısız

**Notlar:**
```

```

---

### 1️⃣1️⃣ Xray Config Cleanup
**Amaç:** Xray config dosyalarının düzgün temizlendiğini doğrulamak

#### Test 1: Config Dosyası Kontrolü
```bash
# Xray config dizinine gir
docker-compose exec api ls -la /etc/xray/

# Config içeriğini kontrol et
docker-compose exec api cat /etc/xray/config.json
```

**Beklenen Sonuç:**
- ✅ Silinen client'lar config'de yok
- ✅ JSON syntax geçerli
- ✅ Aktif inbound'lar mevcut

#### Test 2: Xray Instance Reload
```bash
export ADMIN_TOKEN="your-admin-token"
export INSTANCE_ID="your-instance-id"

curl -X POST http://localhost:8080/api/v1/admin/xray/instances/$INSTANCE_ID/reload \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

**Beklenen Response (200 OK):**
```json
{
  "success": true,
  "message": "Instance reloaded successfully"
}
```

**Durum:** [ ] Başarılı [ ] Başarısız

**Notlar:**
```

```

---

## 📋 Admin Endpoint Listesi

### 👤 User Management (`/api/v1/admin/users`)
| Method | Endpoint | Açıklama | Auth |
|--------|----------|----------|------|
| GET | `/api/v1/admin/users` | Kullanıcıları listele | ✅ Admin |
| GET | `/api/v1/admin/users/:id` | Kullanıcı detayı | ✅ Admin |
| PUT | `/api/v1/admin/users/:id/status` | Kullanıcı durumu güncelle | ✅ Admin |
| PUT | `/api/v1/admin/users/:id/role` | Kullanıcı rolü güncelle | ✅ Admin |

### 🖥️ Xray Instance Management (`/api/v1/admin/xray/instances`)
| Method | Endpoint | Açıklama | Auth |
|--------|----------|----------|------|
| GET | `/api/v1/admin/xray/instances` | Instance'ları listele | ✅ Admin |
| GET | `/api/v1/admin/xray/instances/:id` | Instance detayı | ✅ Admin |
| POST | `/api/v1/admin/xray/instances/:id/start` | Instance'ı başlat | ✅ Admin |
| POST | `/api/v1/admin/xray/instances/:id/stop` | Instance'ı durdur | ✅ Admin |
| POST | `/api/v1/admin/xray/instances/:id/restart` | Instance'ı yeniden başlat | ✅ Admin |
| POST | `/api/v1/admin/xray/instances/:id/reload` | Config'i yeniden yükle | ✅ Admin |
| GET | `/api/v1/admin/xray/instances/:id/health` | Instance sağlık kontrolü | ✅ Admin |
| GET | `/api/v1/admin/xray/instances/:id/stats` | Instance istatistikleri | ✅ Admin |

### 📥 Inbound Management (`/api/v1/admin/xray/inbounds`)
| Method | Endpoint | Açıklama | Auth |
|--------|----------|----------|------|
| GET | `/api/v1/admin/xray/inbounds` | Inbound'ları listele | ✅ Admin |
| GET | `/api/v1/admin/xray/inbounds/:id` | Inbound detayı | ✅ Admin |
| POST | `/api/v1/admin/xray/inbounds` | Yeni inbound oluştur | ✅ Admin |
| PUT | `/api/v1/admin/xray/inbounds/:id` | Inbound güncelle | ✅ Admin |
| DELETE | `/api/v1/admin/xray/inbounds/:id` | Inbound sil | ✅ Admin |
| PUT | `/api/v1/admin/xray/inbounds/:id/enable` | Inbound'ı aktifleştir | ✅ Admin |
| PUT | `/api/v1/admin/xray/inbounds/:id/disable` | Inbound'ı devre dışı bırak | ✅ Admin |

### 👥 Client Management (`/api/v1/admin/xray/clients`)
| Method | Endpoint | Açıklama | Auth |
|--------|----------|----------|------|
| GET | `/api/v1/admin/xray/clients` | Client'ları listele | ✅ Admin |
| GET | `/api/v1/admin/xray/clients/:id` | Client detayı | ✅ Admin |
| POST | `/api/v1/admin/xray/clients` | Yeni client oluştur | ✅ Admin |
| DELETE | `/api/v1/admin/xray/clients/:id` | Client sil | ✅ Admin |
| PUT | `/api/v1/admin/xray/clients/:id/enable` | Client'ı aktifleştir | ✅ Admin |
| PUT | `/api/v1/admin/xray/clients/:id/disable` | Client'ı devre dışı bırak | ✅ Admin |
| POST | `/api/v1/admin/xray/clients/:id/regenerate-uuid` | UUID yeniden oluştur | ✅ Admin |
| POST | `/api/v1/admin/xray/clients/:id/reprovision` | Client'ı yeniden provision et | ✅ Admin |

### 📝 Audit Log Management (`/api/v1/admin/audit`)
| Method | Endpoint | Açıklama | Auth |
|--------|----------|----------|------|
| GET | `/api/v1/admin/audit/logs` | Audit logları listele | ✅ Admin |
| GET | `/api/v1/admin/audit/logs/:id` | Audit log detayı | ✅ Admin |
| GET | `/api/v1/admin/audit/stats` | Audit istatistikleri | ✅ Admin |

### ⚙️ System Admin (`/api/v1/admin/system`)
| Method | Endpoint | Açıklama | Auth |
|--------|----------|----------|------|
| GET | `/api/v1/admin/system/health` | Sistem sağlık kontrolü | ✅ Admin |
| GET | `/api/v1/admin/system/stats` | Sistem istatistikleri | ✅ Admin |
| GET | `/api/v1/admin/system/version` | Backend versiyonu | ✅ Admin |
| GET | `/api/v1/admin/system/database` | Database durumu | ✅ Admin |
| GET | `/api/v1/admin/system/xray` | Xray sistem durumu | ✅ Admin |

---

## 🎯 Test Sonuçları Özeti

| # | Test Adımı | Durum | Notlar |
|---|------------|-------|--------|
| 1 | Server Startup | [ ] ✅ [ ] ❌ | |
| 2 | Health Endpoint | [ ] ✅ [ ] ❌ | |
| 3 | Auth Register | [ ] ✅ [ ] ❌ | |
| 4 | Auth Login | [ ] ✅ [ ] ❌ | |
| 5 | JWT Refresh | [ ] ✅ [ ] ❌ | |
| 6 | Admin User | [ ] ✅ [ ] ❌ | |
| 7 | Client Create | [ ] ✅ [ ] ❌ | |
| 8 | Inbound Create | [ ] ✅ [ ] ❌ | |
| 9 | VLESS Config | [ ] ✅ [ ] ❌ | |
| 10 | Client Delete | [ ] ✅ [ ] ❌ | |
| 11 | Config Cleanup | [ ] ✅ [ ] ❌ | |

**Toplam Başarı Oranı:** __% ( __ / 11 )

---

## 📌 Önemli Notlar

### 🔐 Güvenlik
- [ ] Admin token'lar güvenli bir yerde saklanıyor
- [ ] Production'da güçlü şifreler kullanılıyor
- [ ] HTTPS kullanılıyor (production için)

### 🗄️ Database
- [ ] Migration'lar başarıyla çalıştı
- [ ] Foreign key constraint'ler aktif
- [ ] Backup stratejisi hazır

### 🐳 Docker
- [ ] Container'lar health check geçiyor
- [ ] Volume mount'lar doğru
- [ ] Network ayarları doğru

### 📊 Monitoring
- [ ] Metrics endpoint'i erişilebilir
- [ ] Log'lar düzgün yazılıyor
- [ ] Error tracking aktif

---

## 🚨 Kritik Sorunlar

**Bulunan Sorunlar:**
```
1. 
2. 
3. 
```

**Aksiyonlar:**
```
1. 
2. 
3. 
```

---

## ✅ Onay

**Test Eden:** ___________________  
**Tarih:** ___________________  
**İmza:** ___________________  

**Production'a Geçiş Onayı:**
- [ ] Tüm testler başarılı
- [ ] Kritik sorun yok
- [ ] Dokümantasyon güncel
- [ ] Rollback planı hazır

**Onaylayan:** ___________________  
**Tarih:** ___________________
