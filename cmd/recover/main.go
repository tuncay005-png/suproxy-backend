package main

import (
"fmt"
"os"
"github.com/suproxy/backend/internal/infrastructure/config"
"github.com/suproxy/backend/internal/infrastructure/database"
"github.com/suproxy/backend/internal/infrastructure/logger"
)

func main() {
cfg, err := config.Load()
if err != nil {
fmt.Printf("Config error: %v\n", err)
os.Exit(1)
}

log := logger.New(&cfg.Log)
migrator := database.NewMigrator(cfg, log)

fmt.Println("=== Migration Recovery ===")

version, dirty, err := migrator.Version()
if err != nil {
fmt.Printf("Version error: %v\n", err)
os.Exit(1)
}

fmt.Printf("Current: version=%d dirty=%v\n", version, dirty)

if !dirty {
fmt.Println("✓ Clean, no recovery needed")
return
}

fmt.Println("Forcing to version 10...")
if err := migrator.Force(10); err != nil {
fmt.Printf("Force error: %v\n", err)
os.Exit(1)
}

fmt.Println("Running migration 11...")
if err := migrator.Up(); err != nil {
fmt.Printf("Up error: %v\n", err)
os.Exit(1)
}

version, dirty, _ = migrator.Version()
fmt.Printf("✓ Final: version=%d dirty=%v\n", version, dirty)
fmt.Println("✓ Recovery complete!")
}