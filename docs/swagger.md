# Swagger Security Configuration

## Security Definitions

This API uses JWT Bearer token authentication.

### BearerAuth
- **Type**: apiKey
- **Name**: Authorization
- **In**: header
- **Format**: Bearer {token}

## Example Usage

```bash
# Login to get token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'

# Response
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer",
    "expires_in": 900
  }
}

# Use token in protected endpoints
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

## Protected Endpoints

All endpoints under the following paths require authentication:
- `/api/v1/auth/me`
- `/api/v1/auth/sessions`
- `/api/v1/users/*`
- `/api/v1/subscriptions/*`
- `/api/v1/servers/*`
- `/api/v1/nodes/*`
- `/api/v1/xray/*`

## Public Endpoints

The following endpoints do not require authentication:
- `/health`
- `/ready`
- `/api/v1/auth/register`
- `/api/v1/auth/login`
- `/api/v1/auth/refresh`
- `/api/v1/plans/*`

## Admin Endpoints

Admin endpoints require both authentication and admin role:
- Future: `/api/v1/admin/*`
- Future: Plan CRUD operations (POST, PUT, DELETE)
- Future: Server management
- Future: Node management

## Token Expiration

- **Access Token**: 15 minutes
- **Refresh Token**: 7 days (168 hours)

## Error Codes

| Code | Status | Description |
|------|--------|-------------|
| UNAUTHORIZED | 401 | Missing or invalid authentication |
| TOKEN_EXPIRED | 401 | Access token has expired |
| INVALID_TOKEN | 401 | Token is malformed or invalid |
| INVALID_SIGNATURE | 401 | Token signature verification failed |
| FORBIDDEN | 403 | User doesn't have required permissions |
