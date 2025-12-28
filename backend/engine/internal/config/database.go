package config

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
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

// NewDB creates a new PostgreSQL database connection
func (c *DatabaseConfig) NewDB() (*sql.DB, error) {
	// Create connection string
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode,
	)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(c.QueryTimeout)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

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
