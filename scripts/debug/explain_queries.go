package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/lib/pq"
)

func main() {
    connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        "localhost", "5432", "suproxy", "suproxy", "suproxy", "disable",
    )

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    queries := []string{
        "EXPLAIN ANALYZE SELECT count(*) FROM audit_logs;",
        "EXPLAIN ANALYZE SELECT count(*) FROM audit_logs WHERE created_at >= now() - interval '24 hours';",
        "EXPLAIN ANALYZE SELECT * FROM servers ORDER BY country ASC, name ASC LIMIT 20;",
        "EXPLAIN ANALYZE SELECT * FROM plans ORDER BY price ASC;",
    }

    for i, query := range queries {
        fmt.Printf("\n========== QUERY %d ==========\n", i+1)
        fmt.Printf("%s\n\n", query)

        rows, err := db.Query(query)
        if err != nil {
            log.Printf("Error running query %d: %v\n", i+1, err)
            continue
        }

        fmt.Println("EXPLAIN ANALYZE OUTPUT:")
        for rows.Next() {
            var line string
            if err := rows.Scan(&line); err != nil {
                log.Printf("Error scanning row: %v\n", err)
                continue
            }
            fmt.Println(line)
        }
        rows.Close()
    }
}
