package main

import (
	"Tournament/internal/config"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	// Define subcommands
	generateCmd := flag.NewFlagSet("migrations:generate", flag.ExitOnError)
	migrateCmd := flag.NewFlagSet("migrations:migrate", flag.ExitOnError)
	executeCmd := flag.NewFlagSet("migrations:execute", flag.ExitOnError)

	// Define flags for execute command
	executeUp := executeCmd.Bool("up", false, "Execute migration in up direction")
	executeDown := executeCmd.Bool("down", false, "Execute migration in down direction")

	// Check if a subcommand is provided
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Get migrations directory path
	migrationsDir := getMigrationsDir()

	// Parse subcommand
	switch os.Args[1] {
	case "migrations:generate":
		if generateCmd.NArg() != 0 {
			fmt.Printf("Error: Expected 0 arguments but exactly %d passed", generateCmd.NArg())
			fmt.Println("Usage: migrations:generate")
			os.Exit(1)
		}
		generateMigration(migrationsDir)

	case "migrations:migrate":
		if migrateCmd.NArg() != 0 {
			fmt.Printf("Error: Expected 0 arguments but exactly %d passed", generateCmd.NArg())
			fmt.Println("Usage: migrations:generate")
			os.Exit(1)
		}
		runAllMigrations(migrationsDir)

	case "migrations:execute":
		executeCmd.Parse(os.Args[2:])
		if executeCmd.NArg() == 0 {
			fmt.Println("Error: Migration name is required")
			fmt.Println("Usage: migrations:execute <migration_name> --up|--down")
			os.Exit(1)
		}
		if !*executeUp && !*executeDown {
			fmt.Println("Error: Either --up or --down flag must be specified")
			fmt.Println("Usage: migrations:execute <migration_name> --up|--down")
			os.Exit(1)
		}
		if *executeUp && *executeDown {
			fmt.Println("Error: Only one of --up or --down can be specified")
			fmt.Println("Usage: migrations:execute <migration_name> --up|--down")
			os.Exit(1)
		}

		isUp := *executeUp
		executeMigration(migrationsDir, executeCmd.Arg(0), isUp)

	default:
		printUsage()
		os.Exit(1)
	}
}

// printUsage prints the usage information for the migrate command
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  migrations:generate  					- Generate a new migration")
	fmt.Println("  migrations:migrate                   			- Execute all migrations")
	fmt.Println("  migrations:execute <migration_name> --up|--down 	- Execute a specific migration")
}

// getMigrationsDir determines the migrations directory path based on the current working directory
func getMigrationsDir() string {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current working directory: %v\n", err)
		os.Exit(1)
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

	// Create migrations directory if it doesn't exist
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		fmt.Printf("Error creating migrations directory: %v\n", err)
		os.Exit(1)
	}

	return migrationsDir
}

// generateMigration generates new migration files with the given name
func generateMigration(migrationsDir string) {
	timestamp := time.Now().Format("20060102150405") // Format: YYYYMMDDHHMMSS

	upMigrationPath := filepath.Join(migrationsDir, fmt.Sprintf("%s.up.sql", timestamp))
	downMigrationPath := filepath.Join(migrationsDir, fmt.Sprintf("%s.down.sql", timestamp))

	// Create up migration file with template content
	upFile, err := os.Create(upMigrationPath)
	if err != nil {
		fmt.Printf("Error creating up migration file: %v\n", err)
		os.Exit(1)
	}
	defer upFile.Close()

	upTemplate := fmt.Sprintf("-- Migration: %s\n-- Created at: %s\n\n-- Write your UP migration SQL here\n\n",
		timestamp,
		time.Now().Format("2006-01-02 15:04:05"))

	if _, err := upFile.WriteString(upTemplate); err != nil {
		fmt.Printf("Error writing to up migration file: %v\n", err)
		os.Exit(1)
	}

	// Create down migration file with template content
	downFile, err := os.Create(downMigrationPath)
	if err != nil {
		fmt.Printf("Error creating down migration file: %v\n", err)
		os.Exit(1)
	}
	defer downFile.Close()

	downTemplate := fmt.Sprintf("-- Migration: %s (revert)\n-- Created at: %s\n\n-- Write your DOWN migration SQL here\n\n",
		timestamp,
		time.Now().Format("2006-01-02 15:04:05"))

	if _, err := downFile.WriteString(downTemplate); err != nil {
		fmt.Printf("Error writing to down migration file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created migration files:\n")
	fmt.Printf("  %s\n", upMigrationPath)
	fmt.Printf("  %s\n", downMigrationPath)
}

// runAllMigrations executes all migrations
func runAllMigrations(migrationsDir string) {
	m, err := createMigrateInstance(migrationsDir)
	if err != nil {
		fmt.Printf("Error creating migrate instance: %v\n", err)
		os.Exit(1)
	}
	// Run migrations
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		fmt.Printf("Error running migrations: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("All migrations executed successfully")
}

// executeMigration executes a specific migration
func executeMigration(migrationsDir, migrationName string, isUp bool) {
	migrateInstance, err := createMigrateInstance(migrationsDir)
	if err != nil {
		fmt.Printf("Error creating migrate instance: %v\n", err)
		os.Exit(1)
	}

	// Find the migration version
	_, err = findMigrationVersion(migrationsDir, migrationName)
	if err != nil {
		fmt.Printf("Error finding migration: %v\n", err)
		os.Exit(1)
	}

	// Execute the migration
	var migrationErr error
	if isUp {
		migrationErr = migrateInstance.Steps(1) // Execute one step up
	} else {
		migrationErr = migrateInstance.Steps(-1) // Execute one step down
	}

	if migrationErr != nil && !errors.Is(migrationErr, migrate.ErrNoChange) {
		fmt.Printf("Error executing migration: %v\n", migrationErr)
		os.Exit(1)
	}

	directionStr := "up"
	if !isUp {
		directionStr = "down"
	}
	fmt.Printf("Migration %s executed successfully in %s direction\n", migrationName, directionStr)
}

// createMigrateInstance creates a new migrate instance
func createMigrateInstance(migrationsDir string) (*migrate.Migrate, error) {
	// Create database configuration
	dbConfig := config.NewDatabaseConfig()

	// Connect to the database
	db, err := sql.Open("postgres", dbConfig.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Create a new postgres driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Create a new migrate instance
	sourcePath := fmt.Sprintf("file://%s", migrationsDir)
	m, err := migrate.NewWithDatabaseInstance(sourcePath, "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return m, nil
}

// findMigrationVersion finds the version number for a migration by name
func findMigrationVersion(migrationsDir, migrationName string) (uint, error) {
	// Read the migrations directory
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return 0, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Look for a file that contains the migration name
	for _, file := range files {
		if strings.Contains(file.Name(), migrationName) && strings.HasSuffix(file.Name(), ".up.sql") {
			// Extract the version number from the filename
			versionStr := strings.TrimPrefix(file.Name(), "Version")
			versionStr = strings.TrimSuffix(versionStr, ".up.sql")

			// Parse the version number
			var version uint
			_, err := fmt.Sscanf(versionStr, "%d", &version)
			if err != nil {
				return 0, fmt.Errorf("failed to parse version number: %w", err)
			}

			return version, nil
		}
	}

	return 0, fmt.Errorf("migration %s not found", migrationName)
}
