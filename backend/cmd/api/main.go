package main

import (
	"Tournament/internal/adapters/driven/postgres"
	httpHandler "Tournament/internal/adapters/driving/http"
	"Tournament/internal/application"
	"Tournament/internal/config"
	"Tournament/internal/ports/output"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"path/filepath"
)

// initRepository initializes the PostgreSQL repository
func initRepository() (output.TournamentRepository, error) {
	log.Println("Using PostgreSQL repository")
	dbConfig := config.NewDatabaseConfig()
	return postgres.NewPostgresTournamentRepository(dbConfig.ConnectionString())
}

func main() {
	// Load .env file if it exists
	envPath := filepath.Join(".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Warning: .env file not found or could not be loaded: %v", err)
	}

	// Create repositories
	tournamentRepository, err := initRepository()
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	// Create services
	tournamentService := application.NewTournamentService(tournamentRepository)

	// Create handlers
	tournamentHandler := httpHandler.NewTournamentHandler(tournamentService)

	// Create router
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Create API subrouter with "/api" prefix
	apiRouter := chi.NewRouter()
	router.Mount("/api", apiRouter)

	// Register routes
	tournamentHandler.RegisterRoutes(apiRouter)

	// Start server
	log.Println("Starting server on :3000")

	serverErr := http.ListenAndServe(":3000", router)

	if serverErr != nil {
		log.Fatal(serverErr)
	}
}
