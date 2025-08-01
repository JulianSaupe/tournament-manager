package config

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// MigrationRunner handles database migrations
type MigrationRunner struct {
	db *sql.DB
}

// NewMigrationRunner creates a new migration runner
func NewMigrationRunner(db *sql.DB) (*MigrationRunner, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection cannot be nil")
	}

	return &MigrationRunner{
		db: db,
	}, nil
}

// InitMigrationTable ensures the migration tracking table exists
func (m *MigrationRunner) InitMigrationTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`
	_, err := m.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	return nil
}

// GetAppliedMigrations returns a list of already applied migrations
func (m *MigrationRunner) GetAppliedMigrations(ctx context.Context) (map[string]bool, error) {
	// Initialize migration table
	if err := m.InitMigrationTable(ctx); err != nil {
		return nil, err
	}

	query := `SELECT version FROM schema_migrations ORDER BY version`
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("failed to scan migration version: %w", err)
		}
		applied[version] = true
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating migrations: %w", err)
	}

	return applied, nil
}

// DiscoverMigrations finds all migration files in the embedded filesystem
func DiscoverMigrations(migrationFiles embed.FS) ([]string, error) {
	entries, err := fs.ReadDir(migrationFiles, "migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrations []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			migrations = append(migrations, entry.Name())
		}
	}

	// Sort migrations by version
	sort.Strings(migrations)
	return migrations, nil
}

// ApplyMigration applies a single migration
func (m *MigrationRunner) ApplyMigration(ctx context.Context, version string, content string) error {
	// Start a transaction
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Execute the migration
	_, err = tx.ExecContext(ctx, content)
	if err != nil {
		return fmt.Errorf("failed to execute migration %s: %w", version, err)
	}

	// Record the migration
	_, err = tx.ExecContext(ctx, "INSERT INTO schema_migrations (version, applied_at) VALUES ($1, $2)", version, time.Now())
	if err != nil {
		return fmt.Errorf("failed to record migration %s: %w", version, err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// RunMigrations applies all pending migrations
func (m *MigrationRunner) RunMigrations(ctx context.Context, migrationFiles embed.FS) error {
	// Initialize migration table
	if err := m.InitMigrationTable(ctx); err != nil {
		return err
	}

	// Get applied migrations
	applied, err := m.GetAppliedMigrations(ctx)
	if err != nil {
		return err
	}

	// Discover migrations
	migrations, err := DiscoverMigrations(migrationFiles)
	if err != nil {
		return err
	}

	// Apply pending migrations
	for _, migration := range migrations {
		// Extract version from filename (e.g., "20250613164422_create_matches.up.sql" -> "20250613164422")
		version := strings.Split(migration, "_")[0]

		// Skip if already applied
		if applied[version] {
			log.Printf("Migration %s already applied, skipping", version)
			continue
		}

		// Read migration content
		content, err := migrationFiles.ReadFile(filepath.Join("migrations", migration))
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", migration, err)
		}

		// Apply migration
		log.Printf("Applying migration %s", migration)
		if err := m.ApplyMigration(ctx, version, string(content)); err != nil {
			return err
		}
		log.Printf("Migration %s applied successfully", migration)
	}

	return nil
}
