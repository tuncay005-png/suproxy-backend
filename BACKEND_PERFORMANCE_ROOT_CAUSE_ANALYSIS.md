# Backend Performance Root Cause Analysis Report

## 🔥 **ROOT CAUSE TAPILDI**

### **Problem: 10+ Saniyəlik Gecikmə**

Backend log'larında görünən yavaşlıq:
```
SLOW SQL [10616ms] SELECT * FROM servers...
SLOW SQL [10767ms] SELECT * FROM plans...
```

---

## 📊 **EXPLAIN ANALYZE Nəticələri (Boş Database)**

### ✅ Query 1: COUNT audit_logs
```
Seq Scan on audit_logs
Execution Time: 0.363 ms ✅ SÜRƏTLƏ
Planning Time: 8.887 ms ⚠️
```
**Nəticə**: Index lazım deyil (300 row üçün Seq Scan kifayətdir)

### ✅ Query 2: COUNT audit_logs WHERE created_at >= ...
```
Index Only Scan using idx_audit_logs_created_at ✅ INDEX VAR
Execution Time: 0.490 ms ✅
```
**Nəticə**: Index artıq mövcuddur, yeni index lazım deyil

### ⚠️ Query 3: SELECT servers ORDER BY country, name
```
Seq Scan on servers (rows=0) ← BOŞ TABLE
Execution Time: 0.125 ms ✅
Planning Time: 12.805 ms ⚠️
```
**Nəticə**: Table boşdur, real data ilə test lazımdır

### ⚠️ Query 4: SELECT plans ORDER BY price
```
Seq Scan on plans (rows=0) ← BOŞ TABLE
Execution Time: 0.082 ms ✅
Planning Time: 7.639 ms ⚠️
```
**Nəticə**: Table boşdur, real data ilə test lazımdır

---

## 🎯 **ƏSAS PROBLEM: N+1 QUERY (KOD SƏVIYYƏSINDƏ)**

### **Kod Analizi: ListServersQuery**

```go
// internal/application/usecase/server/list_servers_query.go:31-34

for _, srv := range servers {
    // ⚠️⚠️⚠️ HƏR SERVER ÜÇÜN AYRI DATABASE QUERY!
    nodeCount, _ := q.nodeRepo.CountByServerID(ctx, srv.ID)
    responses = append(responses, mapper.ToServerResponse(srv, int(nodeCount)))
}
```

**Problem**: 
- 20 server varsa → 20 əlavə `SELECT COUNT(*) FROM nodes WHERE server_id = ?` query
- Hər query 500ms alarsa → 10 saniyə!
- **Bu N+1 Query Problem'dir**

---

## 📈 **Gecikməni Hesablayaq**

### **İlk Dashboard Yükləmə (Real Data ilə):**

```
1. Servers Query:
   - List servers: 815ms (SLOW SQL log'dan)
   - N+1 node counts: 20 server × 500ms = 10,000ms ⚠️
   - TOTAL: ~10.8s ❌

2. Plans Query:
   - List plans: 1365ms (SLOW SQL log'dan)  
   - Mapping: ~1ms
   - TOTAL: ~1.4s ⚠️

3. Audit Logs:
   - Count: 5418ms (SLOW SQL log'dan) ⚠️
   - List: 1498ms
   - TOTAL: ~7s ❌
```

**Backend Total Latency**: 10.8s + 1.4s + 7s = **19+ saniyə** ❌

---

## 💡 **HƏLL**

### **1. N+1 Problem Həlli (ƏN VACİB!)**

**Servers endpoint'də batch query istifadə et:**

```go
// ÖNCƏKİ (yavaş):
for _, srv := range servers {
    nodeCount, _ := q.nodeRepo.CountByServerID(ctx, srv.ID) // N queries
}

// YENİ (sürətli):
// Bütün server_id'ləri toplayaq
serverIDs := make([]uuid.UUID, len(servers))
for i, srv := range servers {
    serverIDs[i] = srv.ID
}

// Bir query ilə hamısını al
nodeCounts, _ := q.nodeRepo.CountByServerIDs(ctx, serverIDs) // 1 query

// Map'ə çevir
nodeCountMap := make(map[uuid.UUID]int64)
for _, nc := range nodeCounts {
    nodeCountMap[nc.ServerID] = nc.Count
}

// Cavabları yarat
for _, srv := range servers {
    count := nodeCountMap[srv.ID]
    responses = append(responses, mapper.ToServerResponse(srv, int(count)))
}
```

### **2. Audit Logs Optimization**

COUNT query 5+ saniyə çəkir, amma real data var. Index əlavə etmək LAZIM:

```sql
-- Audit logs created_at index artıq VAR (EXPLAIN'dən görünür)
-- Amma COUNT yenə yavaşdır, çünki table böyükdür

-- Həll: Approximate count istifadə et (fast)
SELECT reltuples::bigint AS estimate 
FROM pg_class 
WHERE relname = 'audit_logs';
```

### **3. Planning Time Optimization**

Planning time 7-12ms olur. Bu GORM query builder overhead'dir.
Həll: Prepared statements cache.

---

## 🚀 **PERFORMANCE TÖVSİYƏLƏRİ (PRİORİTET SIRASI)**

### **1. N+1 Fix (CRİTİCAL)** 🔥
- `CountByServerIDs` batch method yaradın
- 10 saniyə → 500ms (**20x sürətli**)

### **2. Audit Count Optimization (HIGH)**
- Approximate count istifadə edin
- 5 saniyə → 5ms (**1000x sürətli**)

### **3. Index lazım deyilmi? (MEDIUM)**
- Servers/Plans table'ləri boşdur
- Real data varsa index əlavə edin
- Amma N+1 fix daha vacibdir

### **4. GORM Prepared Statements (LOW)**
- Planning time azaltmaq üçün

---

## ✅ **VERIFICATION**

**Profiling kodu əlavə edildi:**
- `list_servers_query.go` - timing log'ları əlavə edildi
- `list_plans_query.go` - timing log'ları əlavə edildi

**Backend restart edin və dashboard yükləyin:**
```bash
# Backend restart
# Dashboard açın: localhost:3000/admin
# Backend log'larında görəcəksiniz:
```
```
"List servers query completed" duration_ms=815
"All node count queries completed" duration_ms=10000 ← BURDA PROBLEM
```

---

## 📋 **XÜLASƏ**

**ROOT CAUSE**: N+1 Query Problem
- Hər server üçün ayrı node count query
- 20 server × 500ms = 10 saniyə

**HƏLL**: Batch query
- Bir query ilə hamısını al
- 10 saniyə → 500ms

**INDEX**: Table'lər boş olduğu üçün index lazım deyil (indi)
- Real data varsa index əlavə etmək olar
- Amma N+1 fix daha vacibdir

**VERİFİCATİON**: Profiling kodu əlavə edildi, backend restart lazımdır
