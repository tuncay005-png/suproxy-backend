package main

import (
"context"
"database/sql"
"fmt"
"log"
"os"

_ "github.com/lib/pq"
)

func main() {
connStr := "host=localhost port=5432 user=suproxy password=suproxy dbname=suproxy sslmode=disable"
db, err := sql.Open("postgres", connStr)
if err != nil {
log.Fatal(err)
}
defer db.Close()

sqlContent, err := os.ReadFile("seed_data.sql")
if err != nil {
log.Fatal(err)
}

ctx := context.Background()
_, err = db.ExecContext(ctx, string(sqlContent))
if err != nil {
log.Fatal("Error executing seed SQL:", err)
}

fmt.Println("✅ Seed data inserted successfully")
}
