package input

import (
	"Tournament/internal/domain"
)

// TournamentService defines the interface for tournament business operations
type TournamentService interface {
	// CreateTournament creates a new tournament
	CreateTournament(name, description, startDate, endDate string) *domain.Tournament

	// GetTournament retrieves a tournament by Id
	GetTournament(id string) *domain.Tournament

	// ListTournaments retrieves all tournaments
	ListTournaments() []*domain.IndexTournament

	// UpdateTournamentStatus updates the status of a tournament
	UpdateTournamentStatus(id string, status domain.TournamentStatus) *domain.Tournament

	// DeleteTournament removes a tournament
	DeleteTournament(id string)
}
