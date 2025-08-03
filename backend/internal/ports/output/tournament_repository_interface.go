package output

import (
	"Tournament/internal/domain"
	"context"
)

// TournamentRepositoryInterface defines the interface for tournament data access
type TournamentRepositoryInterface interface {
	// FindByID retrieves a tournament by its Id
	FindByID(ctx context.Context, id string) (*domain.Tournament, error)

	// FindAll retrieves all tournaments
	FindAll(ctx context.Context) ([]*domain.IndexTournament, error)

	//InsertNewTournament persists a tournament
	InsertNewTournament(ctx context.Context, tournament *domain.Tournament) (*domain.Tournament, error)

	// Delete removes a tournament
	Delete(ctx context.Context, id string) error

	// Update updates a tournament
	Update(ctx context.Context, tournament *domain.Tournament) (*domain.Tournament, error)
}
