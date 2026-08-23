package main

import (
"database/sql"
"fmt"
"log"
_ "github.com/lib/pq"
)

func main() {
dsn := "host=localhost port=5432 user=suproxy password=suproxy dbname=suproxy sslmode=disable"

db, err := sql.Open("postgres", dsn)
if err != nil {
log.Fatal("Connection failed:", err)
}
defer db.Close()

// Test connection
if err := db.Ping(); err != nil {
log.Fatal("Ping failed:", err)
}

fmt.Println("✓ Connected to database")

// Check current state
var version int
var dirty bool
err = db.QueryRow("SELECT version, dirty FROM schema_migrations").Scan(&version, &dirty)
if err != nil {
log.Fatal("Query failed:", err)
}

fmt.Printf("Current state: version=%d dirty=%v\n", version, dirty)

if !dirty {
fmt.Println("✓ Database is already clean")
return
}

// Fix: Set version to 10 and dirty to false
fmt.Println("Fixing dirty state...")
_, err = db.Exec("UPDATE schema_migrations SET version = 10, dirty = false")
if err != nil {
log.Fatal("Update failed:", err)
}

// Verify
err = db.QueryRow("SELECT version, dirty FROM schema_migrations").Scan(&version, &dirty)
if err != nil {
log.Fatal("Verify failed:", err)
}

fmt.Printf("✓ Fixed state: version=%d dirty=%v\n", version, dirty)
fmt.Println("")
fmt.Println("Now run: .\\api.exe")
fmt.Println("Migration 11 will be applied automatically on startup")
}