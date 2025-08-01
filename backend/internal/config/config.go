package config

import (
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database *DatabaseConfig
}

// ServerConfig holds the configuration for the HTTP server
type ServerConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// Load loads the configuration from environment variables
func Load() *Config {
	return &Config{
		Server:   loadServerConfig(),
		Database: NewDatabaseConfig(), // Reuse existing function from database.go
	}
}

// loadServerConfig loads the server configuration from environment variables
func loadServerConfig() ServerConfig {
	readTimeoutSecs, _ := strconv.Atoi(getEnv("SERVER_READ_TIMEOUT", "30"))
	writeTimeoutSecs, _ := strconv.Atoi(getEnv("SERVER_WRITE_TIMEOUT", "30"))
	shutdownTimeoutSecs, _ := strconv.Atoi(getEnv("SERVER_SHUTDOWN_TIMEOUT", "30"))

	return ServerConfig{
		Port:            getEnv("SERVER_PORT", "3000"),
		ReadTimeout:     time.Duration(readTimeoutSecs) * time.Second,
		WriteTimeout:    time.Duration(writeTimeoutSecs) * time.Second,
		ShutdownTimeout: time.Duration(shutdownTimeoutSecs) * time.Second,
	}
}