package main

import (
	"Tournament/internal/adapters/driven/postgres"
	httpHandler "Tournament/internal/adapters/driving/http"
	"Tournament/internal/application"
	"Tournament/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	dbConfig := config.NewDatabaseConfig()
	tournamentRepository, err := postgres.NewPostgresTournamentRepository(dbConfig.ConnectionString())

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
	router.Use(httpHandler.CustomRecoverer)
	router.Use(httpHandler.RequestStartTimeMiddleware)

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
