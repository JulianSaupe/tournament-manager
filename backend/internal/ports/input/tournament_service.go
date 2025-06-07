package input

import (
	"Tournament/internal/domain"
)

// TournamentService defines the interface for tournament business operations
type TournamentService interface {
	// CreateTournament creates a new tournament
	CreateTournament(name, description, startDate, endDate string) (*domain.Tournament, error)
	
	// GetTournament retrieves a tournament by ID
	GetTournament(id string) (*domain.Tournament, error)
	
	// ListTournaments retrieves all tournaments
	ListTournaments() ([]*domain.Tournament, error)
	
	// UpdateTournamentStatus updates the status of a tournament
	UpdateTournamentStatus(id string, status domain.TournamentStatus) (*domain.Tournament, error)
}