package main

import (
"context"
"database/sql"
"fmt"
"log"
"strings"
"time"

_ "github.com/lib/pq"
)

func main() {
connStr := "host=localhost port=5432 user=suproxy password=suproxy dbname=suproxy sslmode=disable"
db, err := sql.Open("postgres", connStr)
if err != nil {
log.Fatal("Connection error:", err)
}
defer db.Close()

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

if err := db.PingContext(ctx); err != nil {
log.Fatal("Ping error:", err)
}

// Step 1: Check current counts
fmt.Println("=== STEP 1: Current Data ===")
var userCount, serverCount, planCount, auditCount int64
db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
db.QueryRow("SELECT COUNT(*) FROM servers").Scan(&serverCount)
db.QueryRow("SELECT COUNT(*) FROM plans").Scan(&planCount)
db.QueryRow("SELECT COUNT(*) FROM audit_logs").Scan(&auditCount)
fmt.Printf("users: %d, servers: %d, plans: %d, audit_logs: %d\n\n", userCount, serverCount, planCount, auditCount)

// Step 2: Seed data if needed
if serverCount < 20 || planCount < 20 || auditCount < 1000 {
fmt.Println("=== STEP 2: Seeding Data ===")

// Seed servers
if serverCount < 20 {
fmt.Println("Seeding servers...")
countries := []string{"US", "UK", "DE", "FR", "JP", "SG", "CA", "AU"}
for i := int(serverCount); i < 20; i++ {
country := countries[i%len(countries)]
_, err := db.Exec(`
INSERT INTO servers (id, name, hostname, country, region, status, ip_address, port, created_at, updated_at)
VALUES (gen_random_uuid(), $1, $2, $3, 'Region1', 'active', $4, 8080, NOW(), NOW())`,
fmt.Sprintf("Server-%s-%d", country, i),
fmt.Sprintf("server%d.example.com", i),
country,
fmt.Sprintf("192.168.1.%d", i+1))
if err != nil {
log.Printf("Error seeding server: %v", err)
}
}
fmt.Println("✓ Servers seeded")
}

// Seed plans
if planCount < 20 {
fmt.Println("Seeding plans...")
for i := int(planCount); i < 20; i++ {
_, err := db.Exec(`
INSERT INTO plans (id, name, description, price, duration_days, bandwidth_limit, max_devices, active, created_at, updated_at)
VALUES (gen_random_uuid(), $1, $2, $3, 30, 100, 3, true, NOW(), NOW())`,
fmt.Sprintf("Plan %d", i),
fmt.Sprintf("Test plan %d", i),
10+(i*5))
if err != nil {
log.Printf("Error seeding plan: %v", err)
}
}
fmt.Println("✓ Plans seeded")
}

// Seed audit logs
if auditCount < 1000 {
fmt.Println("Seeding audit logs...")
actions := []string{"user.login", "user.logout", "user.create", "server.update", "plan.create"}
for i := int(auditCount); i < 1000; i++ {
action := actions[i%len(actions)]
_, err := db.Exec(`
INSERT INTO audit_logs (id, action, entity_type, entity_id, actor_type, actor_id, ip_address, user_agent, created_at)
VALUES (gen_random_uuid(), $1, 'user', gen_random_uuid(), 'user', gen_random_uuid(), '127.0.0.1', 'test-agent', NOW() - ($2 || ' hours')::INTERVAL)`,
action, i)
if err != nil {
log.Printf("Error seeding audit log: %v", err)
}
}
fmt.Println("✓ Audit logs seeded")
}

// Update counts
db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
db.QueryRow("SELECT COUNT(*) FROM servers").Scan(&serverCount)
db.QueryRow("SELECT COUNT(*) FROM plans").Scan(&planCount)
db.QueryRow("SELECT COUNT(*) FROM audit_logs").Scan(&auditCount)
fmt.Printf("\nFinal: users=%d, servers=%d, plans=%d, audit_logs=%d\n\n", userCount, serverCount, planCount, auditCount)
}

// Step 3: Run ANALYZE
fmt.Println("=== STEP 3: Running ANALYZE ===")
_, err = db.Exec("ANALYZE users")
if err != nil {
log.Printf("Error analyzing users: %v", err)
} else {
fmt.Println("✓ ANALYZE users")
}
_, err = db.Exec("ANALYZE servers")
if err != nil {
log.Printf("Error analyzing servers: %v", err)
} else {
fmt.Println("✓ ANALYZE servers")
}
_, err = db.Exec("ANALYZE plans")
if err != nil {
log.Printf("Error analyzing plans: %v", err)
} else {
fmt.Println("✓ ANALYZE plans")
}
_, err = db.Exec("ANALYZE audit_logs")
if err != nil {
log.Printf("Error analyzing audit_logs: %v", err)
} else {
fmt.Println("✓ ANALYZE audit_logs")
}
fmt.Println()

// Step 4: EXPLAIN ANALYZE queries
fmt.Println("=== STEP 4: EXPLAIN ANALYZE ===\n")

queries := []struct {
name  string
query string
}{
{"users COUNT", "EXPLAIN (ANALYZE, BUFFERS) SELECT count(*) FROM users;"},
{"servers LIST", "EXPLAIN (ANALYZE, BUFFERS) SELECT * FROM servers ORDER BY country ASC, name ASC LIMIT 20;"},
{"plans LIST", "EXPLAIN (ANALYZE, BUFFERS) SELECT * FROM plans ORDER BY price ASC;"},
{"audit_logs COUNT", "EXPLAIN (ANALYZE, BUFFERS) SELECT count(*) FROM audit_logs;"},
{"audit_logs DATE filter", "EXPLAIN (ANALYZE, BUFFERS) SELECT count(*) FROM audit_logs WHERE created_at >= now() - interval '24 hours';"},
{"servers + nodes (batch)", "EXPLAIN (ANALYZE, BUFFERS) SELECT server_id, COUNT(*) FROM nodes WHERE server_id IN (SELECT id FROM servers LIMIT 20) GROUP BY server_id;"},
}

sep := strings.Repeat("=", 80)
for _, q := range queries {
fmt.Printf("%s\n%s\n%s\n\n", sep, q.name, sep)

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
rows, err := db.QueryContext(ctx, q.query)
cancel()

if err != nil {
log.Printf("ERROR: %v\n\n", err)
continue
}

for rows.Next() {
var line string
if err := rows.Scan(&line); err != nil {
log.Printf("Scan error: %v\n", err)
continue
}
fmt.Println(line)
}
rows.Close()
fmt.Println()
}

fmt.Println(sep)
fmt.Println("VERIFICATION COMPLETE")
fmt.Println(sep)
}