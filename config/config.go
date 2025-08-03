package config

import (
    "fmt"
    "log"
    "os"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/joho/godotenv"
)

type DBConfig struct {
    User     string
    Password string
    Host     string
    Port     string
    Name     string
}

func LoadDBConfig() DBConfig {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    return DBConfig{
        User:     os.Getenv("DB_USER"),
        Password: os.Getenv("DB_PASS"),
        Host:     os.Getenv("DB_HOST"),
        Port:     os.Getenv("DB_PORT"),
        Name:     os.Getenv("DB_NAME"),
    }
}

func (c DBConfig) DSN() string {
    return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
        c.User, c.Password, c.Host, c.Port, c.Name)
}

func InitDB() *sql.DB {
	cfg := LoadDBConfig()
	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Optional: ping to test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	return db
}
