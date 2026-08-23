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
        log.Fatal("Failed to connect:", err)
    }
    defer db.Close()
    
    // Check current state
    var version int
    var dirty bool
    err = db.QueryRow("SELECT version, dirty FROM schema_migrations").Scan(&version, &dirty)
    if err != nil {
        log.Fatal("Failed to query:", err)
    }
    fmt.Printf("Current state: version=%d, dirty=%v\n", version, dirty)
    
    // Fix dirty state
    _, err = db.Exec("UPDATE schema_migrations SET version = 10, dirty = false WHERE version = 11")
    if err != nil {
        log.Fatal("Failed to update:", err)
    }
    
    // Verify
    err = db.QueryRow("SELECT version, dirty FROM schema_migrations").Scan(&version, &dirty)
    if err != nil {
        log.Fatal("Failed to verify:", err)
    }
    fmt.Printf("Fixed state: version=%d, dirty=%v\n", version, dirty)
    fmt.Println("✓ Migration state fixed. Now restart backend to apply migration 11.")
}
