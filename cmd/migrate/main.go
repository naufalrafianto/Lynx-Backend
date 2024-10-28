package main

import (
    "flag"
    "fmt"
    "log"
    "os"

    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    "github.com/joho/godotenv"
)

func main() {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    // Parse command line flags
    var command string
    flag.StringVar(&command, "command", "", "migrate command (up/down)")
    flag.Parse()

    if command == "" {
        log.Fatal("Please provide a command: -command=up or -command=down")
    }

    // Construct database URL
    dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"),
    )

    // Initialize migrator
    m, err := migrate.New(
        "file://migrations",
        dbURL,
    )
    if err != nil {
        log.Fatal(err)
    }
    defer m.Close()

    // Execute migration command
    switch command {
    case "up":
        if err := m.Up(); err != nil && err != migrate.ErrNoChange {
            log.Fatal(err)
        }
        log.Println("Successfully applied up migrations")
    case "down":
        if err := m.Down(); err != nil && err != migrate.ErrNoChange {
            log.Fatal(err)
        }
        log.Println("Successfully applied down migrations")
    default:
        log.Fatal("Invalid command. Use 'up' or 'down'")
    }
}