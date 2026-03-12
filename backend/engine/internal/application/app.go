package application

import (
	"context"
	"database/sql"
	"engine/backend/shared/proto/authorization"
	"engine/internal/adapters/driven/event"
	"engine/internal/adapters/driven/postgres"
	"engine/internal/adapters/driving/handler"
	"engine/internal/adapters/driving/response"
	"engine/internal/application/service"
	"engine/internal/config"
	"engine/internal/middleware"
	"engine/internal/ports/input"
	"engine/internal/ports/output"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
)

// App represents the application with all its dependencies
type App struct {
	config     *config.Config
	server     *http.Server
	grpcServer *grpc.Server
	router     *chi.Mux
	db         *sql.DB

	// Repositories
	tournamentRepository output.TournamentRepositoryInterface
	userRepository       output.UserRepositoryInterface
	playerRepository     output.PlayerRepositoryInterface
	qualifyingRepository output.QualifyingRepositoryInterface

	// Services
	tournamentService        input.TournamentServiceInterface
	userService              input.UserServiceInterface
	playerService            input.PlayerServiceInterface
	qualifyingService        input.QualifyingServiceInterface
	authorizationGRPCService *service.AuthorizationGRPCService

	// Handlers
	tournamentHandler *handler.TournamentHandler
	eventHandler      *handler.EventHandler

	// Broker
	broker *event.Broker
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
	// Start gRPC server in a goroutine
	go func() {
		if err := a.startGRPCServer(); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	// Start HTTP server
	log.Printf("Starting HTTP server on :%s", a.config.Server.Port)
	return a.server.ListenAndServe()
}

// startGRPCServer starts the gRPC server
func (a *App) startGRPCServer() error {
	lis, err := net.Listen("tcp", ":"+a.config.GRPC.Port)
	if err != nil {
		return fmt.Errorf("failed to listen on gRPC port: %w", err)
	}

	log.Printf("Starting gRPC server on :%s", a.config.GRPC.Port)
	return a.grpcServer.Serve(lis)
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown(ctx context.Context) error {
	log.Println("Shutting down servers...")

	// Shutdown gRPC server
	if a.grpcServer != nil {
		log.Println("Shutting down gRPC server...")
		a.grpcServer.GracefulStop()
	}

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

	// Initialize and start broker
	a.broker = event.NewBroker()
	a.broker.Start()

	// Initialize database
	a.db, err = a.config.Database.NewDB()
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

	a.qualifyingRepository, err = postgres.NewQualifyingRepository(a.db)
	if err != nil {
		return fmt.Errorf("failed to initialize qualifying repository: %w", err)
	}

	// Initialize services
	a.tournamentService = service.NewTournamentService(a.tournamentRepository, a.broker)
	a.userService = service.NewUserService(a.userRepository)
	a.playerService = service.NewPlayerService(a.playerRepository)
	a.qualifyingService = service.NewQualifyingService(a.qualifyingRepository)
	a.authorizationGRPCService = service.NewAuthorizationGRPCService()

	// Initialize handlers
	a.tournamentHandler = handler.NewTournamentHandler(a.tournamentService, a.playerService, a.qualifyingService)
	a.eventHandler = handler.NewEventHandler(a.broker)

	// Initialize gRPC server
	a.grpcServer = grpc.NewServer()
	authorization.RegisterAuthorizationServiceServer(a.grpcServer, a.authorizationGRPCService)

	return nil
}

// registerRoutes registers all HTTP routes
func (a *App) registerRoutes() {
	// Global middleware
	a.router.Use(chiMiddleware.RequestID)
	a.router.Use(chiMiddleware.RealIP)
	a.router.Use(chiMiddleware.Logger)
	// a.router.Use(middleware.AuthMiddleware(a.userService))
	a.router.Use(middleware.CustomRecoverer)
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
	a.eventHandler.RegisterRoutes(apiRouter)
}
