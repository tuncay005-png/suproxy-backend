# Phase 2 Security Hardening - Verification Guide

## Overview
This guide provides comprehensive step-by-step manual testing procedures to verify the Phase 2 security hardening implementation, including token family tracking, reuse detection, family-wide revocation, rate limiting, and frontend multi-tab synchronization.

---

## Table of Contents
1. [Security Fix Verification](#1-security-fix-verification)
2. [Prerequisites](#2-prerequisites)
3. [Migration Schema Verification](#3-migration-schema-verification)
4. [Normal Token Refresh Test](#4-normal-token-refresh-test)
5. [Token Reuse Attack Scenario Test](#5-token-reuse-attack-scenario-test)
6. [Family Revocation Test](#6-family-revocation-test)
7. [Rate Limiting Test](#7-rate-limiting-test)
8. [Multi-Tab Sync Browser Test](#8-multi-tab-sync-browser-test)
9. [Backend Log Verification](#9-backend-log-verification)
10. [Compilation & Type Check](#10-compilation--type-check)
11. [Security Fix Validation](#11-security-fix-validation)
12. [Rollback Procedures](#12-rollback-procedures)

---

## 1. Security Fix Verification

### Critical Security Issue Fixed

#### ❌ OLD BEHAVIOR (VULNERABLE)
**File:** `internal/application/usecase/auth/refresh_token_command.go` (Before Fix)

```go
// VULNERABLE CODE - 5 MINUTE LIMITATION
if storedToken.IsRevoked {
    timeSinceRevoked := time.Since(*storedToken.RevokedAt)
    
    // ❌ BUG: Only detects reuse within 5 minutes
    if timeSinceRevoked < 5*time.Minute {
        // Security alert and family revocation
        c.logger.Error("SECURITY ALERT: Recent token reuse detected")
        c.refreshTokenRepo.RevokeByFamilyID(ctx, storedToken.FamilyID)
    } else {
        // ⚠️ VULNERABILITY: After 5 minutes, reuse is allowed!
        c.logger.Warn("Old revoked token reuse (>5min, ignoring)")
    }
    
    return nil, jwt.ErrInvalidToken
}
```

**Attack Scenario:**
```
Time 10:00 - Token revoked during normal rotation
Time 10:04 - Attacker reuses token → ❌ Detected, family revoked
Time 10:06 - Attacker waits and reuses → ✅ NOT detected, attack succeeds!
```

#### ✅ NEW BEHAVIOR (SECURE)
**File:** `internal/application/usecase/auth/refresh_token_command.go` (After Fix)

```go
// SECURE CODE - NO TIME LIMITATION
if storedToken.IsRevoked {
    var timeSinceRevoked time.Duration
    if storedToken.RevokedAt != nil {
        timeSinceRevoked = time.Since(*storedToken.RevokedAt)
    }
    
    // ✅ FIX: ALWAYS log security alert for ANY revoked token reuse
    c.logger.Error("SECURITY ALERT: Revoked token reuse detected - possible token theft",
        "token_id", storedToken.ID,
        "family_id", storedToken.FamilyID,
        "user_id", storedToken.UserID,
        "revoked_at", storedToken.RevokedAt,
        "time_since_revoked_seconds", int(timeSinceRevoked.Seconds()))

    // ✅ FIX: ALWAYS revoke entire family on ANY revoked token reuse
    if err := c.refreshTokenRepo.RevokeByFamilyID(ctx, storedToken.FamilyID); err != nil {
        c.logger.Error("Failed to revoke token family", "error", err, "family_id", storedToken.FamilyID)
    } else {
        c.logger.Warn("Token family revoked due to reuse detection",
            "family_id", storedToken.FamilyID,
            "user_id", storedToken.UserID)
    }

    // ✅ FIX: Create security incident audit log for ANY reuse
    securityLog := audit.NewLog(
        storedToken.UserID,
        "security.incident",
        "token_reuse_detected",
        storedToken.ID,
        storedToken.IPAddress,
        storedToken.UserAgent,
    )
    securityLog.AddMetadata("family_id", storedToken.FamilyID.String())
    securityLog.AddMetadata("time_since_revoked_seconds", int(timeSinceRevoked.Seconds()))
    securityLog.AddMetadata("severity", "high")
    
    if err := c.auditRepo.Create(ctx, securityLog); err != nil {
        c.logger.Warn("Failed to create security incident audit log", "error", err)
    }
    
    // Always return invalid token
    return nil, jwt.ErrInvalidToken
}
```

**Attack Prevention:**
```
Time 10:00 - Token revoked during normal rotation
Time 10:04 - Attacker reuses token → ❌ Detected, family revoked
Time 10:06 - Attacker reuses token → ❌ Detected, family revoked
Time 11:00 - Attacker reuses token → ❌ Detected, family revoked
Next day  - Attacker reuses token → ❌ Detected, family revoked
```

### Key Differences

| Aspect | OLD (Vulnerable) | NEW (Secure) |
|--------|------------------|--------------|
| Detection Window | 5 minutes only | ✅ Forever (no time limit) |
| After 5 Minutes | ⚠️ Reuse allowed | ✅ Always blocked |
| Security Logging | Only < 5 min | ✅ All reuses logged |
| Family Revocation | Only < 5 min | ✅ Always revoked |
| Audit Trail | Incomplete | ✅ Complete forensics |

---

## 2. Prerequisites

### Database Access
Ensure PostgreSQL is running and accessible:
```powershell
# Test database connection
psql -U suproxy -d suproxy -h localhost -c "SELECT version();"
```

**Expected Output:**
```
PostgreSQL 14.x or higher
```

### Backend Service
```powershell
cd c:\Users\Tuncay\Desktop\suproxy-backend

# Build
go build -o bin\api.exe ./cmd/api

# Run with logging
.\bin\api.exe > backend.log 2>&1
```

### Frontend Service
```powershell
cd c:\Users\Tuncay\Desktop\suproxy-admin

# Install dependencies
npm install

# Start dev server
npm run dev
```

### Test User Accounts
Create test users in database:
```sql
-- Admin user
INSERT INTO users (id, email, password_hash, role, is_active) 
VALUES (
    gen_random_uuid(),
    'admin@suproxy.com',
    '$2a$10$yourhashedpassword', -- Hash for 'Admin123!'
    'admin',
    true
);

-- Test users for attack scenarios
INSERT INTO users (id, email, password_hash, role, is_active) 
VALUES 
(gen_random_uuid(), 'test@suproxy.com', '$2a$10$hash', 'user', true),
(gen_random_uuid(), 'user@suproxy.com', '$2a$10$hash', 'user', true),
(gen_random_uuid(), 'ratelimit@suproxy.com', '$2a$10$hash', 'user', true);
```

---

## 3. Migration Schema Verification

### Step 1: Check Current Schema (Before Migration)

```powershell
# Check if migration already applied
psql -U suproxy -d suproxy -h localhost -c "SELECT column_name FROM information_schema.columns WHERE table_name = 'refresh_tokens' AND column_name = 'family_id';"
```

**Expected Output (if not applied):**
```
(0 rows)
```

### Step 2: View Current Migration Version

```powershell
psql -U suproxy -d suproxy -h localhost -c "SELECT version, dirty FROM schema_migrations ORDER BY version DESC LIMIT 1;"
```

**Expected Output:**
```
 version | dirty 
---------+-------
       9 | f
```

### Step 3: View Current refresh_tokens Schema

```powershell
psql -U suproxy -d suproxy -h localhost -c "\d refresh_tokens"
```

**Expected Output (before migration):**
```
                       Table "public.refresh_tokens"
     Column     |            Type             | Nullable | Default 
----------------+-----------------------------+----------+---------
 id             | uuid                        | not null | 
 user_id        | uuid                        | not null | 
 token_hash     | text                        | not null | 
 expires_at     | timestamp without time zone | not null | 
 is_revoked     | boolean                     | not null | false
 revoked_at     | timestamp without time zone |          | 
 ip_address     | text                        |          | 
 user_agent     | text                        |          | 
 created_at     | timestamp without time zone | not null | 
```

### Step 4: Backup Database (IMPORTANT!)

```powershell
# Create backup
$timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
pg_dump -U suproxy -d suproxy -h localhost -f "backup_before_phase2_$timestamp.sql"

# Verify backup created
Get-Item "backup_before_phase2_*.sql"
```

### Step 5: Apply Migration

```powershell
# Navigate to backend directory
cd c:\Users\Tuncay\Desktop\suproxy-backend

# Apply migration using psql
psql -U suproxy -d suproxy -h localhost -f migrations\000010_add_token_family.up.sql
```

**Expected Output:**
```
ALTER TABLE
CREATE INDEX
CREATE INDEX
UPDATE 42  -- (number of existing tokens backfilled)
ALTER TABLE
COMMENT
```

### Step 6: Verify family_id Column Created

```powershell
psql -U suproxy -d suproxy -h localhost -c "SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_name = 'refresh_tokens' AND column_name = 'family_id';"
```

**Expected Output:**
```
 column_name | data_type | is_nullable 
-------------+-----------+-------------
 family_id   | uuid      | NO
```

### Step 7: Verify Indexes Created

```powershell
psql -U suproxy -d suproxy -h localhost -c "SELECT indexname, indexdef FROM pg_indexes WHERE tablename = 'refresh_tokens' AND indexname LIKE '%family%';"
```

**Expected Output:**
```
            indexname             |                                    indexdef                                    
----------------------------------+-------------------------------------------------------------------------------
 idx_refresh_tokens_family_id     | CREATE INDEX idx_refresh_tokens_family_id ON refresh_tokens USING btree (family_id)
 idx_refresh_tokens_revoked_family| CREATE INDEX idx_refresh_tokens_revoked_family ON refresh_tokens USING btree (is_revoked, family_id, revoked_at) WHERE (is_revoked = true)
```

### Step 8: Verify Backfill Worked (No NULL family_id)

```powershell
psql -U suproxy -d suproxy -h localhost -c "SELECT COUNT(*) as null_count FROM refresh_tokens WHERE family_id IS NULL;"
```

**Expected Output:**
```
 null_count 
------------
          0
```

### Step 9: Verify Backfill Logic (Existing Tokens Are Own Family)

```powershell
psql -U suproxy -d suproxy -h localhost -c "SELECT id, family_id, (id = family_id) as is_family_root FROM refresh_tokens LIMIT 5;"
```

**Expected Output:**
```
                  id                  |              family_id              | is_family_root 
--------------------------------------+-------------------------------------+----------------
 a1b2c3d4-e5f6-7890-abcd-ef1234567890 | a1b2c3d4-e5f6-7890-abcd-ef1234567890 | t
 b2c3d4e5-f6a7-8901-bcde-f12345678901 | b2c3d4e5-f6a7-8901-bcde-f12345678901 | t
...
```

### Step 10: View Updated Schema

```powershell
psql -U suproxy -d suproxy -h localhost -c "\d refresh_tokens"
```

**Expected Output (after migration):**
```
                       Table "public.refresh_tokens"
     Column     |            Type             | Nullable | Default 
----------------+-----------------------------+----------+---------
 id             | uuid                        | not null | 
 user_id        | uuid                        | not null | 
 token_hash     | text                        | not null | 
 expires_at     | timestamp without time zone | not null | 
 is_revoked     | boolean                     | not null | false
 revoked_at     | timestamp without time zone |          | 
 ip_address     | text                        |          | 
 user_agent     | text                        |          | 
 created_at     | timestamp without time zone | not null | 
 family_id      | uuid                        | not null |   ← NEW!
```

---

## 4. Normal Token Refresh Test

This test verifies that normal token rotation works correctly with family tracking.

### Step 1: Login and Get Initial Tokens

```powershell
# Login
curl.exe -X POST http://localhost:8080/api/v1/auth/login `
  -H "Content-Type: application/json" `
  -d '{\"email\": \"admin@suproxy.com\", \"password\": \"Admin123!\"}' `
  -c cookies.txt `
  -v
```

**Expected Response:**
```json
{
  "user": {
    "id": "uuid",
    "email": "admin@suproxy.com",
    "role": "admin"
  }
}
```

**Check cookies.txt:**
```powershell
Get-Content cookies.txt
```

**Expected Output:**
```
localhost   FALSE   /   FALSE   0   access_token    eyJhbGc...
localhost   FALSE   /   FALSE   0   refresh_token   eyJhbGc...
```

### Step 2: Extract User ID for Queries

```powershell
psql -U suproxy -d suproxy -h localhost -c "SELECT id FROM users WHERE email = 'admin@suproxy.com';"
```

**Save this UUID for later queries.**

### Step 3: Check Initial Token in Database

```powershell
# Replace USER_ID with actual UUID from Step 2
psql -U suproxy -d suproxy -h localhost -c "SELECT id, family_id, (id = family_id) as is_root, is_revoked, created_at FROM refresh_tokens WHERE user_id = 'USER_ID' ORDER BY created_at DESC LIMIT 1;"
```

**Expected Output:**
```
                  id                  |              family_id              | is_root | is_revoked |         created_at         
--------------------------------------+-------------------------------------+---------+------------+----------------------------
 new-token-uuid                       | new-token-uuid                      | t       | f          | 2025-01-15 10:00:00.123456
```

**✅ Verify:** Initial token is its own family root (`id = family_id`, `is_root = true`)

### Step 4: Use Refresh Token (First Rotation)

```powershell
curl.exe -X POST http://localhost:8080/api/v1/auth/refresh `
  -H "Content-Type: application/json" `
  -b cookies.txt `
  -c cookies2.txt `
  -v
```

**Expected Response:**
```json
{
  "message": "Token refreshed successfully"
}
```

### Step 5: Verify Token Rotation in Database

```powershell
psql -U suproxy -d suproxy -h localhost -c "SELECT id, family_id, (id = family_id) as is_root, is_revoked, created_at FROM refresh_tokens WHERE user_id = 'USER_ID' ORDER BY created_at DESC LIMIT 2;"
```

**Expected Output:**
```
                  id                  |              family_id              | is_root | is_revoked |         created_at         
--------------------------------------+-------------------------------------+---------+------------+----------------------------
 token-2-uuid                         | token-1-uuid                        | f       | f          | 2025-01-15 10:01:00.123456  ← NEW (active)
 token-1-uuid                         | token-1-uuid                        | t       | t          | 2025-01-15 10:00:00.123456  ← OLD (revoked)
```

**✅ Verify:**
- Old token (Token 1): `is_revoked = true`
- New token (Token 2): `is_revoked = false`
- Both tokens share same `family_id` (Token 1's ID)
- New token is NOT root (`is_root = false`)

### Step 6: Multiple Rotations

```powershell
# Rotation 2
curl.exe -X POST http://localhost:8080/api/v1/auth/refresh -b cookies2.txt -c cookies3.txt

# Rotation 3
curl.exe -X POST http://localhost:8080/api/v1/auth/refresh -b cookies3.txt -c cookies4.txt
```

### Step 7: Verify Family Chain

```powershell
psql -U suproxy -d suproxy -h localhost -c "SELECT id, family_id, is_revoked, created_at FROM refresh_tokens WHERE user_id = 'USER_ID' ORDER BY created_at;"
```

**Expected Output:**
```
                  id                  |              family_id              | is_revoked |         created_at         
--------------------------------------+-------------------------------------+------------+----------------------------
 token-1-uuid (ROOT)                  | token-1-uuid                        | t          | 2025-01-15 10:00:00
 token-2-uuid                         | token-1-uuid                        | t          | 2025-01-15 10:01:00
 token-3-uuid                         | token-1-uuid                        | t          | 2025-01-15 10:02:00
 token-4-uuid                         | token-1-uuid                        | f          | 2025-01-15 10:03:00  ← ACTIVE
```

**✅ Verify:**
- All 4 tokens share the same `family_id`
- First 3 tokens revoked (normal rotation)
- Last token active

---

## 5. Token Reuse Attack Scenario Test

This test simulates an attacker reusing a stolen revoked token to verify the security fix.

### Step 1: Create Fresh Test User Session

```powershell
# Login as test user
curl.exe -X POST http://localhost:8080/api/v1/auth/login `
  -H "Content-Type: application/json" `
  -d '{\"email\": \"test@suproxy.com\", \"password\": \"Test123!\"}' `
  -c attack_cookies1.txt
```

### Step 2: Extract Refresh Token (Attacker Captures Token)

```powershell
# View cookies
Get-Content attack_cookies1.txt | Select-String "refresh_token"
```

**Save this cookie file - it represents the "stolen" token.**

### Step 3: Legitimate User Refreshes (Token Gets Revoked)

```powershell
# User performs normal refresh
curl.exe -X POST http://localhost:8080/api/v1/auth/refresh `
  -b attack_cookies1.txt `
  -c attack_cookies2.txt
```

**Expected Response:**
```json
{
  "message": "Token refreshed successfully"
}
```

### Step 4: Verify Token Revoked in Database

```powershell
# Get test user ID
psql -U suproxy -d suproxy -h localhost -c "SELECT id FROM users WHERE email = 'test@suproxy.com';"

# Check tokens
psql -U suproxy -d suproxy -h localhost -c "SELECT id, family_id, is_revoked, revoked_at FROM refresh_tokens WHERE user_id = 'TEST_USER_ID' ORDER BY created_at DESC LIMIT 2;"
```

**Expected Output:**
```
                  id                  |              family_id              | is_revoked |         revoked_at         
--------------------------------------+-------------------------------------+------------+----------------------------
 token-2-uuid                         | token-1-uuid                        | f          | NULL                        ← ACTIVE
 token-1-uuid                         | token-1-uuid                        | t          | 2025-01-15 10:01:00.123456  ← REVOKED
```

### Step 5: ATTACK - Reuse Revoked Token (Immediately)

```powershell
# Attacker tries to use stolen token (immediately after revocation)
curl.exe -X POST http://localhost:8080/api/v1/auth/refresh `
  -b attack_cookies1.txt `
  -v
```

**Expected Response:**
```
HTTP/1.1 401 Unauthorized
Content-Type: application/json

{
  "error": {
    "code": "INVALID_TOKEN",
    "message": "Invalid or expired token"
  }
}
```

**✅ Verify:** Request rejected (no details revealed to attacker)

### Step 6: Check Backend Logs for Security Alert

```powershell
# View recent logs
Get-Content backend.log | Select-String "SECURITY ALERT" -Context 0,3 | Select-Object -Last 1
```

**Expected Log Output:**
```
ERROR SECURITY ALERT: Revoked token reuse detected - possible token theft
  token_id=token-1-uuid
  family_id=token-1-uuid
  user_id=TEST_USER_ID
  revoked_at=2025-01-15T10:01:00Z
  time_since_revoked_seconds=5
```

### Step 7: Verify Family Revocation Occurred

```powershell
psql -U suproxy -d suproxy -h localhost -c "SELECT id, family_id, is_revoked, revoked_at FROM refresh_tokens WHERE user_id = 'TEST_USER_ID' ORDER BY created_at;"
```

**Expected Output:**
```
                  id                  |              family_id              | is_revoked |         revoked_at         
--------------------------------------+-------------------------------------+------------+----------------------------
 token-1-uuid                         | token-1-uuid                        | t          | 2025-01-15 10:01:00  ← Original revocation
 token-2-uuid                         | token-1-uuid                        | t          | 2025-01-15 10:01:05  ← Family revoked!
```

**✅ Verify:** 
- Both Token 1 AND Token 2 are now revoked
- Token 2's `revoked_at` is AFTER the reuse attempt

### Step 8: Verify Security Incident Audit Log

```powershell
psql -U suproxy -d suproxy -h localhost -c "SELECT id, user_id, action, resource_type, metadata, created_at FROM audit_logs WHERE action = 'security.incident' ORDER BY created_at DESC LIMIT 1;"
```

**Expected Output:**
```
                  id                  |              user_id                |      action       |   resource_type      |                    metadata                     |         created_at         
--------------------------------------+-------------------------------------+-------------------+----------------------+-------------------------------------------------+----------------------------
 log-uuid                             | TEST_USER_ID                        | security.incident | token_reuse_detected | {"family_id": "...", "severity": "high", ...}   | 2025-01-15 10:01:05
```

**✅ Verify:**
- `action = 'security.incident'`
- `resource_type = 'token_reuse_detected'`
- `metadata` contains `family_id` and `severity: high`

### Step 9: CRITICAL TEST - Attack After 5 Minutes (Verifies Fix)

**This test verifies that the 5-minute limitation bug is fixed.**

```powershell
# Wait 6 minutes (or manually adjust revoked_at in database to simulate)
# Update revoked_at to 6 minutes ago
psql -U suproxy -d suproxy -h localhost -c "UPDATE refresh_tokens SET revoked_at = NOW() - INTERVAL '6 minutes', is_revoked = true WHERE id = 'token-1-uuid';"

# Create new active token for the family
psql -U suproxy -d suproxy -h localhost -c "UPDATE refresh_tokens SET is_revoked = false, revoked_at = NULL WHERE id = 'token-2-uuid';"

# ATTACK: Reuse old token after 6 minutes
curl.exe -X POST http://localhost:8080/api/v1/auth/refresh `
  -b attack_cookies1.txt `
  -v
```

**Expected Response:**
```
HTTP/1.1 401 Unauthorized

{
  "error": {
    "code": "INVALID_TOKEN",
    "message": "Invalid or expired token"
  }
}
```

**✅ CRITICAL VERIFICATION:**
- Request MUST be rejected (not allowed after 5 minutes)
- Check logs for "SECURITY ALERT" with `time_since_revoked_seconds > 300`
- Verify family revoked AGAIN in database

```powershell
# Verify family revoked
psql -U suproxy -d suproxy -h localhost -c "SELECT is_revoked FROM refresh_tokens WHERE family_id = 'token-1-uuid';"
```

**Expected Output:**
```
 is_revoked 
------------
 t          ← ALL tokens revoked
 t
```

**✅ This confirms the 5-minute limitation bug is FIXED!**

---

## 6. Family Revocation Test

This test verifies that reusing ANY token in a family revokes ALL tokens in that family.

### Step 1: Create Token Family with Multiple Rotations

```powershell
# Login
curl.exe -X POST http://localhost:8080/api/v1/auth/login `
  -H "Content-Type: application/json" `
  -d '{\"email\": \"user@suproxy.com\", \"password\": \"User123!\"}' `
  -c family_cookies_r1.txt

# Rotation 1
curl.exe -X POST http://localhost:8080/api/v1/auth/refresh `
  -b family_cookies_r1.txt -c family_cookies_r2.txt

# Rotation 2
curl.exe -X POST http://localhost:8080/api/v1/auth/refresh `
  -b family_cookies_r2.txt -c family_cookies_r3.txt

# Rotation 3
curl.exe -X POST http://localhost:8080/api/v1/auth/refresh `
  -b family_cookies_r3.txt -c family_cookies_r4.txt

# Rotation 4
curl.exe -X POST http://localhost:8080/api/v1/auth/refresh `
  -b family_cookies_r4.txt -c family_cookies_r5.txt
```

### Step 2: Verify Family Chain Created

```powershell
# Get user ID
psql -U suproxy -d suproxy -h localhost -c "SELECT id FROM users WHERE email = 'user@suproxy.com';"

# View family
psql -U suproxy -d suproxy -h localhost -c "SELECT id, family_id, is_revoked, created_at FROM refresh_tokens WHERE user_id = 'USER_ID' ORDER BY created_at;"
```

**Expected Output:**
```
                  id                  |              family_id              | is_revoked |         created_at         
--------------------------------------+-------------------------------------+------------+----------------------------
 token-1-uuid (ROOT)                  | token-1-uuid                        | t          | 10:00:00
 token-2-uuid                         | token-1-uuid                        | t          | 10:00:15
 token-3-uuid                         | token-1-uuid                        | t          | 10:00:30
 token-4-uuid                         | token-1-uuid                        | t          | 10:00:45
 token-5-uuid                         | token-1-uuid                        | f          | 10:01:00  ← ONLY ACTIVE TOKEN
```

**✅ Verify:**
- 5 tokens in same family
- First 4 revoked (normal rotation)
- Last token (Token 5) is active

### Step 3: Count Active vs Revoked

```powershell
psql -U suproxy -d suproxy -h localhost -c "SELECT family_id, COUNT(*) as total, SUM(CASE WHEN is_revoked THEN 1 ELSE 0 END) as revoked, SUM(CASE WHEN NOT is_revoked THEN 1 ELSE 0 END) as active FROM refresh_tokens WHERE user_id = 'USER_ID' GROUP BY family_id;"
```

**Expected Output:**
```
              family_id              | total | revoked | active 
-------------------------------------+-------+---------+--------
 token-1-uuid                        |     5 |       4 |      1
```

### Step 4: ATTACK - Reuse Rotation 1 Token

```powershell
# Attacker reuses Token 2 (first rotation, revoked 4 rotations ago)
curl.exe -X POST http://localhost:8080/api/v1/auth/refresh `
  -b family_cookies_r2.txt `
  -v
```

**Expected Response:**
```
HTTP/1.1 401 Unauthorized
```

### Step 5: Verify ENTIRE Family Revoked

```powershell
psql -U suproxy -d suproxy -h localhost -c "SELECT id, family_id, is_revoked, revoked_at FROM refresh_tokens WHERE user_id = 'USER_ID' ORDER BY created_at;"
```

**Expected Output:**
```
                  id                  |              family_id              | is_revoked |         revoked_at         
--------------------------------------+-------------------------------------+------------+----------------------------
 token-1-uuid                         | token-1-uuid                        | t          | 10:00:15
 token-2-uuid                         | token-1-uuid                        | t          | 10:00:30  ← Reused token
 token-3-uuid                         | token-1-uuid                        | t          | 10:00:45
 token-4-uuid                         | token-1-uuid                        | t          | 10:01:00
 token-5-uuid                         | token-1-uuid                        | t          | 10:01:15  ← WAS ACTIVE, NOW REVOKED!
```

**✅ CRITICAL VERIFICATION:**
- ALL 5 tokens are now revoked (`is_revoked = true`)
- Token 5 (which was active) is now revoked
- Token 5's `revoked_at` matches the reuse attack time

### Step 6: Verify Legitimate User Session Terminated

```powershell
# Legitimate user (with Token 5) tries to refresh
curl.exe -X POST http://localhost:8080/api/v1/auth/refresh `
  -b family_cookies_r5.txt `
  -v
```

**Expected Response:**
```
HTTP/1.1 401 Unauthorized

{
  "error": {
    "code": "INVALID_TOKEN",
    "message": "Invalid or expired token"
  }
}
```

**✅ Verify:** Even legitimate user must re-authenticate (security over convenience)

### Step 7: User Must Re-authenticate

```powershell
# User logs in again (creates new family)
curl.exe -X POST http://localhost:8080/api/v1/auth/login `
  -H "Content-Type: application/json" `
  -d '{\"email\": \"user@suproxy.com\", \"password\": \"User123!\"}' `
  -c family_cookies_new.txt
```

### Step 8: Verify New Family Created

```powershell
psql -U suproxy -d suproxy -h localhost -c "SELECT family_id, COUNT(*) as total, SUM(CASE WHEN is_revoked THEN 1 ELSE 0 END) as revoked, SUM(CASE WHEN NOT is_revoked THEN 1 ELSE 0 END) as active FROM refresh_tokens WHERE user_id = 'USER_ID' GROUP BY family_id ORDER BY MIN(created_at);"
```

**Expected Output:**
```
              family_id              | total | revoked | active 
-------------------------------------+-------+---------+--------
 token-1-uuid (OLD FAMILY)           |     5 |       5 |      0  ← Compromised family
 token-6-uuid (NEW FAMILY)           |     1 |       0 |      1  ← Fresh start
```

**✅ Verify:** New family isolated from compromised family

---

## 7. Rate Limiting Test

This test verifies that the `/api/v1/auth/refresh` endpoint is protected against brute force attacks.

### Step 1: Login to Get Valid Refresh Token

```powershell
curl.exe -X POST http://localhost:8080/api/v1/auth/login `
  -H "Content-Type: application/json" `
  -d '{\"email\": \"ratelimit@suproxy.com\", \"password\": \"Test123!\"}' `
  -c ratelimit_cookies.txt
```

### Step 2: Manual Rate Limit Test (Sequential)

```powershell
# Make 11 sequential requests
1..11 | ForEach-Object {
    $response = try {
        curl.exe -X POST http://localhost:8080/api/v1/auth/refresh `
          -b ratelimit_cookies.txt `
          -c "ratelimit_cookies_$_.txt" `
          -w "%{http_code}" `
          -s -o response_$_.json
        
        $statusCode = $response[-3..-1] -join ''
        Write-Host "Request $_`: HTTP $statusCode"
        
    } catch {
        Write-Host "Request $_`: FAILED"
    }
    
    Start-Sleep -Milliseconds 100
}
```

**Expected Output:**
```
Request 1: HTTP 200
Request 2: HTTP 200
Request 3: HTTP 200
...
Request 10: HTTP 200
Request 11: HTTP 429  ← RATE LIMITED!
```

### Step 3: Verify Rate Limit Error Response

```powershell
Get-Content response_11.json
```

**Expected Response:**
```json
{
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Too many refresh requests. Please try again later.",
    "retry_after": 300
  }
}
```

### Step 4: Verify Rate Limit Persists

```powershell
# Try again immediately
curl.exe -X POST http://localhost:8080/api/v1/auth/refresh `
  -b ratelimit_cookies.txt `
  -v
```

**Expected Response:**
```
HTTP/1.1 429 Too Many Requests
```

### Step 5: PowerShell Concurrent Rate Limit Test

```powershell
# Create PowerShell script for concurrent requests
$script = @'
$headers = @{
    "Content-Type" = "application/json"
}

# Login first
$loginBody = @{
    email = "ratelimit@suproxy.com"
    password = "Test123!"
} | ConvertTo-Json

$session = $null
$loginResponse = Invoke-WebRequest `
    -Uri "http://localhost:8080/api/v1/auth/login" `
    -Method POST `
    -Headers $headers `
    -Body $loginBody `
    -SessionVariable session `
    -ErrorAction Stop

Write-Host "Login successful, starting rate limit test..."
Write-Host ""

# Make 15 concurrent refresh requests
$jobs = 1..15 | ForEach-Object {
    Start-Job -ScriptBlock {
        param($requestNum, $sessionCookies)
        
        try {
            $response = Invoke-WebRequest `
                -Uri "http://localhost:8080/api/v1/auth/refresh" `
                -Method POST `
                -Headers @{"Content-Type" = "application/json"} `
                -WebSession $sessionCookies `
                -ErrorAction Stop
            
            return @{
                Request = $requestNum
                Status = $response.StatusCode
                Result = "Success"
            }
        } catch {
            return @{
                Request = $requestNum
                Status = $_.Exception.Response.StatusCode.Value__
                Result = "Rate Limited"
            }
        }
    } -ArgumentList $_, $session
}

# Wait for all jobs and collect results
$results = $jobs | Wait-Job | Receive-Job
$jobs | Remove-Job

# Display results
$results | Sort-Object Request | ForEach-Object {
    Write-Host "Request $($_.Request): HTTP $($_.Status) - $($_.Result)"
}

# Summary
Write-Host ""
Write-Host "Summary:"
$successful = ($results | Where-Object { $_.Status -eq 200 }).Count
$rateLimited = ($results | Where-Object { $_.Status -eq 429 }).Count
Write-Host "  Successful: $successful"
Write-Host "  Rate Limited: $rateLimited"
Write-Host ""
Write-Host "Expected: 10 successful, 5 rate limited"
'@

# Save and execute script
$script | Out-File -FilePath rate_limit_test.ps1 -Encoding UTF8
powershell.exe -ExecutionPolicy Bypass -File rate_limit_test.ps1
```

**Expected Output:**
```
Login successful, starting rate limit test...

Request 1: HTTP 200 - Success
Request 2: HTTP 200 - Success
...
Request 10: HTTP 200 - Success
Request 11: HTTP 429 - Rate Limited
Request 12: HTTP 429 - Rate Limited
Request 13: HTTP 429 - Rate Limited
Request 14: HTTP 429 - Rate Limited
Request 15: HTTP 429 - Rate Limited

Summary:
  Successful: 10
  Rate Limited: 5

Expected: 10 successful, 5 rate limited
```

### Step 6: Verify Rate Limit Resets After 5 Minutes

```powershell
Write-Host "Waiting 5 minutes for rate limit reset..."
Start-Sleep -Seconds 300

# Try again after 5 minutes
curl.exe -X POST http://localhost:8080/api/v1/auth/refresh `
  -b ratelimit_cookies.txt `
  -v
```

**Expected Response:**
```
HTTP/1.1 200 OK
```

**✅ Verify:** Rate limit resets after window expires

### Step 7: Test Different IP Addresses (Optional)

If you have access to multiple IPs or a proxy:

```powershell
# Request from IP 1 (10 times) - succeeds
# Request from IP 2 (10 times) - also succeeds
# Rate limits are per-IP, not global
```

---

## 8. Multi-Tab Sync Browser Test

This test verifies frontend multi-tab synchronization for token refresh and logout events.

### Prerequisites

- Frontend running at `http://localhost:3000`
- Modern browser (Chrome, Firefox, Edge)
- Browser DevTools

### Test 1: BroadcastChannel Initialization (Modern Browsers)

**Steps:**
1. Open Chrome or Firefox
2. Navigate to `http://localhost:3000/admin`
3. Login with valid credentials
4. Open DevTools (F12) → Console tab
5. Look for initialization message

**Expected Console Output:**
```
[MULTI-TAB-SYNC] BroadcastChannel initialized
```

**✅ Verify:** BroadcastChannel API is being used (modern browser path)

### Test 2: localStorage Fallback (Legacy Browsers)

**Steps:**
1. Open Safari < 15.4 or simulate by disabling BroadcastChannel
2. Open Console
3. To simulate fallback in Chrome:
   ```javascript
   // In console, before page load
   window.BroadcastChannel = undefined;
   // Then refresh page
   ```

**Expected Console Output:**
```
[MULTI-TAB-SYNC] BroadcastChannel not supported, using localStorage fallback
[MULTI-TAB-SYNC] localStorage fallback initialized
```

**✅ Verify:** Fallback mechanism works

### Test 3: Token Refresh Notification

**Method 1: Wait for Natural Refresh (15 minutes)**

**Steps:**
1. Open Tab 1: `http://localhost:3000/admin`
2. Open Tab 2: `http://localhost:3000/admin`
3. Keep DevTools Console open in both tabs
4. Wait for token to expire (15 minutes)
5. Watch for automatic refresh

**Expected Console Output:**

**Tab 1 (where refresh happens):**
```
[API-CLIENT] Access token expired, attempting refresh...
[API-CLIENT] Token refresh successful
[MULTI-TAB-SYNC] Notifying other tabs: token_refreshed
```

**Tab 2 (receives notification):**
```
[MULTI-TAB-SYNC] Token refreshed in another tab
```

**Method 2: Force Refresh (Testing)**

Create a test button to trigger manual refresh:

**Steps:**
1. Open `http://localhost:3000/admin`
2. Open browser console
3. Execute:
   ```javascript
   // Force token refresh
   fetch('/api/v1/auth/refresh', {
     method: 'POST',
     credentials: 'include'
   }).then(r => r.json()).then(console.log);
   ```
4. Open second tab and watch console

**Expected:** Second tab logs "Token refreshed in another tab"

### Test 4: Logout Synchronization

**Steps:**
1. Open Tab 1: `http://localhost:3000/admin`
2. Open Tab 2: `http://localhost:3000/admin/plans`
3. Open Tab 3: `http://localhost:3000/admin/servers`
4. Open DevTools Console in all tabs
5. In Tab 1, click "Logout" button

**Expected Behavior:**

**Tab 1 Console:**
```
[API-CLIENT] Logout initiated
[API-CLIENT] Redirecting to login...
[MULTI-TAB-SYNC] Notifying other tabs: logout (reason: session_expired)
```

**Tab 2 & Tab 3 Console:**
```
[MULTI-TAB-SYNC] Logout detected in another tab (reason: session_expired)
```

**Visual Behavior:**
- Tab 1: Immediate redirect to `/login?reason=session_expired`
- Tab 2: Immediate redirect to `/login?reason=session_expired`
- Tab 3: Immediate redirect to `/login?reason=session_expired`

**✅ Verify:** All tabs redirect simultaneously (within ~100ms)

### Test 5: Multiple Tabs (Stress Test)

**Steps:**
1. Open 10 tabs with `http://localhost:3000/admin`
2. Login in all tabs
3. Open Console in Tab 1 and Tab 10
4. Logout in Tab 5

**Expected Behavior:**
- All 10 tabs redirect to login page immediately
- Console in Tab 1 and Tab 10 both show logout notification

**Timing Test:**
```
Tab 5 logout: 10:00:00.000
Tab 1 redirect: 10:00:00.050 (50ms delay)
Tab 10 redirect: 10:00:00.080 (80ms delay)
```

**✅ Verify:** All tabs sync within 100ms

### Test 6: Cross-Window Scenario

**Steps:**
1. Open Window 1: `http://localhost:3000/admin`
2. Open Window 2 (new browser window): `http://localhost:3000/admin`
3. Position windows side-by-side
4. Logout in Window 1
5. Watch Window 2

**Expected:** Window 2 redirects immediately (BroadcastChannel works across windows)

### Test 7: Private/Incognito Mode

**Steps:**
1. Open incognito window
2. Navigate to `http://localhost:3000/admin`
3. Open Console
4. Login

**Expected Console:**
```
[MULTI-TAB-SYNC] BroadcastChannel initialized
```

**Open second incognito tab:**

**Expected:** Tabs in same incognito session sync, but isolated from regular tabs

### Test 8: Event Deduplication

**Steps:**
1. Open 2 tabs
2. Open Console in Tab 2
3. In browser console (Tab 2), manually trigger duplicate events:
   ```javascript
   // Simulate duplicate notification
   localStorage.setItem('multi-tab-sync-event', JSON.stringify({
     type: 'logout',
     timestamp: Date.now(),
     reason: 'test'
   }));
   
   // Try same event again (< 5 seconds)
   setTimeout(() => {
     localStorage.setItem('multi-tab-sync-event', JSON.stringify({
       type: 'logout',
       timestamp: Date.now() - 3000, // Old timestamp
       reason: 'test'
     }));
   }, 1000);
   ```

**Expected Console:**
```
[MULTI-TAB-SYNC] Logout detected in another tab (reason: test)
[MULTI-TAB-SYNC] Ignoring stale event (5023ms old)
```

**✅ Verify:** Duplicate/stale events are ignored

### Test 9: SSR Safety (Next.js)

**Steps:**
1. Check server console (where `npm run dev` is running)
2. Look for any errors related to `window` or `BroadcastChannel`

**Expected:** No SSR errors (module guards against `typeof window === 'undefined'`)

### Test 10: Cleanup on Unmount

**Steps:**
1. Open tab at `http://localhost:3000/admin`
2. Open Console
3. Navigate away (e.g., go to `/login`)
4. Navigate back to `/admin`

**Expected Console:**
```
[MULTI-TAB-SYNC] BroadcastChannel initialized  ← First visit
[MULTI-TAB-SYNC] Cleanup completed              ← On navigate away
[MULTI-TAB-SYNC] BroadcastChannel initialized  ← On navigate back
```

**✅ Verify:** No memory leaks or duplicate listeners

---

## 9. Backend Log Verification

### Step 1: Configure Logging

Ensure backend is running with console logging:

```powershell
cd c:\Users\Tuncay\Desktop\suproxy-backend

# Run with logging to file
.\bin\api.exe 2>&1 | Tee-Object -FilePath backend.log
```

### Step 2: Monitor Real-Time Security Events

```powershell
# In separate terminal, monitor security events
Get-Content backend.log -Wait | Select-String "SECURITY ALERT", "Token family revoked", "security.incident"
```

### Step 3: Trigger Security Event

```powershell
# In another terminal, trigger token reuse attack
curl.exe -X POST http://localhost:8080/api/v1/auth/login `
  -H "Content-Type: application/json" `
  -d '{\"email\": \"test@suproxy.com\", \"password\": \"Test123!\"}' `
  -c attack.txt

# Refresh once
curl.exe -X POST http://localhost:8080/api/v1/auth/refresh -b attack.txt -c attack2.txt

# Reuse old token
curl.exe -X POST http://localhost:8080/api/v1/auth/refresh -b attack.txt
```

### Step 4: Verify Log Output

**Expected Log Output:**
```
[2025-01-15 10:00:05] ERROR SECURITY ALERT: Revoked token reuse detected - possible token theft
  token_id=a1b2c3d4-e5f6-7890-abcd-ef1234567890
  family_id=a1b2c3d4-e5f6-7890-abcd-ef1234567890
  user_id=b2c3d4e5-f6a7-8901-bcde-f12345678901
  revoked_at=2025-01-15T10:00:00Z
  time_since_revoked_seconds=5

[2025-01-15 10:00:05] WARN Token family revoked due to reuse detection
  family_id=a1b2c3d4-e5f6-7890-abcd-ef1234567890
  user_id=b2c3d4e5-f6a7-8901-bcde-f12345678901
```

### Step 5: Parse and Analyze Logs

Create PowerShell script to analyze security incidents:

```powershell
# Save as analyze_security_logs.ps1
$logFile = "backend.log"
$securityEvents = Get-Content $logFile | Select-String "SECURITY ALERT" -Context 0,5

Write-Host "Security Incident Summary"
Write-Host "========================="
Write-Host "Total Incidents: $($securityEvents.Count)"
Write-Host ""

foreach ($event in $securityEvents) {
    Write-Host "Incident at: $($event.Line -match '\[(.*?)\]')" 
    $event.Context.PostContext | ForEach-Object {
        if ($_ -match "family_id=(.+)") {
            Write-Host "  Family ID: $($matches[1])"
        }
        if ($_ -match "time_since_revoked_seconds=(\d+)") {
            Write-Host "  Time Since Revoked: $($matches[1]) seconds"
        }
    }
    Write-Host ""
}
```

**Execute:**
```powershell
powershell.exe -ExecutionPolicy Bypass -File analyze_security_logs.ps1
```

### Step 6: Verify Rate Limit Logs

```powershell
Get-Content backend.log | Select-String "rate limit", "Too many requests"
```

**Expected Output:**
```
[2025-01-15 10:05:00] WARN Rate limit exceeded for /api/v1/auth/refresh
  ip=127.0.0.1
  requests=11
  window=300s
```

### Step 7: Export Security Events to CSV

```powershell
# Create CSV report
$events = Get-Content backend.log | Select-String "SECURITY ALERT" | ForEach-Object {
    $line = $_.Line
    [PSCustomObject]@{
        Timestamp = ($line -match '\[(.*?)\]') ? $matches[1] : ''
        Event = 'Token Reuse'
        Severity = 'High'
        Details = $line
    }
}

$events | Export-Csv -Path security_incidents.csv -NoTypeInformation
Write-Host "Security report exported to security_incidents.csv"
```

---

## 10. Compilation & Type Check

### Backend Compilation

```powershell
cd c:\Users\Tuncay\Desktop\suproxy-backend

# Clean previous builds
Remove-Item -Recurse -Force bin -ErrorAction SilentlyContinue

# Build
go build -o bin\api.exe ./cmd/api

# Check exit code
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Backend compilation successful" -ForegroundColor Green
} else {
    Write-Host "❌ Backend compilation failed" -ForegroundColor Red
    exit 1
}

# Verify binary created
Get-Item bin\api.exe
```

**Expected Output:**
```
✅ Backend compilation successful

    Directory: c:\Users\Tuncay\Desktop\suproxy-backend\bin

Mode                 LastWriteTime         Length Name
----                 -------------         ------ ----
-a---          1/15/2025  10:00 AM       25678912 api.exe
```

### Backend Tests (Optional)

```powershell
# Run all tests
go test ./... -v

# Run specific package tests
go test ./internal/application/usecase/auth -v
go test ./internal/infrastructure/repository -v
```

### Frontend Type Check

```powershell
cd c:\Users\Tuncay\Desktop\suproxy-admin

# Type check (no emit)
npx tsc --noEmit

# Check exit code
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Frontend type check successful" -ForegroundColor Green
} else {
    Write-Host "❌ Frontend type check failed" -ForegroundColor Red
    exit 1
}
```

**Expected Output:**
```
✅ Frontend type check successful
```

### Frontend Unit Tests

```powershell
# Run multi-tab sync tests
npm test -- lib/auth/multi-tab-sync.test.ts --run

# Run all tests
npm test --run
```

**Expected Output:**
```
✓ lib/auth/multi-tab-sync.test.ts (5 tests) 18ms
  ✓ MultiTabSync (5)
    ✓ should export multiTabSync singleton
    ✓ should have notifyTokenRefreshed method
    ✓ should have notifyLogout method
    ✓ should have destroy method
    ✓ should not throw when calling methods in test environment

Test Files  1 passed (1)
     Tests  5 passed (5)
```

### Frontend Build

```powershell
# Production build
npm run build

# Check exit code
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Frontend build successful" -ForegroundColor Green
} else {
    Write-Host "❌ Frontend build failed" -ForegroundColor Red
    exit 1
}
```

**Expected Output:**
```
✓ Compiled successfully
✓ Linting and checking validity of types
✓ Collecting page data
✓ Generating static pages
✓ Finalizing page optimization

Route (app)                              Size     First Load JS
┌ ○ /                                    ...
└ ○ /admin                               ...

✅ Frontend build successful
```

### Lint Check

```powershell
# Backend
cd c:\Users\Tuncay\Desktop\suproxy-backend
golangci-lint run

# Frontend
cd c:\Users\Tuncay\Desktop\suproxy-admin
npm run lint
```

---

## 11. Security Fix Validation

### Summary Table

| Test Scenario | OLD Behavior (Bug) | NEW Behavior (Fixed) | Status |
|---------------|-------------------|----------------------|--------|
| Token reused < 5 min | ❌ Family revoked | ✅ Family revoked | ✅ Same |
| Token reused > 5 min | ✅ Allowed (VULNERABILITY) | ❌ Family revoked | ✅ FIXED |
| Token reused > 1 hour | ✅ Allowed (VULNERABILITY) | ❌ Family revoked | ✅ FIXED |
| Token reused > 1 day | ✅ Allowed (VULNERABILITY) | ❌ Family revoked | ✅ FIXED |
| Security logging | Only < 5 min | Always | ✅ FIXED |
| Audit trail | Incomplete | Complete | ✅ FIXED |

### Critical Validation Steps

#### Step 1: Verify No Time Check in Code

```powershell
cd c:\Users\Tuncay\Desktop\suproxy-backend

# Search for 5-minute check (should NOT exist)
Select-String -Path "internal\application\usecase\auth\refresh_token_command.go" -Pattern "5.*time.Minute"
```

**Expected Output:**
```
(No matches found)
```

**✅ If no results, the 5-minute check has been removed!**

#### Step 2: Verify ALWAYS Revoke Logic

```powershell
# Verify code contains ALWAYS revoke
Select-String -Path "internal\application\usecase\auth\refresh_token_command.go" -Pattern "ALWAYS revoke" -Context 0,2
```

**Expected Output:**
```
Line: // ✅ FIX: ALWAYS revoke entire family on ANY revoked token reuse
Line: if err := c.refreshTokenRepo.RevokeByFamilyID(ctx, storedToken.FamilyID); err != nil {
```

#### Step 3: Manual Code Review

Open `internal\application\usecase\auth\refresh_token_command.go` and verify:

**Required Code Pattern:**
```go
if storedToken.IsRevoked {
    // Calculate time (for logging only)
    var timeSinceRevoked time.Duration
    if storedToken.RevokedAt != nil {
        timeSinceRevoked = time.Since(*storedToken.RevokedAt)
    }
    
    // ✅ ALWAYS log (no time check)
    c.logger.Error("SECURITY ALERT: Revoked token reuse detected - possible token theft", ...)
    
    // ✅ ALWAYS revoke family (no time check)
    if err := c.refreshTokenRepo.RevokeByFamilyID(ctx, storedToken.FamilyID); err != nil {
        // Handle error
    }
    
    // ✅ ALWAYS create audit log (no time check)
    securityLog := audit.NewLog(...)
    c.auditRepo.Create(ctx, securityLog)
    
    return nil, jwt.ErrInvalidToken
}
```

**❌ MUST NOT contain:**
```go
if timeSinceRevoked < 5*time.Minute {  // ❌ This check MUST NOT exist
    // Revoke family
}
```

#### Step 4: End-to-End Validation

```powershell
# Create token and revoke it
curl.exe -X POST http://localhost:8080/api/v1/auth/login -d '...' -c test.txt
curl.exe -X POST http://localhost:8080/api/v1/auth/refresh -b test.txt -c test2.txt

# Get token IDs from database
psql -U suproxy -d suproxy -h localhost -c "SELECT id, family_id, is_revoked FROM refresh_tokens WHERE user_id = 'USER_ID' ORDER BY created_at DESC LIMIT 2;"

# Manually set revoked_at to 10 minutes ago
psql -U suproxy -d suproxy -h localhost -c "UPDATE refresh_tokens SET revoked_at = NOW() - INTERVAL '10 minutes' WHERE id = 'OLD_TOKEN_ID';"

# Reset new token to active (simulate attacker scenario)
psql -U suproxy -d suproxy -h localhost -c "UPDATE refresh_tokens SET is_revoked = false, revoked_at = NULL WHERE id = 'NEW_TOKEN_ID';"

# CRITICAL TEST: Reuse old token (revoked 10 minutes ago)
curl.exe -X POST http://localhost:8080/api/v1/auth/refresh -b test.txt -v
```

**Expected:**
```
HTTP/1.1 401 Unauthorized
```

**Verify in logs:**
```powershell
Get-Content backend.log | Select-String "SECURITY ALERT" | Select-Object -Last 1
```

**Expected:**
```
SECURITY ALERT: Revoked token reuse detected - possible token theft
  time_since_revoked_seconds=600  ← 10 minutes (600 seconds)
```

**Verify family revoked:**
```powershell
psql -U suproxy -d suproxy -h localhost -c "SELECT id, is_revoked FROM refresh_tokens WHERE family_id = 'FAMILY_ID';"
```

**Expected:**
```
                  id                  | is_revoked 
--------------------------------------+------------
 old-token-uuid                       | t          ← Was revoked
 new-token-uuid                       | t          ← Now revoked (family revocation!)
```

**✅ CRITICAL: If new token is revoked, the fix is working correctly!**

---

## 12. Rollback Procedures

If issues occur during verification, follow these rollback steps.

### Rollback Database Migration

```powershell
cd c:\Users\Tuncay\Desktop\suproxy-backend

# Apply down migration
psql -U suproxy -d suproxy -h localhost -f migrations\000010_add_token_family.down.sql
```

**Expected Output:**
```
DROP INDEX
DROP INDEX
ALTER TABLE
```

**Verify rollback:**
```powershell
psql -U suproxy -d suproxy -h localhost -c "SELECT column_name FROM information_schema.columns WHERE table_name = 'refresh_tokens' AND column_name = 'family_id';"
```

**Expected:**
```
(0 rows)
```

### Restore from Backup

```powershell
# Stop backend service
Stop-Process -Name "api" -Force -ErrorAction SilentlyContinue

# Drop current database
psql -U postgres -h localhost -c "DROP DATABASE suproxy;"

# Recreate database
psql -U postgres -h localhost -c "CREATE DATABASE suproxy;"

# Restore from backup
psql -U suproxy -d suproxy -h localhost -f backup_before_phase2_20250115_100000.sql
```

### Rollback Code Changes

```powershell
cd c:\Users\Tuncay\Desktop\suproxy-backend

# View recent commits
git log --oneline -5

# Rollback to previous commit (replace COMMIT_HASH)
git reset --hard COMMIT_HASH

# Rebuild
go build -o bin\api.exe ./cmd/api
```

### Rollback Frontend Changes

```powershell
cd c:\Users\Tuncay\Desktop\suproxy-admin

# Rollback multi-tab sync
git checkout HEAD -- lib/auth/multi-tab-sync.ts
git checkout HEAD -- lib/api/client.ts

# Remove test file
Remove-Item lib\auth\multi-tab-sync.test.ts

# Rebuild
npm run build
```

---

## Verification Checklist

Use this checklist to track your verification progress:

### Database Migration
- [ ] Backup created successfully
- [ ] Migration applied without errors
- [ ] `family_id` column exists
- [ ] `family_id` is NOT NULL
- [ ] Indexes created successfully
- [ ] Existing tokens backfilled
- [ ] All existing tokens have `id = family_id`

### Normal Token Refresh
- [ ] Login successful
- [ ] Initial token is family root
- [ ] Refresh creates new token
- [ ] Old token revoked
- [ ] New token inherits `family_id`
- [ ] Multiple rotations maintain family

### Token Reuse Detection
- [ ] Reused token returns 401
- [ ] Security alert logged
- [ ] Family revoked
- [ ] Audit log created
- [ ] **CRITICAL:** Reuse after 5+ minutes also detected

### Family Revocation
- [ ] Created token family (5+ tokens)
- [ ] Reused old revoked token
- [ ] ALL family tokens revoked
- [ ] Legitimate user session terminated
- [ ] User must re-authenticate
- [ ] New family created after re-login

### Rate Limiting
- [ ] 10 requests succeed
- [ ] 11th request returns 429
- [ ] Error message correct
- [ ] Rate limit persists
- [ ] Resets after 5 minutes
- [ ] Concurrent requests handled

### Multi-Tab Sync
- [ ] BroadcastChannel initializes
- [ ] localStorage fallback works
- [ ] Token refresh syncs across tabs
- [ ] Logout syncs across tabs
- [ ] Multiple tabs (5+) sync
- [ ] No SSR errors
- [ ] Event deduplication works

### Security Fix Validation
- [ ] No 5-minute time check in code
- [ ] ALWAYS revoke logic present
- [ ] ALWAYS log logic present
- [ ] End-to-end test passes
- [ ] Logs show time_since_revoked for ALL reuse
- [ ] Family revoked for reuse at ANY time

### Compilation & Tests
- [ ] Backend compiles successfully
- [ ] Frontend type check passes
- [ ] Frontend unit tests pass
- [ ] Frontend builds successfully
- [ ] No lint errors

### Documentation
- [ ] Implementation report reviewed
- [ ] Verification guide completed
- [ ] Security fix understood
- [ ] Rollback procedures documented

---

## Success Criteria

Phase 2 Security Hardening is successfully verified if:

1. ✅ All database migrations applied correctly
2. ✅ Token family tracking works across rotations
3. ✅ **Token reuse detected at ANY time (no 5-minute limitation)**
4. ✅ Family-wide revocation works correctly
5. ✅ Rate limiting blocks after 10 requests
6. ✅ Multi-tab sync works in modern and legacy browsers
7. ✅ All compilation and type checks pass
8. ✅ All tests pass
9. ✅ No regressions in existing functionality
10. ✅ Security fix validated end-to-end

---

## Troubleshooting

### Issue: Migration fails with "column already exists"

**Solution:**
```powershell
# Check if migration already applied
psql -U suproxy -d suproxy -h localhost -c "\d refresh_tokens"

# If family_id exists, skip migration or manually verify
```

### Issue: Token reuse not detected

**Solution:**
```powershell
# Check backend logs
Get-Content backend.log | Select-String "SECURITY ALERT"

# Verify code has ALWAYS revoke logic (no time check)
Select-String -Path "internal\application\usecase\auth\refresh_token_command.go" -Pattern "5.*time.Minute"

# Should return NO matches
```

### Issue: Multi-tab sync not working

**Solution:**
```powershell
# Check browser console for errors
# Verify BroadcastChannel or localStorage initialized
# Check for Content Security Policy blocking

# Test in different browser
# Chrome: BroadcastChannel
# Safari < 15.4: localStorage fallback
```

### Issue: Rate limit not working

**Solution:**
```powershell
# Check middleware applied in router
Select-String -Path "internal\interfaces\http\router\router.go" -Pattern "RefreshTokenRateLimiter"

# Should find: auth.POST("/refresh", middleware.RefreshTokenRateLimiter(), ...)
```

### Issue: Compilation errors

**Solution:**
```powershell
# Backend
go mod tidy
go build ./cmd/api

# Frontend
npm install
npx tsc --noEmit
npm run build
```

---

## Contact & Support

For issues or questions:
1. Check backend logs: `backend.log`
2. Check database state with provided SQL queries
3. Review implementation reports:
   - `PHASE_2_SECURITY_HARDENING_IMPLEMENTATION.md`
   - `MULTI_TAB_SYNC_IMPLEMENTATION.md`

---

**End of Phase 2 Verification Guide**

*Last Updated: January 2025*
