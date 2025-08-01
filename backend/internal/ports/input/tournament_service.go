package input

import (
	"Tournament/internal/domain"
	"context"
)

// TournamentService defines the interface for tournament business operations
type TournamentService interface {
	// CreateTournament creates a new tournament
	CreateTournament(ctx context.Context, name, description, startDate, endDate string) *domain.Tournament

	// GetTournament retrieves a tournament by Id
	GetTournament(ctx context.Context, id string) *domain.Tournament

	// ListTournaments retrieves all tournaments
	ListTournaments(ctx context.Context) []*domain.IndexTournament

	// UpdateTournamentStatus updates the status of a tournament
	UpdateTournamentStatus(ctx context.Context, id string, status domain.TournamentStatus) *domain.Tournament

	// DeleteTournament removes a tournament
	DeleteTournament(ctx context.Context, id string)
}
