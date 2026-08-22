package main
import (
"database/sql"
"fmt"
"time"
_ "github.com/lib/pq"
)
func main() {
dsn := "host=localhost port=5432 user=suproxy password=suproxy dbname=suproxy sslmode=disable"

fmt.Println("Connecting to PostgreSQL...")
start := time.Now()

db, err := sql.Open("postgres", dsn)
if err != nil {
fmt.Printf("Open error: %v\n", err)
return
}
defer db.Close()

// ACTUAL CONNECTION HAPPENS HERE
err = db.Ping()
duration := time.Since(start)

fmt.Printf("\n=== RESULT ===\n")
fmt.Printf("Connection time: %v\n", duration)
fmt.Printf("Error: %v\n", err)

if duration > 500*time.Millisecond {
fmt.Println("\n🔴 PROBLEM FOUND: Connection is SLOW (>500ms)")
fmt.Println("   This explains GORM [730ms] logs!")
} else {
fmt.Println("\n✅ Connection is FAST (<500ms)")
fmt.Println("   Problem is NOT connection acquisition")
}
}
