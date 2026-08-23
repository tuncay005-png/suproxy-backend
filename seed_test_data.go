package main

import (
"context"
"database/sql"
"fmt"
"log"
"time"

"github.com/google/uuid"
_ "github.com/lib/pq"
"golang.org/x/crypto/bcrypt"
)

func main() {
connStr := "host=localhost port=5432 user=suproxy password=suproxy dbname=suproxy sslmode=disable"
db, err := sql.Open("postgres", connStr)
if err != nil {
log.Fatal(err)
}
defer db.Close()

ctx := context.Background()

// Check current counts
var userCount, serverCount, planCount, auditCount int
db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
db.QueryRow("SELECT COUNT(*) FROM servers").Scan(&serverCount)
db.QueryRow("SELECT COUNT(*) FROM plans").Scan(&planCount)
db.QueryRow("SELECT COUNT(*) FROM audit_logs").Scan(&auditCount)

fmt.Printf("Current counts: users=%d, servers=%d, plans=%d, audit_logs=%d\n", userCount, serverCount, planCount, auditCount)

// Seed users (50 total)
if userCount < 50 {
fmt.Println("Seeding users...")
hashedPw, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
for i := userCount; i < 50; i++ {
email := fmt.Sprintf("user%d@example.com", i)
name := fmt.Sprintf("User %d", i)
_, err := db.ExecContext(ctx, `
INSERT INTO users (id, email, email_lower, name, password_hash, role, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (email) DO NOTHING`,
uuid.New(), email, email, name, string(hashedPw), "user", "active", time.Now(), time.Now())
if err != nil {
log.Printf("Error seeding user: %v", err)
}
}
fmt.Println("Users seeded.")
}

// Seed plans (20 total)
if planCount < 20 {
fmt.Println("Seeding plans...")
for i := planCount; i < 20; i++ {
name := fmt.Sprintf("Plan %d", i)
price := float64(10 + i*5)
_, err := db.ExecContext(ctx, `
INSERT INTO plans (id, name, price, duration_days, bandwidth_limit, max_devices, active, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
uuid.New(), name, price, 30, 100, 3, true, time.Now(), time.Now())
if err != nil {
log.Printf("Error seeding plan: %v", err)
}
}
fmt.Println("Plans seeded.")
}

// Seed servers (20 total)
if serverCount < 20 {
fmt.Println("Seeding servers...")
countries := []string{"US", "UK", "DE", "FR", "JP", "SG", "CA", "AU"}
for i := serverCount; i < 20; i++ {
country := countries[i%len(countries)]
name := fmt.Sprintf("Server-%s-%d", country, i)
_, err := db.ExecContext(ctx, `
INSERT INTO servers (id, name, hostname, country, region, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
uuid.New(), name, fmt.Sprintf("server%d.example.com", i), country, "Region1", "active", time.Now(), time.Now())
if err != nil {
log.Printf("Error seeding server: %v", err)
}
}
fmt.Println("Servers seeded.")
}

// Seed audit logs (1000 total)
if auditCount < 1000 {
fmt.Println("Seeding audit logs...")
actions := []string{"user.login", "user.logout", "user.create", "server.update", "plan.create"}
for i := auditCount; i < 1000; i++ {
action := actions[i%len(actions)]
createdAt := time.Now().Add(-time.Hour * time.Duration(i))
_, err := db.ExecContext(ctx, `
INSERT INTO audit_logs (id, action, entity_type, entity_id, actor_type, actor_id, ip_address, user_agent, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
uuid.New(), action, "user", uuid.New(), "user", uuid.New(), "127.0.0.1", "test-agent", createdAt)
if err != nil {
log.Printf("Error seeding audit log: %v", err)
}
}
fmt.Println("Audit logs seeded.")
}

// Final counts
db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
db.QueryRow("SELECT COUNT(*) FROM servers").Scan(&serverCount)
db.QueryRow("SELECT COUNT(*) FROM plans").Scan(&planCount)
db.QueryRow("SELECT COUNT(*) FROM audit_logs").Scan(&auditCount)

fmt.Printf("\nFinal counts: users=%d, servers=%d, plans=%d, audit_logs=%d\n", userCount, serverCount, planCount, auditCount)
fmt.Println("Seed complete!")
}
