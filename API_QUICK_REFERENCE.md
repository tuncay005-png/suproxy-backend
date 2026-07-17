# 🚀 API Quick Reference - Client Oluşturma

## 📝 Client Oluşturma Endpoint

### HTTP Method
**POST**

### URL Path
```
/api/v1/admin/xray/clients
```

### Authentication
✅ **Gerekli** - Admin JWT Token

**Header:**
```
Authorization: Bearer {admin-access-token}
```

### Request Body
```json
{
  "inbound_id": "uuid-string",      // Required: Xray inbound ID (UUID format)
  "user_id": "uuid-string",         // Required: User ID (UUID format)
  "email": "client@example.com",    // Required: Valid email format
  "flow": "xtls-rprx-vision"        // Optional: VLESS flow (default: "")
}
```

### Validation Rules
- `inbound_id`: UUID formatında olmalı, inbound mevcut olmalı
- `user_id`: UUID formatında olmalı, user mevcut olmalı
- `email`: Geçerli email formatı olmalı
- `flow`: Boş string veya geçerli flow değeri ("xtls-rprx-vision", "xtls-rprx-direct")

### Success Response (201 Created)
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "inbound_id": "660e8400-e29b-41d4-a716-446655440000",
  "user_id": "770e8400-e29b-41d4-a716-446655440000",
  "uuid": "880e8400-e29b-41d4-a716-446655440000",
  "flow": "xtls-rprx-vision",
  "email": "client@example.com",
  "enabled": true,
  "created_at": "2026-07-16T12:00:00Z",
  "updated_at": "2026-07-16T12:00:00Z"
}
```

### Error Responses

#### 400 Bad Request - Invalid Input
```json
{
  "error": "invalid request body"
}
```

#### 401 Unauthorized - Token Eksik/Geçersiz
```json
{
  "error": "unauthorized"
}
```

#### 403 Forbidden - Admin Yetkisi Yok
```json
{
  "error": "forbidden: admin access required"
}
```

#### 404 Not Found - Inbound Bulunamadı
```json
{
  "error": "inbound not found"
}
```

#### 404 Not Found - User Bulunamadı
```json
{
  "error": "user not found"
}
```

#### 500 Internal Server Error
```json
{
  "error": "internal server error"
}
```

### cURL Örneği
```bash
curl -X POST http://localhost:8080/api/v1/admin/xray/clients \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "inbound_id": "660e8400-e29b-41d4-a716-446655440000",
    "user_id": "770e8400-e29b-41d4-a716-446655440000",
    "email": "newclient@example.com",
    "flow": "xtls-rprx-vision"
  }'
```

### Postman/Insomnia Body
```json
{
  "inbound_id": "{{INBOUND_ID}}",
  "user_id": "{{USER_ID}}",
  "email": "client@example.com",
  "flow": "xtls-rprx-vision"
}
```

---

## 🔄 Diğer Client İşlemleri

### Client Listele
```bash
GET /api/v1/admin/xray/clients
Authorization: Bearer {token}

# Query Parameters (optional):
?offset=0&limit=20&inbound_id={uuid}&enabled=true
```

### Client Detayı
```bash
GET /api/v1/admin/xray/clients/{client-id}
Authorization: Bearer {token}
```

### Client Sil
```bash
DELETE /api/v1/admin/xray/clients/{client-id}
Authorization: Bearer {token}
```

### Client Aktifleştir
```bash
PUT /api/v1/admin/xray/clients/{client-id}/enable
Authorization: Bearer {token}
```

### Client Devre Dışı Bırak
```bash
PUT /api/v1/admin/xray/clients/{client-id}/disable
Authorization: Bearer {token}
```

### Client UUID Yenile
```bash
POST /api/v1/admin/xray/clients/{client-id}/regenerate-uuid
Authorization: Bearer {token}
```

### Client Yeniden Provision Et
```bash
POST /api/v1/admin/xray/clients/{client-id}/reprovision
Authorization: Bearer {token}
```

---

## 📚 Admin Endpoint'leri Hızlı Erişim

### Health & System
- `GET /health` - Public health check
- `GET /ready` - Readiness probe
- `GET /api/v1/admin/health` - Admin health (requires auth)
- `GET /api/v1/admin/system/health` - System health details
- `GET /api/v1/admin/system/stats` - System statistics
- `GET /api/v1/admin/system/version` - API version
- `GET /api/v1/admin/system/database` - Database status
- `GET /api/v1/admin/system/xray` - Xray system status

### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - Logout
- `GET /api/v1/auth/me` - Get current user

### User Management (Admin)
- `GET /api/v1/admin/users` - List users
- `GET /api/v1/admin/users/:id` - Get user
- `PUT /api/v1/admin/users/:id/status` - Update status
- `PUT /api/v1/admin/users/:id/role` - Update role

### Xray Instance Management (Admin)
- `GET /api/v1/admin/xray/instances` - List instances
- `GET /api/v1/admin/xray/instances/:id` - Get instance
- `POST /api/v1/admin/xray/instances/:id/start` - Start instance
- `POST /api/v1/admin/xray/instances/:id/stop` - Stop instance
- `POST /api/v1/admin/xray/instances/:id/restart` - Restart instance
- `POST /api/v1/admin/xray/instances/:id/reload` - Reload config
- `GET /api/v1/admin/xray/instances/:id/health` - Instance health
- `GET /api/v1/admin/xray/instances/:id/stats` - Instance stats

### Inbound Management (Admin)
- `GET /api/v1/admin/xray/inbounds` - List inbounds
- `GET /api/v1/admin/xray/inbounds/:id` - Get inbound
- `POST /api/v1/admin/xray/inbounds` - Create inbound
- `PUT /api/v1/admin/xray/inbounds/:id` - Update inbound
- `DELETE /api/v1/admin/xray/inbounds/:id` - Delete inbound
- `PUT /api/v1/admin/xray/inbounds/:id/enable` - Enable inbound
- `PUT /api/v1/admin/xray/inbounds/:id/disable` - Disable inbound

### Audit Logs (Admin)
- `GET /api/v1/admin/audit/logs` - List audit logs
- `GET /api/v1/admin/audit/logs/:id` - Get audit log
- `GET /api/v1/admin/audit/stats` - Audit statistics

---

## 🎯 Hızlı Test Senaryosu

### 1. Admin Token Al
```bash
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"AdminPass123!"}'

# Token'ı kaydet
export ADMIN_TOKEN="eyJhbGc..."
```

### 2. Prerequisite ID'leri Al
```bash
# Instance ID
curl -X GET http://localhost:8080/api/v1/admin/xray/instances \
  -H "Authorization: Bearer $ADMIN_TOKEN"

export INSTANCE_ID="instance-uuid"

# Inbound ID
curl -X GET "http://localhost:8080/api/v1/admin/xray/inbounds?instance_id=$INSTANCE_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN"

export INBOUND_ID="inbound-uuid"

# User ID
curl -X GET http://localhost:8080/api/v1/admin/users \
  -H "Authorization: Bearer $ADMIN_TOKEN"

export USER_ID="user-uuid"
```

### 3. Client Oluştur
```bash
curl -X POST http://localhost:8080/api/v1/admin/xray/clients \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"inbound_id\": \"$INBOUND_ID\",
    \"user_id\": \"$USER_ID\",
    \"email\": \"testclient@example.com\",
    \"flow\": \"xtls-rprx-vision\"
  }"
```

### 4. Client'ı Doğrula
```bash
export CLIENT_ID="response-client-id"

curl -X GET http://localhost:8080/api/v1/admin/xray/clients/$CLIENT_ID \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

---

## 💡 Tips & Best Practices

### JWT Token Management
- Access token süresi: 15 dakika
- Refresh token süresi: 7 gün
- Token'ları güvenli sakla (environment variable, secrets manager)
- Token expired olduğunda `/auth/refresh` kullan

### Error Handling
- Her zaman HTTP status code'u kontrol et
- 4xx hatalar: Client hatası, request'i düzelt
- 5xx hatalar: Server hatası, retry logic uygula
- Rate limiting için 429 status code'una dikkat et

### Production Checklist
- [ ] HTTPS kullan
- [ ] Token'ları loglamayın
- [ ] API rate limiting ayarla
- [ ] CORS ayarlarını kontrol et
- [ ] Request timeout'ları ayarla
- [ ] Retry logic implement et
- [ ] Circuit breaker pattern kullan

---

**Not:** Bu dokümantasyon development/staging environment içindir. Production URL ve credential'lar için DevOps ekibiyle iletişime geçin.
