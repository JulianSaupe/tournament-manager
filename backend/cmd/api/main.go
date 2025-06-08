package main

import (
	"Tournament/internal/adapters/driven/postgres"
	"Tournament/internal/adapters/driving/middleware"
	"Tournament/internal/adapters/driving/response"
	httpHandler "Tournament/internal/adapters/driving/tournament"
	"Tournament/internal/application"
	"Tournament/internal/config"
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

	participantsRepository, err := postgres.NewParticipantsRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize participants repository: %v", err)
	}

	// Create services
	tournamentService := application.NewTournamentService(tournamentRepository)
	userService := application.NewUserService(userRepository)
	participantsService := application.NewParticipantsService(participantsRepository)

	// Create handlers
	tournamentHandler := httpHandler.NewTournamentHandler(tournamentService, participantsService)

	// Create router
	router := chi.NewRouter()

	// Middleware
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.RealIP)
	router.Use(chiMiddleware.Logger)
	router.Use(middleware.AuthMiddleware(userService))
	router.Use(middleware.CustomRecoverer)
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
