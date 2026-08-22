# SuProxy Load Testing Guide

## Overview

This guide provides load testing procedures for the SuProxy VPN platform to validate production readiness under high concurrent load.

---

## Prerequisites

### Install Load Testing Tools

**Option A: hey (Simple HTTP load generator)**
```bash
# macOS
brew install hey

# Linux
go install github.com/rakyll/hey@latest

# Windows
go install github.com/rakyll/hey@latest
```

**Option B: k6 (Advanced scripting)**
```bash
# macOS
brew install k6

# Linux
sudo gpg -k
sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update
sudo apt-get install k6

# Windows
choco install k6
```

---

## Load Test Scenarios

### 1. Login Endpoint Test

**Target**: `/api/v1/auth/login`

**Expected Behavior**:
- Rate limit: 5 requests/minute/IP
- After 5 failed attempts: HTTP 429

#### Using hey

```bash
# 100 requests, 10 concurrent
hey -n 100 -c 10 -m POST \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@suproxy.com","password":"Admin123!"}' \
  http://localhost:8080/api/v1/auth/login
```

**Expected Output**:
```
Summary:
  Total:        2.5403 secs
  Slowest:      0.5231 secs
  Fastest:      0.0123 secs
  Average:      0.2456 secs

Status code distribution:
  [200] 100 responses
```

#### Using k6

```javascript
// login-load-test.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '30s', target: 50 },  // Ramp up to 50 users
    { duration: '1m', target: 50 },   // Stay at 50 users
    { duration: '30s', target: 0 },   // Ramp down to 0
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests under 500ms
    http_req_failed: ['rate<0.01'],   // Less than 1% failures
  },
};

export default function () {
  const url = 'http://localhost:8080/api/v1/auth/login';
  const payload = JSON.stringify({
    email: 'admin@suproxy.com',
    password: 'Admin123!',
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  let res = http.post(url, payload, params);
  
  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });

  sleep(1);
}
```

**Run**:
```bash
k6 run login-load-test.js
```

---

### 2. Users List Endpoint Test

**Target**: `/api/v1/admin/users`

**Expected Behavior**:
- Paginated response
- Admin authentication required
- Should handle 100+ concurrent requests

#### Using hey (with auth token)

```bash
# First, get admin token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@suproxy.com","password":"Admin123!"}' \
  | jq -r '.data.access_token')

# Load test
hey -n 1000 -c 100 \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/admin/users?page=1&limit=20
```

#### Using k6

```javascript
// users-list-load-test.js
import http from 'k6/http';
import { check, group } from 'k6';

export let options = {
  stages: [
    { duration: '1m', target: 100 },  // Ramp up to 100 users
    { duration: '3m', target: 100 },  // Stay at 100 users
    { duration: '1m', target: 0 },    // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<1000'], // 95% under 1s
    http_req_failed: ['rate<0.05'],    // Less than 5% failures
  },
};

const BASE_URL = 'http://localhost:8080';
let authToken = '';

export function setup() {
  // Login once to get token
  const loginRes = http.post(`${BASE_URL}/api/v1/auth/login`, 
    JSON.stringify({
      email: 'admin@suproxy.com',
      password: 'Admin123!',
    }),
    { headers: { 'Content-Type': 'application/json' } }
  );
  
  authToken = loginRes.json('data.access_token');
  return { token: authToken };
}

export default function (data) {
  const params = {
    headers: {
      'Authorization': `Bearer ${data.token}`,
    },
  };

  group('Users List', function () {
    let res = http.get(`${BASE_URL}/api/v1/admin/users?page=1&limit=20`, params);
    
    check(res, {
      'status is 200': (r) => r.status === 200,
      'has data array': (r) => r.json('data') !== undefined,
      'response time < 1s': (r) => r.timings.duration < 1000,
    });
  });
}
```

**Run**:
```bash
k6 run users-list-load-test.js
```

---

### 3. Plans List Endpoint Test

**Target**: `/api/v1/admin/plans`

```bash
# Using hey
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@suproxy.com","password":"Admin123!"}' \
  | jq -r '.data.access_token')

hey -n 500 -c 50 \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/admin/plans?page=1&limit=20
```

---

### 4. Servers List Endpoint Test

**Target**: `/api/v1/admin/servers`

```bash
# Using hey
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@suproxy.com","password":"Admin123!"}' \
  | jq -r '.data.access_token')

hey -n 500 -c 50 \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/admin/servers?page=1&limit=20
```

---

### 5. Rate Limiter Test

**Target**: Verify rate limiting works

```bash
# Should see HTTP 429 after 5 attempts
for i in {1..10}; do
  curl -X POST http://localhost:8080/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"test@test.com","password":"wrong"}' \
    -w "\nStatus: %{http_code}\n" \
    -s -o /dev/null
  sleep 0.5
done
```

**Expected Output**:
```
Status: 401
Status: 401
Status: 401
Status: 401
Status: 401
Status: 429  <-- Rate limit triggered
Status: 429
...
```

---

### 6. Mixed Workload Test

```javascript
// mixed-workload.js
import http from 'k6/http';
import { check, group, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '2m', target: 200 },  // Ramp to 200 users
    { duration: '5m', target: 200 },  // Hold 200 users
    { duration: '2m', target: 0 },    // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(99)<2000'], // 99% under 2s
    http_req_failed: ['rate<0.1'],     // Less than 10% failures
  },
};

const BASE_URL = 'http://localhost:8080';

export function setup() {
  const loginRes = http.post(`${BASE_URL}/api/v1/auth/login`, 
    JSON.stringify({
      email: 'admin@suproxy.com',
      password: 'Admin123!',
    }),
    { headers: { 'Content-Type': 'application/json' } }
  );
  
  return { token: loginRes.json('data.access_token') };
}

export default function (data) {
  const params = {
    headers: { 'Authorization': `Bearer ${data.token}` },
  };

  // 40% users list
  if (Math.random() < 0.4) {
    group('Users', () => {
      http.get(`${BASE_URL}/api/v1/admin/users?page=1&limit=20`, params);
    });
  }
  // 30% plans list
  else if (Math.random() < 0.7) {
    group('Plans', () => {
      http.get(`${BASE_URL}/api/v1/admin/plans?page=1&limit=20`, params);
    });
  }
  // 30% servers list
  else {
    group('Servers', () => {
      http.get(`${BASE_URL}/api/v1/admin/servers?page=1&limit=20`, params);
    });
  }

  sleep(Math.random() * 2); // Random think time 0-2s
}
```

**Run**:
```bash
k6 run mixed-workload.js
```

---

## Performance Targets

### Response Time Targets

| Endpoint | p50 | p95 | p99 |
|----------|-----|-----|-----|
| Login | <100ms | <300ms | <500ms |
| Users List | <200ms | <500ms | <1000ms |
| Plans List | <150ms | <400ms | <800ms |
| Servers List | <200ms | <500ms | <1000ms |

### Throughput Targets

- **Login**: 100 req/s (with rate limiting)
- **Admin API**: 500 req/s (authenticated)
- **Database**: 1000+ queries/s

### Concurrent Users

- **Target**: 1000 concurrent users
- **Database connections**: 50 max open, 25 idle

---

## Monitoring During Load Tests

### 1. Application Metrics

```bash
# Watch metrics endpoint
watch -n 1 'curl -s http://localhost:8080/metrics'
```

### 2. Database Connections

```sql
-- Active connections
SELECT count(*) as active_connections 
FROM pg_stat_activity 
WHERE datname = 'suproxy';

-- Connection pool status
SELECT 
  max_conn,
  used,
  res_for_super,
  max_conn - used - res_for_super as available
FROM 
  (SELECT count(*) used FROM pg_stat_activity) t1,
  (SELECT setting::int res_for_super FROM pg_settings WHERE name = 'superuser_reserved_connections') t2,
  (SELECT setting::int max_conn FROM pg_settings WHERE name = 'max_connections') t3;
```

### 3. System Resources

```bash
# CPU and Memory
htop

# Network
iftop

# Disk I/O
iostat -x 1
```

---

## Interpreting Results

### Good Results ✅

```
Summary:
  Total:        10.2314 secs
  Slowest:      0.8521 secs
  Fastest:      0.0231 secs
  Average:      0.1024 secs
  Requests/sec: 97.74
  
Status code distribution:
  [200] 1000 responses

Response time histogram:
  0.023 [1]     |
  0.106 [756]   |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.189 [201]   |■■■■■■■■■■
  0.272 [32]    |■■
  0.355 [7]     |
  0.438 [2]     |
  0.521 [1]     |
```

### Warning Signs ⚠️

- **High p99 latency**: >2s indicates bottleneck
- **Error rate >5%**: Connection or timeout issues
- **Increasing response times**: Memory leak or resource exhaustion
- **HTTP 500 errors**: Application crashes

### Action Items

If performance issues found:
1. Check database slow query log
2. Verify connection pool settings
3. Review application logs for errors
4. Check system resources (CPU, memory, disk I/O)
5. Verify network latency
6. Consider horizontal scaling

---

## Production Load Test Checklist

- [ ] Test on staging environment first
- [ ] Notify team before production load test
- [ ] Set up monitoring dashboards
- [ ] Prepare rollback plan
- [ ] Document baseline metrics
- [ ] Start with small load, gradually increase
- [ ] Monitor database connections
- [ ] Monitor application logs
- [ ] Monitor system resources
- [ ] Test rate limiting behavior
- [ ] Test graceful degradation
- [ ] Verify auto-scaling triggers (if configured)

---

## Emergency Response

If load test causes production issues:

```bash
# 1. Stop load test immediately
Ctrl+C

# 2. Check application health
curl http://localhost:8080/ready

# 3. Check database connections
psql -c "SELECT count(*) FROM pg_stat_activity;"

# 4. Restart if needed
sudo systemctl restart suproxy-backend

# 5. Review logs
sudo journalctl -u suproxy-backend -n 100
```

---

## Next Steps

After successful load testing:
1. Document baseline performance metrics
2. Set up continuous performance monitoring
3. Configure alerting thresholds
4. Plan capacity scaling strategy
5. Schedule regular load tests (monthly/quarterly)
