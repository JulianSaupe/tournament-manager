package main

import (
	"Tournament/internal/adapters/driven/postgres"
	"Tournament/internal/adapters/driving/response"
	httpHandler "Tournament/internal/adapters/driving/tournament"
	"Tournament/internal/application"
	"Tournament/internal/config"
	middleware2 "Tournament/internal/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	dbConfig := config.NewDatabaseConfig()

	// Initialize Bun DB
	db, err := dbConfig.NewBunDB()
	if err != nil {
		log.Fatalf("Failed to initialize Bun DB: %v", err)
	}
	defer db.Close()

	// Create repositories
	tournamentRepository, err := postgres.NewTournamentRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	userRepository, err := postgres.NewPostgresUserRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize user repository: %v", err)
	}

	playerRepository, err := postgres.NewPlayerRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize player repository: %v", err)
	}

	// Create services
	tournamentService := application.NewTournamentService(tournamentRepository)
	userService := application.NewUserService(userRepository)
	playerService := application.NewPlayerService(playerRepository)

	// Create handlers
	tournamentHandler := httpHandler.NewTournamentHandler(tournamentService, playerService)

	// Create router
	router := chi.NewRouter()

	// Middleware
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.RealIP)
	router.Use(chiMiddleware.Logger)
	router.Use(middleware2.AuthMiddleware(userService))
	router.Use(middleware2.CustomRecoverer)
	router.Use(response.RequestStartTimeMiddleware)

	apiRouter := chi.NewRouter()
	router.Mount("/api", apiRouter)

	// Register protected routes
	tournamentHandler.RegisterRoutes(apiRouter)

	// Start server
	log.Println("Starting server on :3000")

	serverErr := http.ListenAndServe(":3000", router)
	if serverErr != nil {
		log.Fatal(serverErr)
	}
}
