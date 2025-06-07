package output

import (
	"Tournament/internal/domain"
)

// TournamentRepository defines the interface for tournament data access
type TournamentRepository interface {
	// FindByID retrieves a tournament by its ID
	FindByID(id string) (*domain.Tournament, error)
	
	// FindAll retrieves all tournaments
	FindAll() ([]*domain.Tournament, error)
	
	// Save persists a tournament
	Save(tournament *domain.Tournament) (*domain.Tournament, error)
	
	// Delete removes a tournament
	Delete(id string) error
}