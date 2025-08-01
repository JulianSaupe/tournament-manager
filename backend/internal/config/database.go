package config

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

// DatabaseConfig holds the configuration for the database connection
type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DBName       string
	SSLMode      string
	QueryTimeout time.Duration // Timeout for database queries
}

// NewDatabaseConfig creates a new database configuration with default values
func NewDatabaseConfig() *DatabaseConfig {
	// Default query timeout is 30 seconds
	queryTimeoutStr := getEnv("DB_QUERY_TIMEOUT", "30")
	queryTimeoutSec, err := strconv.Atoi(queryTimeoutStr)
	if err != nil {
		queryTimeoutSec = 30 // Default to 30 seconds if parsing fails
	}

	return &DatabaseConfig{
		Host:         getEnv("DB_HOST", "localhost"),
		Port:         getEnv("DB_PORT", "5432"),
		User:         getEnv("DB_USER", "postgres"),
		Password:     getEnv("DB_PASSWORD", "postgres"),
		DBName:       getEnv("DB_NAME", "tournament"),
		SSLMode:      getEnv("DB_SSLMODE", "disable"),
		QueryTimeout: time.Duration(queryTimeoutSec) * time.Second,
	}
}

// ConnectionString returns the PostgreSQL connection string
func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// NewBunDB creates a new Bun DB instance
func (c *DatabaseConfig) NewBunDB() (*bun.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode,
	)

	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	// Configure connection pool
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Set query timeout
	sqlDB.SetConnMaxIdleTime(c.QueryTimeout)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create a Bun DB instance
	db := bun.NewDB(sqlDB, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	return db, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
