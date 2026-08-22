package main

import (
"database/sql"
"fmt"
"log"

_ "github.com/lib/pq"
)

func main() {
// Connect to database
dsn := "host=localhost port=5432 user=suproxy password=suproxy dbname=suproxy sslmode=disable"
db, err := sql.Open("postgres", dsn)
if err != nil {
log.Fatal("Failed to connect: ", err)
}
defer db.Close()

// Check current state
var version int
var dirty bool
err = db.QueryRow("SELECT version, dirty FROM schema_migrations").Scan(&version, &dirty)
if err != nil {
log.Fatal("Failed to query schema_migrations: ", err)
}

fmt.Printf("Current state: version=%d, dirty=%v\n", version, dirty)

if !dirty {
fmt.Println("Database is clean, no action needed")
return
}

// Fix dirty state by setting dirty=false
fmt.Println("Fixing dirty state...")
_, err = db.Exec("UPDATE schema_migrations SET dirty = false")
if err != nil {
log.Fatal("Failed to fix dirty state: ", err)
}

fmt.Println("✅ Dirty state fixed! Database is now clean.")
fmt.Println("You can now run the API: .\api.exe")
}
