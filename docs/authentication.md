# Authentication & Authorization

## Overview

SuProxy Backend uses JWT (JSON Web Token) based authentication with role-based access control (RBAC).

## Authentication Flow

### 1. Registration
```
POST /api/v1/auth/register
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}

Response:
{
  "success": true,
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 2. Login
```
POST /api/v1/auth/login
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}

Response:
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer",
    "expires_in": 900
  }
}
```

### 3. Access Protected Resources
```
GET /api/v1/users/me
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...

Response:
{
  "success": true,
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "role": "user"
  }
}
```

### 4. Refresh Token
```
POST /api/v1/auth/refresh
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}

Response:
{
  "success": true,
  "data": {
    "access_token": "new_access_token",
    "refresh_token": "new_refresh_token",
    "token_type": "Bearer",
    "expires_in": 900
  }
}
```

## Token Details

### Access Token
- **Expiry**: 15 minutes
- **Type**: JWT
- **Claims**: user_id, email, role, token_type
- **Usage**: All protected API requests

### Refresh Token
- **Expiry**: 7 days (168 hours)
- **Type**: JWT
- **Storage**: Database (SHA256 hash)
- **Rotation**: New token on each refresh
- **Usage**: Get new access token

## Authorization Levels

### Public Routes
No authentication required:
- `GET /health`
- `GET /ready`
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `GET /api/v1/plans`
- `GET /api/v1/plans/:id`

### User Routes
Requires authentication (any valid user):
- `GET /api/v1/auth/me`
- `GET /api/v1/auth/sessions`
- `DELETE /api/v1/auth/sessions/:id`
- `POST /api/v1/auth/logout-all`
- `GET /api/v1/users/me`
- `PUT /api/v1/users/me`
- `GET /api/v1/subscriptions/me`
- `GET /api/v1/servers`
- `GET /api/v1/nodes`
- `GET /api/v1/xray/instances`

### Admin Routes
Requires authentication + admin role:
- (Future) Plan management (POST, PUT, DELETE)
- (Future) Server management
- (Future) Node management
- (Future) User management
- (Future) Analytics dashboard

## Roles

### User (default)
- Access to own profile
- View subscription
- View servers and nodes
- View Xray instances

### Admin
- All user permissions
- Create/Update/Delete plans
- Manage servers and nodes
- Manage Xray instances
- View all users
- Access analytics

## Security Features

### 1. Token Rotation
- Refresh tokens are rotated on each use
- Old refresh token is revoked
- Prevents token replay attacks

### 2. Session Management
- Multiple active sessions per user
- Device tracking (IP, User-Agent)
- Individual session logout
- Logout all sessions

### 3. Account Security
- Failed login counter
- Account locking (5 failed attempts)
- Lock duration: 15 minutes
- Audit log for security events

### 4. Password Security
- Bcrypt hashing (cost 12)
- Minimum 8 characters
- Complexity requirements

## Error Handling

### Authentication Errors

#### 401 Unauthorized
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "authentication required"
  }
}
```

#### 401 Token Expired
```json
{
  "success": false,
  "error": {
    "code": "TOKEN_EXPIRED",
    "message": "access token has expired"
  }
}
```

#### 401 Invalid Token
```json
{
  "success": false,
  "error": {
    "code": "INVALID_TOKEN",
    "message": "invalid access token"
  }
}
```

### Authorization Errors

#### 403 Forbidden
```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "admin access required"
  }
}
```

## Best Practices

### Client Implementation

1. **Store tokens securely**
   - Use httpOnly cookies or secure storage
   - Never store in localStorage for production

2. **Handle token expiration**
   - Implement automatic token refresh
   - Redirect to login on refresh failure

3. **Set Authorization header**
   ```javascript
   axios.defaults.headers.common['Authorization'] = `Bearer ${accessToken}`;
   ```

4. **Logout properly**
   - Call logout endpoint
   - Clear local tokens
   - Redirect to login

### Server Implementation

1. **Use HTTPS only** (production)
2. **Rotate secrets regularly**
3. **Monitor failed login attempts**
4. **Implement rate limiting**
5. **Log security events**

## Middleware Stack

```
Request
  ↓
CORS Middleware
  ↓
Error Handler
  ↓
Request Logger
  ↓
Auth Middleware (if protected route)
  ↓
  ├─ Validate Bearer token
  ├─ Verify signature
  ├─ Check expiration
  └─ Set user context
  ↓
Admin Middleware (if admin route)
  ↓
  ├─ Check role = admin
  └─ Return 403 if not admin
  ↓
Handler
```

## Testing

### Get Token
```bash
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}' \
  | jq -r '.data.access_token')
```

### Use Token
```bash
curl http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN"
```

### Test Token Expiry
```bash
# Wait 15+ minutes
curl http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN"
# Should return TOKEN_EXPIRED error
```

### Refresh Token
```bash
REFRESH_TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}' \
  | jq -r '.data.refresh_token')

curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\":\"$REFRESH_TOKEN\"}"
```
