package main

import (
	"Tournament/internal/config"
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	// Create database configuration
	dbConfig := config.NewDatabaseConfig()

	// Connect to database
	db, err := dbConfig.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test the connection with a simple query
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to create a test table
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS test_table (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create test table: %v", err)
	}

	// Insert a test record
	_, err = db.ExecContext(ctx, `
		INSERT INTO test_table (name) VALUES ($1)
	`, "Test record")
	if err != nil {
		log.Fatalf("Failed to insert test record: %v", err)
	}

	// Query the test record
	row := db.QueryRowContext(ctx, `
		SELECT id, name, created_at FROM test_table ORDER BY id DESC LIMIT 1
	`)

	var id int
	var name string
	var createdAt time.Time
	err = row.Scan(&id, &name, &createdAt)
	if err != nil {
		log.Fatalf("Failed to query test record: %v", err)
	}

	fmt.Printf("Test record: ID=%d, Name=%s, CreatedAt=%s\n", id, name, createdAt)
	fmt.Println("Database connection and queries are working correctly!")
}
