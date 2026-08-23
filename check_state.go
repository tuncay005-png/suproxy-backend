package main
import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

func main() {
    db, _ := sql.Open("postgres", "host=localhost port=5432 user=suproxy password=suproxy dbname=suproxy sslmode=disable")
    defer db.Close()
    
    // 1. Check schema_migrations
    var version int
    var dirty bool
    db.QueryRow("SELECT version, dirty FROM schema_migrations").Scan(&version, &dirty)
    fmt.Printf("schema_migrations: version=%d dirty=%v\n", version, dirty)
    
    // 2. Check if indexes exist
    var count int
    db.QueryRow("SELECT COUNT(*) FROM pg_indexes WHERE indexname IN ('idx_inbounds_instance_enabled_port', 'idx_clients_inbound_enabled_created', 'idx_audit_logs_created_at')").Scan(&count)
    fmt.Printf("Indexes created: %d/3\n", count)
    
    // 3. List actual indexes
    rows, _ := db.Query("SELECT indexname FROM pg_indexes WHERE indexname LIKE 'idx_%' AND tablename IN ('inbounds', 'clients', 'audit_logs') ORDER BY indexname")
    fmt.Println("\nExisting indexes:")
    for rows.Next() {
        var name string
        rows.Scan(&name)
        fmt.Printf("  - %s\n", name)
    }
}
