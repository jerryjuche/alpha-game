package main

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Absolute path to migrations folder
	absPath, err := filepath.Abs("migrations")
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	// Convert Windows backslashes to forward slashes
	absPath = strings.ReplaceAll(absPath, "\\", "/")

	// Use standard file URL for golang-migrate (Windows safe)
	migrationURL := "file:///" + absPath

	// Database connection string
	dbURL := "postgres://postgres:25563214@localhost:5432/alphaplitz?sslmode=disable"

	m, err := migrate.New(migrationURL, dbURL)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("✅ Migrations applied successfully!")
}