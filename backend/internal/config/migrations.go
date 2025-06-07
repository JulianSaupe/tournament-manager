package config

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// getMigrationsDir determines the migrations directory path based on the current working directory
func getMigrationsDir() string {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting current working directory: %v\n", err)
		return "migrations" // Fallback to default
	}

	// Determine the migrations directory path
	var migrationsDir string
	if filepath.Base(cwd) == "backend" {
		migrationsDir = "migrations"
	} else if filepath.Base(cwd) == "cmd" || filepath.Base(cwd) == "migrate" {
		migrationsDir = "../../migrations"
	} else {
		migrationsDir = "backend/migrations"
	}

	return migrationsDir
}

// MigrateDatabase runs database migrations
func MigrateDatabase(dbConfig *DatabaseConfig) error {
	// Connect to the database
	db, err := sql.Open("postgres", dbConfig.ConnectionString())
	if err != nil {
		return fmt.Errorf("failed to connect to database for migrations: %w", err)
	}
	defer db.Close()

	// Create a new postgres driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Get migrations directory path
	migrationsDir := getMigrationsDir()

	// Create a new migrate instance
	sourcePath := fmt.Sprintf("file://%s", migrationsDir)
	m, err := migrate.NewWithDatabaseInstance(sourcePath, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Run migrations
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}
