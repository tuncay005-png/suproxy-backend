package main

import (
"database/sql"
"fmt"
"log"

_ "github.com/lib/pq"
)

func main() {
// PostgreSQL bağlantısı
dsn := "host=localhost port=5432 user=suproxy password=suproxy dbname=suproxy sslmode=disable"
db, err := sql.Open("postgres", dsn)
if err != nil {
log.Fatal("Bağlantı xətası: ", err)
}
defer db.Close()

// İndiki vəziyyəti yoxla
var version int
var dirty bool
err = db.QueryRow("SELECT version, dirty FROM schema_migrations").Scan(&version, &dirty)
if err != nil {
log.Fatal("schema_migrations oxuna bilmədi: ", err)
}

fmt.Printf("İndiki vəziyyət: version=%d, dirty=%v\n", version, dirty)

if !dirty {
fmt.Println("✅ Database artıq təmizdir, heç nə etməyə ehtiyac yoxdur")
return
}

// Dirty flag-ı təmizlə
fmt.Printf("⚙️  Migration %d dirty flag-ını təmizləyirəm...\n", version)
_, err = db.Exec("UPDATE schema_migrations SET dirty = false")
if err != nil {
log.Fatal("Dirty flag təmizlənə bilmədi: ", err)
}

fmt.Println("✅ Dirty flag təmizləndi!")
fmt.Println("✅ Backend indi işləməlidir: .\api.exe")
}
