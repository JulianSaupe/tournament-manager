package main

import (
	"Tournament/internal/config"
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/*
   Embed every migration file (*.sql) that lives in the migrations folder.
   Adjust the pattern if your folder name is different.
*/
//go:embed migrations/*.sql
var migrationFiles embed.FS

const (
	migrationsDir = "migrations"
	timeFormat    = "20060102150405"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

func main() {
	// Define commands
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)

	// Check if a command was provided
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Parse the command
	switch os.Args[1] {
	case "create":
		createCmd.Parse(os.Args[2:])
		createMigration()

	case "list":
		listCmd.Parse(os.Args[2:])
		listMigrations()

	case "run":
		runCmd.Parse(os.Args[2:])
		runMigrations()

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  migrate create -name <migration_name>  Create a new migration")
	fmt.Println("  migrate list                          List all migrations")
	fmt.Println("  migrate run                           Run all pending migrations")
}

// createMigration creates new up and down migration files with the current timestamp
func createMigration() {
	// Generate timestamp
	timestamp := time.Now().Format(timeFormat)

	// Create migration directory if it doesn't exist
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		log.Fatalf("Failed to create migrations directory: %v", err)
	}

	// Create up migration file
	upFilename := filepath.Join(migrationsDir, fmt.Sprintf("%s.up.sql", timestamp))
	if err := os.WriteFile(upFilename, []byte("-- Write your UP migration SQL here\n"), 0644); err != nil {
		log.Fatalf("Failed to create up migration file: %v", err)
	}

	// Create down migration file
	downFilename := filepath.Join(migrationsDir, fmt.Sprintf("%s.down.sql", timestamp))
	if err := os.WriteFile(downFilename, []byte("-- Write your DOWN migration SQL here\n"), 0644); err != nil {
		log.Fatalf("Failed to create down migration file: %v", err)
	}

	fmt.Printf("Created migration files:\n%s\n%s\n", upFilename, downFilename)
}

// listMigrations lists all migrations and their status
func listMigrations() {
	// Create database configuration
	dbConfig := config.NewDatabaseConfig()

	// Connect to database
	db, err := dbConfig.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create migration runner
	migrationRunner, err := config.NewMigrationRunner(db)
	if err != nil {
		log.Fatalf("Failed to create migration runner: %v", err)
	}

	// Get applied migrations
	ctx := context.Background()
	applied, err := migrationRunner.GetAppliedMigrations(ctx)
	if err != nil {
		log.Fatalf("Failed to get applied migrations: %v", err)
	}

	// Get all migration files from disk
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		// Try to read from embedded FS if directory doesn't exist
		fsEntries, fsErr := fs.ReadDir(migrationFiles, migrationsDir)
		if fsErr != nil {
			log.Fatalf("Failed to read migrations: %v", err)
		}
		entries = fsEntries
	}

	// Filter and sort up migrations
	var migrations []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			migrations = append(migrations, entry.Name())
		}
	}

	// Print migrations with status
	fmt.Println("Migrations:")
	for _, migration := range migrations {
		// Extract version from filename (e.g., "20250613164422_create_matches.up.sql" -> "20250613164422")
		version := strings.Split(migration, "_")[0]

		// Check if migration is applied
		status := colorYellow + "PENDING" + colorReset
		if applied[version] {
			status = colorGreen + "APPLIED" + colorReset
		}

		fmt.Printf("  %s - %s\n", migration, status)
	}
}

// runMigrations runs all pending migrations
func runMigrations() {
	// Create database configuration
	dbConfig := config.NewDatabaseConfig()

	// Connect to database
	db, err := dbConfig.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create migration runner
	migrationRunner, err := config.NewMigrationRunner(db)
	if err != nil {
		log.Fatalf("Failed to create migration runner: %v", err)
	}

	// Run migrations
	ctx := context.Background()
	if err := migrationRunner.RunMigrations(ctx, migrationFiles); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	fmt.Println(colorGreen + "Migrations completed successfully" + colorReset)
}
