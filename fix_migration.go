package main

import (
    "database/sql"
    "fmt"
    "log"
    
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/lib/pq"
)

func main() {
    dsn := "host=localhost port=5432 user=suproxy password=suproxy dbname=suproxy sslmode=disable"
    
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        log.Fatal(err)
    }
    
    m, err := migrate.NewWithDatabaseInstance(
        "file://C:/Users/Tuncay/Desktop/suproxy-backend/migrations",
        "postgres",
        driver,
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Force version 10 (clean state before migration 11)
    if err := m.Force(10); err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Forced migration version to 10")
    
    // Now run Up to apply migration 11
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal(err)
    }
    
    fmt.Println("Migration 11 applied successfully")
}
