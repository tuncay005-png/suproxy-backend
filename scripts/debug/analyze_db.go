package main

import (
"database/sql"
"fmt"
"log"
"time"

_ "github.com/lib/pq"
)

func main() {
connStr := "host=localhost port=5432 user=suproxy password=suproxy dbname=suproxy sslmode=disable"
db, err := sql.Open("postgres", connStr)
if err != nil {
log.Fatal(err)
}
defer db.Close()

fmt.Println("=== BEFORE ANALYZE: Real COUNT timings ===")
measure := func(label, query string) {
start := time.Now()
var n int64
if err := db.QueryRow(query).Scan(&n); err != nil {
log.Printf("ERROR %s: %v", label, err)
return
}
elapsed := time.Since(start)
fmt.Printf("[%6dms] %s => %d rows\n", elapsed.Milliseconds(), label, n)
}

measure("users", "SELECT count(*) FROM users")
measure("audit_logs", "SELECT count(*) FROM audit_logs")

fmt.Println()
fmt.Println("=== RUNNING VACUUM ANALYZE ===")
tables := []string{"users", "audit_logs", "servers", "plans"}
for _, t := range tables {
start := time.Now()
if _, err := db.Exec("VACUUM ANALYZE " + t); err != nil {
log.Printf("ERROR VACUUM ANALYZE %s: %v", t, err)
} else {
fmt.Printf("✓ VACUUM ANALYZE %-15s (%dms)\n", t, time.Since(start).Milliseconds())
}
}

fmt.Println()
fmt.Println("=== AFTER ANALYZE: Real COUNT timings ===")
measure("users", "SELECT count(*) FROM users")
measure("audit_logs", "SELECT count(*) FROM audit_logs")

fmt.Println()
fmt.Println("DONE")
}
