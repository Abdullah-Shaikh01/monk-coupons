package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    "github.com/Abdullah-Shaikh01/monk-coupons/config"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    cfg := config.LoadDBConfig()

    db, err := sql.Open("mysql", cfg.DSN())
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    content, err := os.ReadFile("migrations/migrations.sql")
    if err != nil {
        log.Fatalf("Failed to read migration file: %v", err)
    }

    _, err = db.Exec(string(content))
    if err != nil {
        log.Fatalf("Failed to execute migration: %v", err)
    }

    fmt.Println("Migration completed successfully.")
}
