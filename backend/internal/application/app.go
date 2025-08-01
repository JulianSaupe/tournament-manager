package application

import (
	"Tournament/internal/adapters/driven/postgres"
	"Tournament/internal/adapters/driving/response"
	httpHandler "Tournament/internal/adapters/driving/tournament"
	"Tournament/internal/application/service"
	"Tournament/internal/config"
	middleware2 "Tournament/internal/middleware"
	"Tournament/internal/ports/input"
	"Tournament/internal/ports/output"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/uptrace/bun"
	"log"
	"net/http"
)

// App represents the application with all its dependencies
type App struct {
	config *config.Config
	server *http.Server
	router *chi.Mux
	db     *bun.DB

	// Repositories
	tournamentRepository output.TournamentRepository
	userRepository       output.UserRepository
	playerRepository     output.PlayerRepository

	// Services
	tournamentService input.TournamentService
	userService       input.UserService
	playerService     input.PlayerService

	// Handlers
	tournamentHandler *httpHandler.Handler
}

// NewApp creates a new application instance
func NewApp(cfg *config.Config) (*App, error) {
	app := &App{
		config: cfg,
		router: chi.NewRouter(),
	}

	// Initialize dependencies
	if err := app.initializeDependencies(); err != nil {
		return nil, fmt.Errorf("failed to initialize dependencies: %w", err)
	}

	// Register routes
	app.registerRoutes()

	// Create HTTP server
	app.server = &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      app.router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return app, nil
}

// Start starts the application
func (a *App) Start() error {
	log.Printf("Starting server on :%s", a.config.Server.Port)
	return a.server.ListenAndServe()
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")

	// Close database connection
	if a.db != nil {
		log.Println("Closing database connection...")
		a.db.Close()
	}

	// Shutdown HTTP server
	return a.server.Shutdown(ctx)
}

// initializeDependencies initializes all dependencies
func (a *App) initializeDependencies() error {
	var err error

	// Initialize database
	a.db, err = a.config.Database.NewBunDB()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize repositories
	a.tournamentRepository, err = postgres.NewTournamentRepository(a.db)
	if err != nil {
		return fmt.Errorf("failed to initialize tournament repository: %w", err)
	}

	a.userRepository, err = postgres.NewPostgresUserRepository(a.db)
	if err != nil {
		return fmt.Errorf("failed to initialize user repository: %w", err)
	}

	a.playerRepository, err = postgres.NewPlayerRepository(a.db)
	if err != nil {
		return fmt.Errorf("failed to initialize player repository: %w", err)
	}

	// Initialize services
	a.tournamentService = service.NewTournamentService(a.tournamentRepository)
	a.userService = service.NewUserService(a.userRepository)
	a.playerService = service.NewPlayerService(a.playerRepository)

	// Initialize handlers
	a.tournamentHandler = httpHandler.NewTournamentHandler(a.tournamentService, a.playerService)

	return nil
}

// registerRoutes registers all HTTP routes
func (a *App) registerRoutes() {
	// Global middleware
	a.router.Use(chiMiddleware.RequestID)
	a.router.Use(chiMiddleware.RealIP)
	a.router.Use(chiMiddleware.Logger)
	a.router.Use(middleware2.AuthMiddleware(a.userService))
	a.router.Use(middleware2.CustomRecoverer)
	a.router.Use(response.RequestStartTimeMiddleware)

	// Health check endpoint
	a.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API routes
	apiRouter := chi.NewRouter()
	a.router.Mount("/api", apiRouter)

	// Register protected routes
	a.tournamentHandler.RegisterRoutes(apiRouter)
}
