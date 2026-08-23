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

    queries := []struct {
        name  string
        query string
    }{
        {"Query 1: COUNT audit_logs", "EXPLAIN (ANALYZE, BUFFERS) SELECT count(*) FROM audit_logs;"},
        {"Query 2: COUNT audit_logs with date filter", "EXPLAIN (ANALYZE, BUFFERS) SELECT count(*) FROM audit_logs WHERE created_at >= now() - interval '24 hours';"},
        {"Query 3: SELECT servers with ORDER BY", "EXPLAIN (ANALYZE, BUFFERS) SELECT * FROM servers ORDER BY country ASC, name ASC LIMIT 20;"},
        {"Query 4: SELECT plans with ORDER BY", "EXPLAIN (ANALYZE, BUFFERS) SELECT * FROM plans ORDER BY price ASC;"},
    }

    sep := strings.Repeat("=", 80)

    for _, q := range queries {
        fmt.Printf("\n%s\n", sep)
        fmt.Printf("%s\n", q.name)
        fmt.Printf("%s\n\n", sep)

        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        rows, err := db.QueryContext(ctx, q.query)
        cancel()

        if err != nil {
            log.Printf("ERROR: %v\n", err)
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
    }

    fmt.Printf("\n%s\n", sep)
    fmt.Println("ANALYSIS COMPLETE")
    fmt.Printf("%s\n", sep)
}